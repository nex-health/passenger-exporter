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
	"math"
	"os"
	"testing"
)

func TestParsing(t *testing.T) {
	fixture, err := os.Open("testdata/passenger_xml_output.xml")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	info, err := Parse(fixture)
	if err != nil {
		t.Fatalf("parse xml file failed: %v", err)
	}
	if len(info.SuperGroups) == 0 {
		t.Fatalf("no supergroups in output")
	}

	topLevelQueue := parseFloat(info.TopLevelRequestsInQueue)
	if topLevelQueue == 0 {
		t.Fatalf("no queuing requests parsed from output")
	}

	for _, sg := range info.SuperGroups {
		if want, got := "/src/app/my_app", sg.Group.Options.AppRoot; want != got {
			t.Fatalf("incorrect app_root: wanted %s, got %s", want, got)
		}

		queue := parseFloat(sg.RequestsInQueue)
		if math.IsNaN(queue) {
			t.Fatalf("failed to parse requests in queue")
		}

		if len(sg.Group.Processes) == 0 {
			t.Fatalf("no processes in output")
		}
		for _, proc := range sg.Group.Processes {
			if want, got := "2254", proc.ProcessGroupID; want != got {
				t.Fatalf("incorrect process_group_id: wanted %s, got %s", want, got)
			}
		}
	}
}
