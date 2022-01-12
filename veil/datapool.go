package veil

import (
	"fmt"
	"net/http"
	"net/url"
)

const baseDataPoolUrl string = "/api/data-pools/"

type DataPoolService struct {
	client Client
}

type NameSharedStorage struct {
	Id          string `json:"id,omitempty"`
	VerboseName string `json:"verbose_name,omitempty"`
}

type NameLun struct {
	Id     string `json:"id,omitempty"`
	Device string `json:"device,omitempty"`
	Status string `json:"status,omitempty"`
}

type NameClusterStorage struct {
	Id          string `json:"id,omitempty"`
	VerboseName string `json:"verbose_name,omitempty"`
}

type NodesConnected struct {
	Id               string `json:"id,omitempty"`
	VerboseName      string `json:"verbose_name,omitempty"`
	ConnectionStatus string `json:"connection_status,omitempty"`
}

type DataPoolObjectsList struct {
	Id             string             `json:"id,omitempty"`
	Status         string             `json:"name,omitempty"`
	VerboseName    string             `json:"verbose_name,omitempty"`
	BuiltIn        bool               `json:"built_in,omitempty"`
	Priority       int                `json:"priority,omitempty"`
	FreeSpace      int                `json:"free_space,omitempty"`
	Size           int                `json:"size,omitempty"`
	UsedSpace      int                `json:"used_space,omitempty"`
	SharedStorage  NameSharedStorage  `json:"shared_storage,omitempty"`
	ClusterStorage NameClusterStorage `json:"cluster_storage,omitempty"`
	NodesConnected []NodesConnected   `json:"nodes_connected,omitempty"`
	IsoCount       int                `json:"iso_count,omitempty"`
	FileCount      int                `json:"file_count,omitempty"`
	Type           string             `json:"type,omitempty"`
	VdiskCount     int                `json:"vdisk_count,omitempty"`
	ZfsPool        string             `json:"zfs_pool,omitempty"`
	Tags           []Tags             `json:"tags,omitempty"`
	Hints          int                `json:"hints,omitempty"`
	ResourcePools  []NameResourcePool `json:"resource_pools,omitempty"`
}

type DataPoolObject struct {
	Id             string             `json:"id,omitempty"`
	VerboseName    string             `json:"verbose_name,omitempty"`
	Description    string             `json:"description,omitempty"`
	LockedBy       string             `json:"locked_by,omitempty"`
	BuiltIn        bool               `json:"built_in,omitempty"`
	EntityType     string             `json:"entity_type,omitempty"`
	Status         string             `json:"name,omitempty"`
	Created        string             `json:"created,omitempty"`
	Modified       string             `json:"modified,omitempty"`
	Type           string             `json:"type,omitempty"`
	Path           string             `json:"path,omitempty"`
	Priority       int                `json:"priority,omitempty"`
	FreeSpace      int                `json:"free_space,omitempty"`
	Size           int                `json:"size,omitempty"`
	UsedSpace      int                `json:"used_space,omitempty"`
	SharedStorage  NameSharedStorage  `json:"shared_storage,omitempty"`
	ClusterStorage NameClusterStorage `json:"cluster_storage,omitempty"`
	NodesConnected []NodesConnected   `json:"nodes_connected,omitempty"`
	Permissions    []string           `json:"permissions,omitempty"`
	Lun            NameLun            `json:"lun,omitempty"`
	ZfsPool        string             `json:"zfs_pool,omitempty"`
	Tags           []Tags             `json:"tags,omitempty"`
	Hints          int                `json:"hints,omitempty"`
}

type DataPoolsResponse struct {
	BaseListResponse
	Results []DataPoolObjectsList `json:"results,omitempty"`
}

func (dp *DataPoolObject) Refresh(client *WebClient) (*DataPoolObject, error) {
	_, err := client.ExecuteRequest("GET", fmt.Sprint(baseDataPoolUrl, dp.Id, "/"), []byte{}, dp)
	return dp, err
}

func (d *DataPoolService) List() (*DataPoolsResponse, *http.Response, error) {

	response := new(DataPoolsResponse)

	res, err := d.client.ExecuteRequest("GET", baseDataPoolUrl, []byte{}, response)

	return response, res, err
}

func (d *DataPoolService) ListParams(queryParams map[string]string) (*DataPoolsResponse, *http.Response, error) {
	listUrl := baseDataPoolUrl
	if len(queryParams) != 0 {
		params := url.Values{}
		for k, v := range queryParams {
			params.Add(k, v)
		}
		listUrl += "?"
		listUrl += params.Encode()
	}
	response := new(DataPoolsResponse)
	res, err := d.client.ExecuteRequest("GET", listUrl, []byte{}, response)
	return response, res, err
}

func (d *DataPoolService) Get(Id string) (*DataPoolObject, *http.Response, error) {

	dp := new(DataPoolObject)

	res, err := d.client.ExecuteRequest("GET", fmt.Sprint(baseDataPoolUrl, Id, "/"), []byte{}, dp)

	return dp, res, err
}
