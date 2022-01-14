package veil

import (
	"fmt"
	"net/http"
	"net/url"
)

const baseClusterUrl string = "/api/clusters/"

type ClusterService struct {
	client Client
}

type NameDatacenter struct {
	Id          string `json:"id,omitempty"`
	VerboseName string `json:"verbose_name,omitempty"`
}

type ClusterObjectsList struct {
	Id                 string         `json:"id,omitempty"`
	VerboseName        string         `json:"verbose_name,omitempty"`
	CpuCount           int            `json:"cpu_count,omitempty"`
	MemoryCount        int            `json:"memory_count,omitempty"`
	Status             string         `json:"name,omitempty"`
	Datacenter         NameDatacenter `json:"datacenter,omitempty"`
	NodesCount         int            `json:"nodes_count,omitempty"`
	BuiltIn            bool           `json:"built_in,omitempty"`
	Tags               []Tags         `json:"tags,omitempty"`
	Hints              int            `json:"hints,omitempty"`
	CpuUsedPercentUser string         `json:"cpu_used_percent_user,omitempty"`
	MemUsedPercentUser string         `json:"mem_used_percent_user,omitempty"`
}

type ClusterObject struct {
	CpuCount        int            `json:"cpu_count,omitempty"`
	Created         string         `json:"created,omitempty"`
	Datacenter      NameDatacenter `json:"datacenter,omitempty"`
	Description     string         `json:"description,omitempty"`
	OptimalCpuModel string         `json:"optimal_cpu_model,omitempty"`
	Id              string         `json:"id,omitempty"`
	LockedBy        string         `json:"locked_by,omitempty"`
	MemoryCount     int            `json:"memory_count,omitempty"`
	Modified        string         `json:"modified,omitempty"`
	Nodes           []NameNode     `json:"nodes,omitempty"`
	FencingType     string         `json:"fencing_type,omitempty"`
	HeartbeatType   string         `json:"heartbeat_type,omitempty"`
	Permissions     []string       `json:"permissions,omitempty"`
	Status          string         `json:"name,omitempty"`
	VerboseName     string         `json:"verbose_name,omitempty"`
	BuiltIn         bool           `json:"built_in,omitempty"`
	Tags            []Tags         `json:"tags,omitempty"`
	Hints           int            `json:"hints,omitempty"`
	EntityType      string         `json:"entity_type,omitempty"`

	HaAutoSelect bool `json:"ha_autoselect,omitempty"`
	HaEnabled    bool `json:"ha_enabled,omitempty"`
	// ha_nodepolicy
	HaRetryCount int `json:"ha_retrycount,omitempty"`
	HaTimeout    int `json:"ha_timeout,omitempty"`
	HaBootDelay  int `json:"ha_boot_delay,omitempty"`
	HaBootAgent  int `json:"ha_boot_agent,omitempty"`
	// drs
	// tag
	// quorum

}

type ClustersResponse struct {
	BaseListResponse
	Results []ClusterObjectsList `json:"results,omitempty"`
}

func (d *ClusterService) List() (*ClustersResponse, *http.Response, error) {

	response := new(ClustersResponse)

	res, err := d.client.ExecuteRequest("GET", baseClusterUrl, []byte{}, response)

	return response, res, err
}

func (d *ClusterService) ListParams(queryParams map[string]string) (*ClustersResponse, *http.Response, error) {
	listUrl := baseClusterUrl
	if len(queryParams) != 0 {
		params := url.Values{}
		for k, v := range queryParams {
			params.Add(k, v)
		}
		listUrl += "?"
		listUrl += params.Encode()
	}
	response := new(ClustersResponse)
	res, err := d.client.ExecuteRequest("GET", listUrl, []byte{}, response)
	return response, res, err
}

func (d *ClusterService) Get(Id string) (*ClusterObject, *http.Response, error) {

	entity := new(ClusterObject)

	res, err := d.client.ExecuteRequest("GET", fmt.Sprint(baseClusterUrl, Id, "/"), []byte{}, entity)

	return entity, res, err
}
