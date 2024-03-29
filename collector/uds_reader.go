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
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	UDSPath                   = "agents.s/core_api"
	ReadOnlyAdminUsername     = "ro_admin"
	ReadOnlyAdminPasswordFile = "read_only_admin_password.txt"
)

type UDSReader struct {
	Path string

	sync.Mutex
}

func NewUDSReader(path string) *UDSReader {
	return &UDSReader{
		Path: path,
	}
}

func (r *UDSReader) Read() (io.ReadCloser, error) {
	r.Lock()
	defer r.Unlock()

	registry, err := filepath.Glob(filepath.Join(r.Path, "passenger.???????"))
	if err != nil {
		return nil, err
	}

	if len(registry) == 0 {
		return nil, fmt.Errorf("failed to detect Passenger instance registry directory")
	}

	uds := filepath.Join(registry[0], UDSPath)
	passwordFile := filepath.Join(registry[0], ReadOnlyAdminPasswordFile)

	password, err := os.ReadFile(passwordFile)
	if err != nil {
		return nil, err
	}

	client := http.Client{
		Timeout: 1 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", uds)
			},
		},
	}
	req, err := http.NewRequest(http.MethodGet, "http://unix/pool.xml", nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(ReadOnlyAdminUsername, string(password))

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}
