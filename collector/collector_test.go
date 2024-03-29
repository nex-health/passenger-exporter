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

package collector

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

type fakeReader struct {
	ReaderFunc func() (io.ReadCloser, error)
}

func (r *fakeReader) Read() (io.ReadCloser, error) {
	return r.ReaderFunc()
}

func TestNew(t *testing.T) {
	reader := &fakeReader{}
	c := New(reader)

	if c.reader != reader {
		t.Errorf("expected reader field to equal given argument")
	}

	hostname, _ := os.Hostname()
	c = New(reader)

	if c.hostname != hostname {
		t.Errorf("expected hostname field to equal %q, got %q", hostname, c.hostname)
	}

	t.Setenv("HOSTNAME", "fake-hostname")
	c = New(reader)

	if c.hostname != "fake-hostname" {
		t.Errorf("expected hostname field to equal %q, got %q", "fake-hostname", c.hostname)
	}
}

func TestCollect(t *testing.T) {
	fixture, _ := os.Open("testdata/passenger_xml_output.xml")

	t.Setenv("HOSTNAME", "local-machine")

	for _, tc := range []struct {
		name        string
		readerFunc  func() (io.ReadCloser, error)
		status      int
		wantMetrics string
		wantErr     bool
	}{
		{
			name: "collect with valid response",
			wantMetrics: `# HELP passenger_app_count Number of apps.
# TYPE passenger_app_count gauge
passenger_app_count{hostname="local-machine"} 1
# HELP passenger_app_group_queue Number of requests in app group process queues.
# TYPE passenger_app_group_queue gauge
passenger_app_group_queue{default="true",group="/srv/app/my_app (production)",hostname="local-machine"} 0
# HELP passenger_app_procs_spawning Number of processes spawning.
# TYPE passenger_app_procs_spawning gauge
passenger_app_procs_spawning{hostname="local-machine",name="/srv/app/my_app (production)"} 0
# HELP passenger_app_queue Number of requests in app process queues.
# TYPE passenger_app_queue gauge
passenger_app_queue{hostname="local-machine",name="/srv/app/my_app (production)"} 5
# HELP passenger_current_processes Current number of processes.
# TYPE passenger_current_processes gauge
passenger_current_processes{hostname="local-machine"} 48
# HELP passenger_current_sessions Number of sessions currently being handled by a process.
# TYPE passenger_current_sessions gauge
passenger_current_sessions{hostname="local-machine",id="0",name="/srv/app/my_app (production)"} 1
passenger_current_sessions{hostname="local-machine",id="1",name="/srv/app/my_app (production)"} 1
passenger_current_sessions{hostname="local-machine",id="10",name="/srv/app/my_app (production)"} 1
passenger_current_sessions{hostname="local-machine",id="11",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="12",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="13",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="14",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="15",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="16",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="17",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="18",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="19",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="2",name="/srv/app/my_app (production)"} 1
passenger_current_sessions{hostname="local-machine",id="20",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="21",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="22",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="23",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="24",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="25",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="26",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="27",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="28",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="29",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="3",name="/srv/app/my_app (production)"} 1
passenger_current_sessions{hostname="local-machine",id="30",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="31",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="32",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="33",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="34",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="35",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="36",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="37",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="38",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="39",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="4",name="/srv/app/my_app (production)"} 1
passenger_current_sessions{hostname="local-machine",id="40",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="41",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="42",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="43",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="44",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="45",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="46",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="47",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="5",name="/srv/app/my_app (production)"} 1
passenger_current_sessions{hostname="local-machine",id="6",name="/srv/app/my_app (production)"} 1
passenger_current_sessions{hostname="local-machine",id="7",name="/srv/app/my_app (production)"} 0
passenger_current_sessions{hostname="local-machine",id="8",name="/srv/app/my_app (production)"} 1
passenger_current_sessions{hostname="local-machine",id="9",name="/srv/app/my_app (production)"} 1
# HELP passenger_max_processes Configured maximum number of processes.
# TYPE passenger_max_processes gauge
passenger_max_processes{hostname="local-machine"} 48
# HELP passenger_proc_memory Memory consumed by a process
# TYPE passenger_proc_memory gauge
passenger_proc_memory{hostname="local-machine",id="0",name="/srv/app/my_app (production)"} 330012
passenger_proc_memory{hostname="local-machine",id="1",name="/srv/app/my_app (production)"} 303296
passenger_proc_memory{hostname="local-machine",id="10",name="/srv/app/my_app (production)"} 303984
passenger_proc_memory{hostname="local-machine",id="11",name="/srv/app/my_app (production)"} 289680
passenger_proc_memory{hostname="local-machine",id="12",name="/srv/app/my_app (production)"} 306148
passenger_proc_memory{hostname="local-machine",id="13",name="/srv/app/my_app (production)"} 293128
passenger_proc_memory{hostname="local-machine",id="14",name="/srv/app/my_app (production)"} 322064
passenger_proc_memory{hostname="local-machine",id="15",name="/srv/app/my_app (production)"} 297124
passenger_proc_memory{hostname="local-machine",id="16",name="/srv/app/my_app (production)"} 290364
passenger_proc_memory{hostname="local-machine",id="17",name="/srv/app/my_app (production)"} 292056
passenger_proc_memory{hostname="local-machine",id="18",name="/srv/app/my_app (production)"} 272784
passenger_proc_memory{hostname="local-machine",id="19",name="/srv/app/my_app (production)"} 281176
passenger_proc_memory{hostname="local-machine",id="2",name="/srv/app/my_app (production)"} 288884
passenger_proc_memory{hostname="local-machine",id="20",name="/srv/app/my_app (production)"} 269520
passenger_proc_memory{hostname="local-machine",id="21",name="/srv/app/my_app (production)"} 269404
passenger_proc_memory{hostname="local-machine",id="22",name="/srv/app/my_app (production)"} 275844
passenger_proc_memory{hostname="local-machine",id="23",name="/srv/app/my_app (production)"} 276412
passenger_proc_memory{hostname="local-machine",id="24",name="/srv/app/my_app (production)"} 267316
passenger_proc_memory{hostname="local-machine",id="25",name="/srv/app/my_app (production)"} 265152
passenger_proc_memory{hostname="local-machine",id="26",name="/srv/app/my_app (production)"} 261144
passenger_proc_memory{hostname="local-machine",id="27",name="/srv/app/my_app (production)"} 260224
passenger_proc_memory{hostname="local-machine",id="28",name="/srv/app/my_app (production)"} 243688
passenger_proc_memory{hostname="local-machine",id="29",name="/srv/app/my_app (production)"} 243724
passenger_proc_memory{hostname="local-machine",id="3",name="/srv/app/my_app (production)"} 293316
passenger_proc_memory{hostname="local-machine",id="30",name="/srv/app/my_app (production)"} 261492
passenger_proc_memory{hostname="local-machine",id="31",name="/srv/app/my_app (production)"} 260196
passenger_proc_memory{hostname="local-machine",id="32",name="/srv/app/my_app (production)"} 244720
passenger_proc_memory{hostname="local-machine",id="33",name="/srv/app/my_app (production)"} 261268
passenger_proc_memory{hostname="local-machine",id="34",name="/srv/app/my_app (production)"} 261320
passenger_proc_memory{hostname="local-machine",id="35",name="/srv/app/my_app (production)"} 244740
passenger_proc_memory{hostname="local-machine",id="36",name="/srv/app/my_app (production)"} 244656
passenger_proc_memory{hostname="local-machine",id="37",name="/srv/app/my_app (production)"} 244860
passenger_proc_memory{hostname="local-machine",id="38",name="/srv/app/my_app (production)"} 244752
passenger_proc_memory{hostname="local-machine",id="39",name="/srv/app/my_app (production)"} 244708
passenger_proc_memory{hostname="local-machine",id="4",name="/srv/app/my_app (production)"} 330412
passenger_proc_memory{hostname="local-machine",id="40",name="/srv/app/my_app (production)"} 244684
passenger_proc_memory{hostname="local-machine",id="41",name="/srv/app/my_app (production)"} 255428
passenger_proc_memory{hostname="local-machine",id="42",name="/srv/app/my_app (production)"} 243744
passenger_proc_memory{hostname="local-machine",id="43",name="/srv/app/my_app (production)"} 254432
passenger_proc_memory{hostname="local-machine",id="44",name="/srv/app/my_app (production)"} 243592
passenger_proc_memory{hostname="local-machine",id="45",name="/srv/app/my_app (production)"} 244640
passenger_proc_memory{hostname="local-machine",id="46",name="/srv/app/my_app (production)"} 242576
passenger_proc_memory{hostname="local-machine",id="47",name="/srv/app/my_app (production)"} 255376
passenger_proc_memory{hostname="local-machine",id="5",name="/srv/app/my_app (production)"} 306904
passenger_proc_memory{hostname="local-machine",id="6",name="/srv/app/my_app (production)"} 330644
passenger_proc_memory{hostname="local-machine",id="7",name="/srv/app/my_app (production)"} 315104
passenger_proc_memory{hostname="local-machine",id="8",name="/srv/app/my_app (production)"} 288508
passenger_proc_memory{hostname="local-machine",id="9",name="/srv/app/my_app (production)"} 306520
# HELP passenger_proc_start_time_seconds Number of seconds since processor started.
# TYPE passenger_proc_start_time_seconds gauge
passenger_proc_start_time_seconds{hostname="local-machine",id="0",name="/srv/app/my_app (production)"} 1.462477e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="1",name="/srv/app/my_app (production)"} 1.462477e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="10",name="/srv/app/my_app (production)"} 1.462477e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="11",name="/srv/app/my_app (production)"} 1.462477e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="12",name="/srv/app/my_app (production)"} 1.462477e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="13",name="/srv/app/my_app (production)"} 1.462477e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="14",name="/srv/app/my_app (production)"} 1.462477e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="15",name="/srv/app/my_app (production)"} 1.462477e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="16",name="/srv/app/my_app (production)"} 1.462477e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="17",name="/srv/app/my_app (production)"} 1.462477e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="18",name="/srv/app/my_app (production)"} 1.462477e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="19",name="/srv/app/my_app (production)"} 1.462477e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="2",name="/srv/app/my_app (production)"} 1.462477e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="20",name="/srv/app/my_app (production)"} 1.462477e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="21",name="/srv/app/my_app (production)"} 1.462477e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="22",name="/srv/app/my_app (production)"} 1.462477e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="23",name="/srv/app/my_app (production)"} 1.462477e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="24",name="/srv/app/my_app (production)"} 1.462477e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="25",name="/srv/app/my_app (production)"} 1.462477e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="26",name="/srv/app/my_app (production)"} 1.462477e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="27",name="/srv/app/my_app (production)"} 1.462477e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="28",name="/srv/app/my_app (production)"} 1.462477e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="29",name="/srv/app/my_app (production)"} 1.462477e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="3",name="/srv/app/my_app (production)"} 1.462477e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="30",name="/srv/app/my_app (production)"} 1.462477e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="31",name="/srv/app/my_app (production)"} 1.462477e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="32",name="/srv/app/my_app (production)"} 1.462477e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="33",name="/srv/app/my_app (production)"} 1.462477e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="34",name="/srv/app/my_app (production)"} 1.462477e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="35",name="/srv/app/my_app (production)"} 1.462477e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="36",name="/srv/app/my_app (production)"} 1.462477e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="37",name="/srv/app/my_app (production)"} 1.462478e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="38",name="/srv/app/my_app (production)"} 1.462478e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="39",name="/srv/app/my_app (production)"} 1.462478e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="4",name="/srv/app/my_app (production)"} 1.462477e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="40",name="/srv/app/my_app (production)"} 1.462478e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="41",name="/srv/app/my_app (production)"} 1.462478e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="42",name="/srv/app/my_app (production)"} 1.462478e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="43",name="/srv/app/my_app (production)"} 1.462478e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="44",name="/srv/app/my_app (production)"} 1.462478e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="45",name="/srv/app/my_app (production)"} 1.462478e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="46",name="/srv/app/my_app (production)"} 1.462478e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="47",name="/srv/app/my_app (production)"} 1.462478e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="5",name="/srv/app/my_app (production)"} 1.462477e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="6",name="/srv/app/my_app (production)"} 1.462477e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="7",name="/srv/app/my_app (production)"} 1.462477e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="8",name="/srv/app/my_app (production)"} 1.462477e+06
passenger_proc_start_time_seconds{hostname="local-machine",id="9",name="/srv/app/my_app (production)"} 1.462477e+06
# HELP passenger_requests_processed_total Number of processes served by a process.
# TYPE passenger_requests_processed_total counter
passenger_requests_processed_total{hostname="local-machine",id="0",name="/srv/app/my_app (production)"} 43578
passenger_requests_processed_total{hostname="local-machine",id="1",name="/srv/app/my_app (production)"} 48130
passenger_requests_processed_total{hostname="local-machine",id="10",name="/srv/app/my_app (production)"} 26226
passenger_requests_processed_total{hostname="local-machine",id="11",name="/srv/app/my_app (production)"} 22752
passenger_requests_processed_total{hostname="local-machine",id="12",name="/srv/app/my_app (production)"} 18646
passenger_requests_processed_total{hostname="local-machine",id="13",name="/srv/app/my_app (production)"} 15254
passenger_requests_processed_total{hostname="local-machine",id="14",name="/srv/app/my_app (production)"} 11561
passenger_requests_processed_total{hostname="local-machine",id="15",name="/srv/app/my_app (production)"} 9107
passenger_requests_processed_total{hostname="local-machine",id="16",name="/srv/app/my_app (production)"} 6831
passenger_requests_processed_total{hostname="local-machine",id="17",name="/srv/app/my_app (production)"} 4804
passenger_requests_processed_total{hostname="local-machine",id="18",name="/srv/app/my_app (production)"} 3420
passenger_requests_processed_total{hostname="local-machine",id="19",name="/srv/app/my_app (production)"} 2150
passenger_requests_processed_total{hostname="local-machine",id="2",name="/srv/app/my_app (production)"} 46701
passenger_requests_processed_total{hostname="local-machine",id="20",name="/srv/app/my_app (production)"} 1333
passenger_requests_processed_total{hostname="local-machine",id="21",name="/srv/app/my_app (production)"} 809
passenger_requests_processed_total{hostname="local-machine",id="22",name="/srv/app/my_app (production)"} 504
passenger_requests_processed_total{hostname="local-machine",id="23",name="/srv/app/my_app (production)"} 288
passenger_requests_processed_total{hostname="local-machine",id="24",name="/srv/app/my_app (production)"} 161
passenger_requests_processed_total{hostname="local-machine",id="25",name="/srv/app/my_app (production)"} 99
passenger_requests_processed_total{hostname="local-machine",id="26",name="/srv/app/my_app (production)"} 60
passenger_requests_processed_total{hostname="local-machine",id="27",name="/srv/app/my_app (production)"} 49
passenger_requests_processed_total{hostname="local-machine",id="28",name="/srv/app/my_app (production)"} 24
passenger_requests_processed_total{hostname="local-machine",id="29",name="/srv/app/my_app (production)"} 19
passenger_requests_processed_total{hostname="local-machine",id="3",name="/srv/app/my_app (production)"} 45134
passenger_requests_processed_total{hostname="local-machine",id="30",name="/srv/app/my_app (production)"} 9
passenger_requests_processed_total{hostname="local-machine",id="31",name="/srv/app/my_app (production)"} 5
passenger_requests_processed_total{hostname="local-machine",id="32",name="/srv/app/my_app (production)"} 4
passenger_requests_processed_total{hostname="local-machine",id="33",name="/srv/app/my_app (production)"} 4
passenger_requests_processed_total{hostname="local-machine",id="34",name="/srv/app/my_app (production)"} 2
passenger_requests_processed_total{hostname="local-machine",id="35",name="/srv/app/my_app (production)"} 2
passenger_requests_processed_total{hostname="local-machine",id="36",name="/srv/app/my_app (production)"} 0
passenger_requests_processed_total{hostname="local-machine",id="37",name="/srv/app/my_app (production)"} 0
passenger_requests_processed_total{hostname="local-machine",id="38",name="/srv/app/my_app (production)"} 0
passenger_requests_processed_total{hostname="local-machine",id="39",name="/srv/app/my_app (production)"} 0
passenger_requests_processed_total{hostname="local-machine",id="4",name="/srv/app/my_app (production)"} 42932
passenger_requests_processed_total{hostname="local-machine",id="40",name="/srv/app/my_app (production)"} 0
passenger_requests_processed_total{hostname="local-machine",id="41",name="/srv/app/my_app (production)"} 0
passenger_requests_processed_total{hostname="local-machine",id="42",name="/srv/app/my_app (production)"} 0
passenger_requests_processed_total{hostname="local-machine",id="43",name="/srv/app/my_app (production)"} 0
passenger_requests_processed_total{hostname="local-machine",id="44",name="/srv/app/my_app (production)"} 0
passenger_requests_processed_total{hostname="local-machine",id="45",name="/srv/app/my_app (production)"} 0
passenger_requests_processed_total{hostname="local-machine",id="46",name="/srv/app/my_app (production)"} 0
passenger_requests_processed_total{hostname="local-machine",id="47",name="/srv/app/my_app (production)"} 0
passenger_requests_processed_total{hostname="local-machine",id="5",name="/srv/app/my_app (production)"} 40815
passenger_requests_processed_total{hostname="local-machine",id="6",name="/srv/app/my_app (production)"} 38615
passenger_requests_processed_total{hostname="local-machine",id="7",name="/srv/app/my_app (production)"} 35802
passenger_requests_processed_total{hostname="local-machine",id="8",name="/srv/app/my_app (production)"} 33600
passenger_requests_processed_total{hostname="local-machine",id="9",name="/srv/app/my_app (production)"} 30490
# HELP passenger_top_level_queue Number of requests in the top-level queue.
# TYPE passenger_top_level_queue gauge
passenger_top_level_queue{hostname="local-machine"} 3
# HELP passenger_up Passenger state.
# TYPE passenger_up gauge
passenger_up{hostname="local-machine"} 1
# HELP passenger_version Phusion Passenger version.
# TYPE passenger_version gauge
passenger_version{hostname="local-machine",version="5.0.26"} 1
`,
			readerFunc: func() (io.ReadCloser, error) { return fixture, nil },
			status:     http.StatusOK,
			wantErr:    false,
		},
		{
			name:        "collect with error response",
			wantMetrics: ``,
			readerFunc:  func() (io.ReadCloser, error) { return fixture, nil },
			status:      http.StatusInternalServerError,
			wantErr:     true,
		},
		{
			name:        "collect with reader error",
			wantMetrics: ``,
			readerFunc:  func() (io.ReadCloser, error) { return nil, fmt.Errorf("fake") },
			status:      http.StatusInternalServerError,
			wantErr:     true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			reader := &fakeReader{ReaderFunc: tc.readerFunc}
			collector := New(reader)

			buf := bytes.NewReader([]byte(tc.wantMetrics))
			err := testutil.CollectAndCompare(collector, buf)
			if tc.wantErr && err == nil {
				t.Errorf("expected error, but got %q", err)
			}
			if !tc.wantErr && err != nil {
				t.Errorf("expected no error, but got %q", err)
			}
		})
	}
}

