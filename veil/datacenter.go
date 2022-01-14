package veil

import (
	"fmt"
	"net/http"
	"net/url"
)

const baseDataCenterUrl string = "/api/datacenters/"

type DataCenterService struct {
	client Client
}

type DataCenterObjectsList struct {
	ClustersCount          int    `json:"clusters_count,omitempty"`
	CpuCount               int    `json:"cpu_count,omitempty"`
	Id                     string `json:"id,omitempty"`
	MemoryCount            int    `json:"memory_count,omitempty"`
	Status                 string `json:"name,omitempty"`
	VerboseName            string `json:"verbose_name,omitempty"`
	SharedStoragesCount    int    `json:"shared_storages_count,omitempty"`
	TransportStoragesCount int    `json:"transport_storages_count,omitempty"`
	BuiltIn                bool   `json:"built_in,omitempty"`
	Tags                   []Tags `json:"tags,omitempty"`
	Hints                  int    `json:"hints,omitempty"`
}

type DataCenterObject struct {
	CpuCount        int           `json:"cpu_count,omitempty"`
	Created         string        `json:"created,omitempty"`
	Description     string        `json:"description,omitempty"`
	Id              string        `json:"id,omitempty"`
	LockedBy        string        `json:"locked_by,omitempty"`
	MemoryCount     int           `json:"memory_count,omitempty"`
	Clusters        []NameCluster `json:"clusters,omitempty"`
	Modified        string        `json:"modified,omitempty"`
	Permissions     []string      `json:"permissions,omitempty"`
	Status          string        `json:"name,omitempty"`
	VerboseName     string        `json:"verbose_name,omitempty"`
	BuiltIn         bool          `json:"built_in,omitempty"`
	Tags            []Tags        `json:"tags,omitempty"`
	Hints           int           `json:"hints,omitempty"`
	OptimalCpuModel string        `json:"optimal_cpu_model,omitempty"`
	EntityType      string        `json:"entity_type,omitempty"`
}

type DataCentersResponse struct {
	BaseListResponse
	Results []DataCenterObjectsList `json:"results,omitempty"`
}

func (d *DataCenterService) List() (*DataCentersResponse, *http.Response, error) {

	response := new(DataCentersResponse)

	res, err := d.client.ExecuteRequest("GET", baseDataCenterUrl, []byte{}, response)

	return response, res, err
}

func (d *DataCenterService) ListParams(queryParams map[string]string) (*DataCentersResponse, *http.Response, error) {
	listUrl := baseDataCenterUrl
	if len(queryParams) != 0 {
		params := url.Values{}
		for k, v := range queryParams {
			params.Add(k, v)
		}
		listUrl += "?"
		listUrl += params.Encode()
	}
	response := new(DataCentersResponse)
	res, err := d.client.ExecuteRequest("GET", listUrl, []byte{}, response)
	return response, res, err
}

func (d *DataCenterService) Get(Id string) (*DataCenterObject, *http.Response, error) {

	entity := new(DataCenterObject)

	res, err := d.client.ExecuteRequest("GET", fmt.Sprint(baseDataCenterUrl, Id, "/"), []byte{}, entity)

	return entity, res, err
}
