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

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/nex-health/passenger-exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promslog"
	"github.com/prometheus/common/promslog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"
)

func main() {
	var (
		webConfig   = webflag.AddFlags(kingpin.CommandLine, ":9149")
		metricsPath = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()

		instanceRegistry = kingpin.Flag("passenger.instance-registry", "Path to the instance registry directory.").Default(os.TempDir()).String()
		pidFile          = kingpin.Flag("passenger.pid-file", "Optional path to a file containing the passenger/nginx PID for additional metrics.").Default("").String()
	)

	promslogConfig := &promslog.Config{}
	flag.AddFlags(kingpin.CommandLine, promslogConfig)
	kingpin.Version(version.Print("passenger_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promslog.New(promslogConfig)

	logger.Info("Starting passenger_exporter", "version", version.Info())
	logger.Info("Build context", "context", version.BuildContext())

	if *pidFile != "" {
		pidCollector := collectors.NewProcessCollector(collectors.ProcessCollectorOpts{
			PidFn:     prometheus.NewPidFileFn(*pidFile),
			Namespace: "passenger",
		})
		prometheus.MustRegister(pidCollector)
	}

	udsReader := collector.NewUDSReader(*instanceRegistry)
	collector := collector.New(udsReader)
	prometheus.MustRegister(collector)

	http.Handle(*metricsPath, promhttp.Handler())

	landingConfig := web.LandingConfig{
		Name:        "Phusion Passenger Exporter",
		Description: "Prometheus Exporter for Phusion Passenger",
		Version:     version.Info(),
		Links: []web.LandingLinks{
			{
				Address: *metricsPath,
				Text:    "Metrics",
			},
		},
		ExtraHTML: fmt.Sprintf(`<h2>Options</h2><pre>passenger.instance-registry: "%s", passenger.pid-file: "%s"</pre>`, *instanceRegistry, *pidFile),
	}
	landingHandler, err := web.NewLandingPage(landingConfig)
	if err != nil {
		logger.Error("Error creating landing page", "err", err)
		os.Exit(1)
	}
	http.Handle("/", landingHandler)
	http.HandleFunc("/-/healthy", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})
	http.HandleFunc("/-/ready", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	srv := &http.Server{}
	if err := web.ListenAndServe(srv, webConfig, logger); err != nil {
		logger.Error("Error starting HTTP server", "err", err)
		os.Exit(1)
	}
}