// The below code was copied from https://github.com/stuartnelson3/passenger_exporter/blob/80b16566cdab445f6e68f967019a95b67f608aca/main_test.go

type updateProcessSpec struct {
	name      string
	input     map[string]int
	processes []Process
	output    map[string]int
}

func newUpdateProcessSpec(
	name string,
	input map[string]int,
	processes []Process,
) updateProcessSpec {
	s := updateProcessSpec{
		name:      name,
		input:     input,
		processes: processes,
	}
	s.output = updateProcesses(s.input, s.processes)
	return s
}

func TestUpdateProcessIdentifiers(t *testing.T) {
	for _, spec := range []updateProcessSpec{
		newUpdateProcessSpec(
			"empty input",
			map[string]int{},
			[]Process{
				{PID: "abc"},
				{PID: "cdf"},
				{PID: "dfe"},
			},
		),
		newUpdateProcessSpec(
			"1:1",
			map[string]int{
				"abc": 0,
				"cdf": 1,
				"dfe": 2,
			},
			[]Process{
				{PID: "abc"},
				{PID: "cdf"},
				{PID: "dfe"},
			},
		),
		newUpdateProcessSpec(
			"increase processes",
			map[string]int{
				"abc": 0,
				"cdf": 1,
				"dfe": 2,
			},
			[]Process{
				{PID: "abc"},
				{PID: "cdf"},
				{PID: "dfe"},
				{PID: "ghi"},
				{PID: "jkl"},
				{PID: "lmn"},
			},
		),
		newUpdateProcessSpec(
			"reduce processes",
			map[string]int{
				"abc": 0,
				"cdf": 1,
				"dfe": 2,
				"ghi": 3,
				"jkl": 4,
				"lmn": 5,
			},
			[]Process{
				{PID: "abc"},
				{PID: "cdf"},
				{PID: "dfe"},
			},
		),
	} {
		if len(spec.output) != len(spec.processes) {
			t.Fatalf("case %s: proceses improperly copied to output: len(output) (%d) does not match len(processes) (%d)", spec.name, len(spec.output), len(spec.processes))
		}

		for _, p := range spec.processes {
			if _, ok := spec.output[p.PID]; !ok {
				t.Fatalf("case %s: pid not copied into map", spec.name)
			}
		}

		newOutput := updateProcesses(spec.output, spec.processes)
		if !reflect.DeepEqual(newOutput, spec.output) {
			t.Fatalf("case %s: updateProcesses is not idempotent", spec.name)
		}
	}
}

func TestInsertingNewProcesses(t *testing.T) {
	spec := newUpdateProcessSpec(
		"inserting processes",
		map[string]int{
			"abc": 0,
			"cdf": 1,
			"dfe": 2,
			"efg": 3,
		},
		[]Process{
			{PID: "abc"},
			{PID: "dfe"},
			{PID: "newPID"},
			{PID: "newPID2"},
		},
	)

	if len(spec.output) != len(spec.processes) {
		t.Fatalf("case %s: proceses improperly copied to output: len(output) (%d) does not match len(processes) (%d)", spec.name, len(spec.output), len(spec.processes))
	}

	if want, got := 1, spec.output["newPID"]; want != got {
		t.Fatalf("updateProcesses did not correctly map the new PID: wanted %d, got %d", want, got)
	}
	if want, got := 3, spec.output["newPID2"]; want != got {
		t.Fatalf("updateProcesses did not correctly map the new PID: wanted %d, got %d", want, got)
	}
}

func newTestCollector() *Collector {
	return New(&fakeReader{})
}
