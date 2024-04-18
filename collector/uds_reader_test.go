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
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func TestRead_Files(t *testing.T) {
	tempWithPasswordFile, err := os.MkdirTemp(os.TempDir(), "")
	if err != nil {
		t.Errorf("failed to create temporary directory: %s", err.Error())
	}

	instRegWithPasswordFile := filepath.Join(tempWithPasswordFile, "passenger.aBc0d2z")
	err = os.Mkdir(instRegWithPasswordFile, 0755)
	if err != nil {
		t.Errorf("failed to create instance registry directory: %s", err.Error())
	}

	err = os.WriteFile(filepath.Join(instRegWithPasswordFile, ReadOnlyAdminPasswordFile), []byte("fake"), 0644)
	if err != nil {
		t.Errorf("failed to create password file: %s", err.Error())
	}

	tempWithoutPasswordFile, err := os.MkdirTemp(os.TempDir(), "")
	if err != nil {
		t.Errorf("failed to create temporary directory: %s", err.Error())
	}

	instRegWithoutPasswordFile := filepath.Join(tempWithoutPasswordFile, "passenger.zk9bVz7")
	err = os.Mkdir(instRegWithoutPasswordFile, 0755)
	if err != nil {
		t.Errorf("failed to create instance registry directory: %s", err.Error())
	}

	tests := []struct {
		path    string
		wantErr string
	}{
		{path: "", wantErr: `failed to detect Passenger instance registry directory`},
		{path: "/zzz", wantErr: `failed to detect Passenger instance registry directory`},
		{path: "/tmp", wantErr: `failed to detect Passenger instance registry directory`},
		{path: tempWithoutPasswordFile, wantErr: fmt.Sprintf(`open %s/read_only_admin_password.txt: no such file or directory`, instRegWithoutPasswordFile)},
		{path: tempWithPasswordFile, wantErr: fmt.Sprintf(`Get "http://unix/pool.xml": dial unix %s/agents.s/core_api: connect: no such file or directory`, instRegWithPasswordFile)},
	}

	for _, test := range tests {
		socketPath := filepath.Join(test.path, UDSPath)
		closeFn, _ := socketListerner(socketPath, []byte{})
		if closeFn != nil {
			defer closeFn()
		}

		reader := NewUDSReader(test.path)
		_, err = reader.Read()

		if err == nil {
			t.Errorf("expected error with %q, but got %q", test.path, err)
		}
		if err.Error() != test.wantErr {
			t.Errorf("expected error to match with %q, but got %q", err.Error(), test.wantErr)
		}
	}
}

func TestRead_Data(t *testing.T) {
	temp, err := os.MkdirTemp(os.TempDir(), "")
	if err != nil {
		t.Errorf("failed to create temporary directory: %s", err.Error())
	}
	defer os.RemoveAll(temp)
	fmt.Println(temp)

	instReg := filepath.Join(temp, "passenger.Lj9cMz7")
	err = os.Mkdir(instReg, 0755)
	if err != nil {
		t.Errorf("failed to create instance registry directory: %s", err.Error())
	}

	err = os.WriteFile(filepath.Join(instReg, ReadOnlyAdminPasswordFile), []byte("fake"), 0644)
	if err != nil {
		t.Errorf("failed to create password file: %s", err.Error())
	}

	err = os.Mkdir(filepath.Join(instReg, filepath.Dir(UDSPath)), 0755)
	if err != nil {
		t.Errorf("failed to create temporary directory: %s", err.Error())
	}

	fixture, _ := os.ReadFile("testdata/passenger_xml_output.xml")

	socketPath := filepath.Join(instReg, UDSPath)
	closeFn, err := socketListerner(socketPath, fixture)
	if err != nil {
		t.Fatalf("failed to start test server: %s", err.Error())
	}
	defer closeFn()

	reader := NewUDSReader(temp)
	resp, err := reader.Read()
	if err != nil {
		t.Errorf("failed to read data: %s", err.Error())
	}
	defer resp.Close()

	data, err := io.ReadAll(resp)
	if err != nil {
		t.Errorf("failed to read data: %s", err.Error())
	}

	if string(data) != string(fixture) {
		t.Errorf("read data different from fixture")
	}
}

func socketListerner(path string, body []byte) (func() error, error) {
	server := http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Write(body)
		}),
	}

	unixListener, err := net.Listen("unix", path)
	if err != nil {
		return nil, err
	}
	go func() {
		server.Serve(unixListener)
	}()
	return server.Close, nil
}
