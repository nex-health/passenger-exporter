// Copyright 2024 NexHealth Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//

package collector

import (
	"math"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "passenger"

	nanosecondsPerSecond = 1000000000
)

var (
	up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Passenger state.",
		[]string{}, nil,
	)
	version = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "version"),
		"Phusion Passenger version.",
		[]string{"version"}, nil,
	)
	toplevelQueue = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "top_level_queue"),
		"Number of requests in the top-level queue.",
		[]string{}, nil,
	)
	maxProcessCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "max_processes"),
		"Configured maximum number of processes.",
		[]string{}, nil,
	)
	currentProcessCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "current_processes"),
		"Current number of processes.",
		[]string{}, nil,
	)
	appCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "app_count"),
		"Number of apps.",
		[]string{}, nil,
	)
	appQueue = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "app_queue"),
		"Number of requests in app process queues.",
		[]string{"name"}, nil,
	)
	appGroupQueue = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "app_group_queue"),
		"Number of requests in app group process queues.",
		[]string{"group", "default"}, nil,
	)
	appProcsSpawning = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "app_procs_spawning"),
		"Number of processes spawning.",
		[]string{"name"}, nil,
	)
	requestsProcessed = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "requests_processed_total"),
		"Number of processes served by a process.",
		[]string{"name", "id"}, nil,
	)
	sessions = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "current_sessions"),
		"Number of sessions currently being handled by a process.",
		[]string{"name", "id"}, nil,
	)
	procStartTime = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "proc_start_time_seconds"),
		"Number of seconds since processor started.",
		[]string{"name", "id", "codeRevision"}, nil,
	)
	procMemory = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "proc_memory"),
		"Memory consumed by a process",
		[]string{"name", "id"}, nil,
	)
)

type Collector struct {
	MetricsReader

	PIDFile string
}

func New(reader MetricsReader, pidFile string) *Collector {
	return &Collector{MetricsReader: reader, PIDFile: pidFile}
}

func (c Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	ch <- version
	ch <- toplevelQueue
	ch <- maxProcessCount
	ch <- currentProcessCount
	ch <- appCount
	ch <- appQueue
	ch <- appGroupQueue
	ch <- appProcsSpawning
	ch <- requestsProcessed
	ch <- sessions
	ch <- procStartTime
	ch <- procMemory
}

// Mostly copied from https://github.com/stuartnelson3/passenger_exporter/blob/80b16566cdab445f6e68f967019a95b67f608aca/main.go
func (c Collector) Collect(ch chan<- prometheus.Metric) {
	var processIdentifiers map[string]int
	data, err := c.MetricsReader.Read()
	if err != nil {
		ch <- prometheus.NewInvalidMetric(prometheus.NewDesc(prometheus.BuildFQName(namespace, "read", "error"), "Error reading metrics data.", nil, nil), err)
		return
	}
	defer data.Close()

	info, err := Parse(data)
	if err != nil {
		ch <- prometheus.NewInvalidMetric(prometheus.NewDesc(prometheus.BuildFQName(namespace, "parse", "error"), "Error parsing metrics data.", nil, nil), err)
		return
	}

	ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 1)

	ch <- prometheus.MustNewConstMetric(version, prometheus.GaugeValue, 1, info.PassengerVersion)

	ch <- prometheus.MustNewConstMetric(toplevelQueue, prometheus.GaugeValue, parseFloat(info.TopLevelRequestsInQueue))
	ch <- prometheus.MustNewConstMetric(maxProcessCount, prometheus.GaugeValue, parseFloat(info.MaxProcessCount))
	ch <- prometheus.MustNewConstMetric(currentProcessCount, prometheus.GaugeValue, parseFloat(info.CurrentProcessCount))
	ch <- prometheus.MustNewConstMetric(appCount, prometheus.GaugeValue, parseFloat(info.AppCount))

	for _, sg := range info.SuperGroups {
		ch <- prometheus.MustNewConstMetric(appQueue, prometheus.GaugeValue, parseFloat(sg.RequestsInQueue), sg.Name)
		ch <- prometheus.MustNewConstMetric(appProcsSpawning, prometheus.GaugeValue, parseFloat(sg.Group.ProcessesSpawning), sg.Name)

		ch <- prometheus.MustNewConstMetric(appGroupQueue, prometheus.GaugeValue, parseFloat(sg.Group.GetWaitListSize), sg.Group.Name, sg.Group.Default)

		// Update process identifiers map.
		processIdentifiers := updateProcesses(processIdentifiers, sg.Group.Processes)
		for _, proc := range sg.Group.Processes {
			if bucketID, ok := processIdentifiers[proc.PID]; ok {
				ch <- prometheus.MustNewConstMetric(procMemory, prometheus.GaugeValue, parseFloat(proc.RealMemory), sg.Name, strconv.Itoa(bucketID))
				ch <- prometheus.MustNewConstMetric(requestsProcessed, prometheus.CounterValue, parseFloat(proc.RequestsProcessed), sg.Name, strconv.Itoa(bucketID))
				ch <- prometheus.MustNewConstMetric(sessions, prometheus.GaugeValue, parseFloat(proc.Sessions), sg.Name, strconv.Itoa(bucketID))

				if startTime, err := strconv.Atoi(proc.SpawnStartTime); err == nil {
					ch <- prometheus.MustNewConstMetric(procStartTime, prometheus.GaugeValue, float64(startTime/nanosecondsPerSecond),
						sg.Name, strconv.Itoa(bucketID), proc.CodeRevision,
					)
				}
			}
		}
	}
}

// Copied from https://github.com/stuartnelson3/passenger_exporter/blob/80b16566cdab445f6e68f967019a95b67f608aca/main.go
// updateProcesses updates the global map from process id:exporter id. Process
// TTLs cause new processes to be created on a user-defined cycle. When a new
// process replaces an old process, the new process's statistics will be
// bucketed with those of the process it replaced.
// Processes are restarted at an offset, user-defined interval. The
// restarted process is appended to the end of the status output.  For
// maintaining consistent process identifiers between process starts,
// pids are mapped to an identifier based on process count. When a new
// process/pid appears, it is mapped to either the first empty place
// within the global map storing process identifiers, or mapped to
// pid:id pair in the map.
func updateProcesses(old map[string]int, processes []Process) map[string]int {
	var (
		updated = make(map[string]int)
		found   = make([]string, len(old))
		missing []string
	)

	for _, p := range processes {
		if id, ok := old[p.PID]; ok {
			found[id] = p.PID
			// id also serves as an index.
			// By putting the pid at a certain index, we can loop
			// through the array to find the values that are the 0
			// value (empty string).
			// If index i has the empty value, then it was never
			// updated, so we slot the first of the missingPIDs
			// into that position. Passenger-status orders output
			// by pid, increasing. We can then assume that
			// unclaimed pid positions map in order to the missing
			// pids.
		} else {
			missing = append(missing, p.PID)
		}
	}

	j := 0
	for i, pid := range found {
		if pid == "" {
			if j >= len(missing) {
				continue
			}
			pid = missing[j]
			j++
		}
		updated[pid] = i
	}

	// If the number of elements in missing iterated through is less
	// than len(missing), there are new elements to be added to the map.
	// Unused pids from the last collection are not copied from old to
	// updated, thereby cleaning the return value of unused PIDs.
	if j < len(missing) {
		count := len(found)
		for i, pid := range missing[j:] {
			updated[pid] = count + i
		}
	}

	return updated
}

func parseFloat(val string) float64 {
	v, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return math.NaN()
	}
	return v
}
