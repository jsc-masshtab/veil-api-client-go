package veil

import (
	"fmt"
	"net/http"
	"net/url"
)

const baseVnetUrl = baseApiUrl + "vnetworks/"

type VnetService struct {
	client Client
}

type LinkedLswitchInfo struct {
	Id             string           `json:"id,omitempty"`
	Name           string           `json:"name,omitempty"`
	ConnectedNodes []NodesConnected `json:"connected_nodes,omitempty"`
}

type Lswitch struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type ConnectedNodes struct {
	NodesConnected
	Endpoint   string `json:"endpoint,omitempty"`
	EndpointIp string `json:"endpoint_ip,omitempty"`
}

type LinkedVswitchInfo struct {
	Id              string `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	NodeId          string `json:"node_id,omitempty"`
	NodeVerboseName string `json:"node_verbose_name,omitempty"`
	UplinkState     string `json:"uplink_state,omitempty"`
}

type PortGroup struct {
	Id          string `json:"id,omitempty"`
	VerboseName string `json:"verbose_name,omitempty"`
	VlanMode    string `json:"vlan_mode,omitempty"`
	VlanTag     int    `json:"vlan_tag,omitempty"`
	// VlanTrunks []string `json:"vlan_trunks,omitempty"`
	Mtu                   int `json:"mtu,omitempty"`
	LinkedInterfacesCount int `json:"linked_interfaces_count,omitempty"`
}

type VnetObjectsList struct {
	Id          string `json:"id,omitempty"`
	Status      string `json:"name,omitempty"`
	VerboseName string `json:"verbose_name,omitempty"`
	Management  bool   `json:"management,omitempty"`
	Tags        []Tags `json:"tags,omitempty"`
	DataSubnet  string `json:"data_subnet,omitempty"`
	DataVlan    int    `json:"data_vlan,omitempty"`
	DataUseNat  bool   `json:"data_use_nat,omitempty"`
}

type VnetObject struct {
	Id                string              `json:"id,omitempty"`
	VerboseName       string              `json:"verbose_name,omitempty"`
	Description       string              `json:"description,omitempty"`
	LockedBy          string              `json:"locked_by,omitempty"`
	Created           string              `json:"created,omitempty"`
	Modified          string              `json:"modified,omitempty"`
	Permissions       []string            `json:"permissions,omitempty"`
	DataSubnet        string              `json:"data_subnet,omitempty"`
	DataVlan          int                 `json:"data_vlan,omitempty"`
	DataMtu           int                 `json:"data_mtu,omitempty"`
	DataNni           int                 `json:"data_vni,omitempty"`
	DataUseNat        bool                `json:"data_use_nat,omitempty"`
	LinkedLswitchInfo LinkedLswitchInfo   `json:"linked_lswitch_info,omitempty"`
	LinkedVswitchInfo []LinkedVswitchInfo `json:"linked_vswitch_info,omitempty"`
	// vnservices
	EntityType string    `json:"entity_type,omitempty"`
	Status     string    `json:"name,omitempty"`
	PortGroup  PortGroup `json:"port_group,omitempty"`
	// netflow_config
	Uplinks        []LinkedVswitchInfo `json:"uplinks,omitempty"`
	ConnectedNodes []ConnectedNodes    `json:"connected_nodes,omitempty"`
	Lswitch        Lswitch             `json:"lswitch,omitempty"`
	Tags           []Tags              `json:"tags,omitempty"`
	// vnservices_info
	Management bool `json:"management,omitempty"`
}

type VnetsResponse struct {
	BaseListResponse
	Results []VnetObjectsList `json:"results,omitempty"`
}

func (entity *VnetObject) Refresh(client *WebClient) (*VnetObject, error) {
	_, err := client.ExecuteRequest("GET", fmt.Sprint(baseVnetUrl, entity.Id, "/"), []byte{}, entity)
	return entity, err
}

func (d *VnetService) List() (*VnetsResponse, *http.Response, error) {
	response := new(VnetsResponse)
	res, err := d.client.ExecuteRequest("GET", baseVnetUrl, []byte{}, response)
	return response, res, err
}

func (d *VnetService) ListParams(queryParams map[string]string) (*VnetsResponse, *http.Response, error) {
	listUrl := baseVnetUrl
	if len(queryParams) != 0 {
		params := url.Values{}
		for k, v := range queryParams {
			params.Add(k, v)
		}
		listUrl += "?"
		listUrl += params.Encode()
	}
	response := new(VnetsResponse)
	res, err := d.client.ExecuteRequest("GET", listUrl, []byte{}, response)
	return response, res, err
}

func (d *VnetService) Get(Id string) (*VnetObject, *http.Response, error) {
	entity := new(VnetObject)
	res, err := d.client.ExecuteRequest("GET", fmt.Sprint(baseVnetUrl, Id, "/"), []byte{}, entity)
	return entity, res, err
}
