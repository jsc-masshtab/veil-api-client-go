package veil

import (
	"fmt"
	"net/http"
	"net/url"
)

const baseNodeUrl string = "/api/nodes/"

type NodeService struct {
	client Client
}

type NodeObjectsList struct {
	Id                 string             `json:"id,omitempty"`
	VerboseName        string             `json:"verbose_name,omitempty"`
	CpuCount           int                `json:"cpu_count,omitempty"`
	MemoryCount        int                `json:"memory_count,omitempty"`
	Status             string             `json:"name,omitempty"`
	ManagementIp       string             `json:"management_ip,omitempty"`
	DomainsCount       int                `json:"domains_count,omitempty"`
	DomainsOnCount     int                `json:"domains_on_count,omitempty"`
	Cluster            NameCluster        `json:"cluster,omitempty"`
	BuiltIn            bool               `json:"built_in,omitempty"`
	Tags               []Tags             `json:"tags,omitempty"`
	Hints              int                `json:"hints,omitempty"`
	DatacenterName     string             `json:"datacenter_name,omitempty"`
	DatacenterId       string             `json:"datacenter_id,omitempty"`
	ResourcePools      []NameResourcePool `json:"resource_pools,omitempty"`
	CpuUsedPercentUser string             `json:"cpu_used_percent_user,omitempty"`
	MemUsedPercentUser string             `json:"mem_used_percent_user,omitempty"`
}

type NodeObject struct {
	Id           string   `json:"id,omitempty"`
	VerboseName  string   `json:"verbose_name,omitempty"`
	Description  string   `json:"description,omitempty"`
	LockedBy     string   `json:"locked_by,omitempty"`
	Permissions  []string `json:"permissions,omitempty"`
	Status       string   `json:"name,omitempty"`
	Created      string   `json:"created,omitempty"`
	Modified     string   `json:"modified,omitempty"`
	ManagementIp string   `json:"management_ip,omitempty"`
	BuiltIn      bool     `json:"built_in,omitempty"`
	MemoryCount  int      `json:"memory_count,omitempty"`
}

type NodesResponse struct {
	BaseListResponse
	Results []NodeObjectsList `json:"results,omitempty"`
}

func (d *NodeService) List() (*NodesResponse, *http.Response, error) {

	response := new(NodesResponse)

	res, err := d.client.ExecuteRequest("GET", baseNodeUrl, []byte{}, response)

	return response, res, err
}

func (d *NodeService) ListParams(queryParams map[string]string) (*NodesResponse, *http.Response, error) {
	listUrl := baseNodeUrl
	if len(queryParams) != 0 {
		params := url.Values{}
		for k, v := range queryParams {
			params.Add(k, v)
		}
		listUrl += "?"
		listUrl += params.Encode()
	}
	response := new(NodesResponse)
	res, err := d.client.ExecuteRequest("GET", listUrl, []byte{}, response)
	return response, res, err
}

func (d *NodeService) Get(Id string) (*NodeObject, *http.Response, error) {

	entity := new(NodeObject)

	res, err := d.client.ExecuteRequest("GET", fmt.Sprint(baseNodeUrl, Id, "/"), []byte{}, entity)

	return entity, res, err
}
