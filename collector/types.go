// Copyright 2016 stuart nelson
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Copied from https://github.com/stuartnelson3/passenger_exporter/blob/80b16566cdab445f6e68f967019a95b67f608aca/structs.go

package collector

type Info struct {
	CapacityUsed            string       `xml:"capacity_used"`
	MaxProcessCount         string       `xml:"max"`
	PassengerVersion        string       `xml:"passenger_version"`
	AppCount                string       `xml:"group_count"`
	TopLevelRequestsInQueue string       `xml:"get_wait_list_size"`
	CurrentProcessCount     string       `xml:"process_count"`
	SuperGroups             []SuperGroup `xml:"supergroups>supergroup"`
}

type SuperGroup struct {
	RequestsInQueue string `xml:"get_wait_list_size"`
	CapacityUsed    string `xml:"capacity_used"`
	State           string `xml:"state"`
	Group           Group  `xml:"group"`
	Name            string `xml:"name"`
}

type Group struct {
	Environment           string    `xml:"environment"`
	DisabledProcessCount  string    `xml:"disabled_process_count"`
	UID                   string    `xml:"uid"`
	GetWaitListSize       string    `xml:"get_wait_list_size"`
	CapacityUsed          string    `xml:"capacity_used"`
	Name                  string    `xml:"name"`
	AppType               string    `xml:"app_type"`
	AppRoot               string    `xml:"app_root"`
	User                  string    `xml:"user"`
	ComponentName         string    `xml:"component_name"`
	LifeStatus            string    `xml:"life_status"`
	UUID                  string    `xml:"uuid"`
	Default               string    `xml:"default,attr"`
	DisablingProcessCount string    `xml:"disabling_process_count"`
	EnabledProcessCount   string    `xml:"enabled_process_count"`
	DisableWaitListSize   string    `xml:"disable_wait_list_size"`
	GID                   string    `xml:"gid"`
	ProcessesSpawning     string    `xml:"processes_being_spawned"`
	Options               Options   `xml:"options"`
	Processes             []Process `xml:"processes>process"`
}

type Process struct {
	CodeRevision        string `xml:"code_revision"`
	Enabled             string `xml:"enabled"`
	SpawnEndTime        string `xml:"spawn_end_time"`
	HasMetrics          string `xml:"has_metrics"`
	LifeStatus          string `xml:"life_status"`
	Busyness            string `xml:"busyness"`
	RealMemory          string `xml:"real_memory"`
	StickySessionID     string `xml:"sticky_session_id"`
	PSS                 string `xml:"pss"`
	Command             string `xml:"command"`
	LastUsed            string `xml:"last_used"`
	CPU                 string `xml:"cpu"`
	SpawnerCreationTime string `xml:"spawner_creation_time"`
	LastUsedDesc        string `xml:"last_used_desc"`
	Uptime              string `xml:"uptime"`
	Swap                string `xml:"swap"`
	Sessions            string `xml:"sessions"`
	RSS                 string `xml:"rss"`
	PrivateDirty        string `xml:"private_dirty"`
	RequestsProcessed   string `xml:"processed"`
	ProcessGroupID      string `xml:"process_group_id"`
	PID                 string `xml:"pid"`
	GUPID               string `xml:"gupid"`
	VMSize              string `xml:"vmsize"`
	Concurrency         string `xml:"concurrency"`
	SpawnStartTime      string `xml:"spawn_start_time"`
}

type Options struct {
	DefaultGroup              string `xml:"default_group"`
	RubyBinPath               string `xml:"ruby"`
	USTRouterAddress          string `xml:"ust_router_address"`
	USTRouterPassword         string `xml:"ust_router_password"`
	StartCommand              string `xml:"start_command"`
	USTRouterUsername         string `xml:"ust_router_username"`
	MaxPreloaderIdleTime      string `xml:"max_preloader_idle_time"`
	BaseURI                   string `xml:"base_uri"`
	SpawnMethod               string `xml:"spawn_method"`
	AppType                   string `xml:"app_type"`
	Environment               string `xml:"environment"`
	Analytics                 string `xml:"analytics"`
	MinProcesses              string `xml:"min_processes"`
	StartTimeout              string `xml:"start_timeout"`
	AppRoot                   string `xml:"app_root"`
	ProcessTitle              string `xml:"process_title"`
	Debugger                  string `xml:"debugger"`
	DefaultUser               string `xml:"default_user"`
	MaxOutOfBandWorkInstances string `xml:"max_out_of_band_work_instances"`
	MaxProcesses              string `xml:"max_processes"`
	AppGroupName              string `xml:"app_group_name"`
	StartupFile               string `xml:"startup_file"`
	IntegrationMode           string `xml:"integration_mode"`
	LogLevel                  string `xml:"log_level"`
}
