package veil

import (
	"fmt"
	"net/http"
	"net/url"
)

const baseVMachineInfUrl string = "/api/vmachine-infs/"

type VMachineInfService struct {
	client Client
}

type VmachineInfo struct {
	Id          string `json:"id,omitempty"`
	VerboseName string `json:"verbose_name,omitempty"`
}

type VnetworkInfo struct {
	Id          string `json:"id,omitempty"`
	VerboseName string `json:"verbose_name,omitempty"`
	EntityType  string `json:"entity_type,omitempty"`
}

type VMachineInfObjectsList struct {
	Id           string       `json:"id,omitempty"`
	Name         string       `json:"name,omitempty"`
	VmachineInfo VmachineInfo `json:"vmachine_info,omitempty"`
	MacAddress   string       `json:"mac_address,omitempty"`
	NicDriver    string       `json:"nic_driver,omitempty"`
	Status       string       `json:"status,omitempty"`
	VmachineName string       `json:"vmachine_name,omitempty"`
	VnetworkInfo VnetworkInfo `json:"vnetwork_info,omitempty"`
}

type VMachineInfObject struct {
	Id           string       `json:"id,omitempty"`
	Name         string       `json:"name,omitempty"`
	Description  string       `json:"description,omitempty"`
	LockedBy     string       `json:"locked_by,omitempty"`
	Created      string       `json:"created,omitempty"`
	Modified     string       `json:"modified,omitempty"`
	Permissions  []string     `json:"permissions,omitempty"`
	VmachineInfo VmachineInfo `json:"vmachine_info,omitempty"`
	NodeInfo     NameNode     `json:"node_info,omitempty"`
	MacAddress   string       `json:"mac_address,omitempty"`
	NicDriver    string       `json:"nic_driver,omitempty"`
	Status       string       `json:"status,omitempty"`
	EntityType   string       `json:"entity_type,omitempty"`
	VnetworkInfo VnetworkInfo `json:"vnetwork_info,omitempty"`
	LinkState    string       `json:"link_state,omitempty"`
}

type VMachinesResponse struct {
	BaseListResponse
	Results []VMachineInfObjectsList `json:"results,omitempty"`
}

func (entity *VMachineInfObject) Refresh(client *WebClient) (*VMachineInfObject, error) {
	_, err := client.ExecuteRequest("GET", fmt.Sprint(baseVMachineInfUrl, entity.Id, "/"), []byte{}, entity)
	return entity, err
}

func (d *VMachineInfService) List() (*VMachinesResponse, *http.Response, error) {
	response := new(VMachinesResponse)
	res, err := d.client.ExecuteRequest("GET", baseVMachineInfUrl, []byte{}, response)
	return response, res, err
}

func (d *VMachineInfService) ListParams(queryParams map[string]string) (*VMachinesResponse, *http.Response, error) {
	listUrl := baseVMachineInfUrl
	if len(queryParams) != 0 {
		params := url.Values{}
		for k, v := range queryParams {
			params.Add(k, v)
		}
		listUrl += "?"
		listUrl += params.Encode()
	}
	response := new(VMachinesResponse)
	res, err := d.client.ExecuteRequest("GET", listUrl, []byte{}, response)
	return response, res, err
}

func (d *VMachineInfService) Get(Id string) (*VMachineInfObject, *http.Response, error) {
	entity := new(VMachineInfObject)
	res, err := d.client.ExecuteRequest("GET", fmt.Sprint(baseVMachineInfUrl, Id, "/"), []byte{}, entity)
	return entity, res, err
}
