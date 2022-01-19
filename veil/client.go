package veil

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const baseApiUrl string = "/api/"

// Client Interface of Client for mocking data receiving in tests
type Client interface {
	ExecuteRequest(method, url string, body []byte, object interface{}) (*http.Response, error)
	Execute(req *http.Request) (*http.Response, error)
	RetClient() *WebClient
}

type WebClient struct {
	Token      string
	HTTPClient *http.Client

	BaseURL string

	// Services which is used for accessing API
	Domain      *DomainService
	Node        *NodeService
	Cluster     *ClusterService
	DataCenter  *DataCenterService
	DataPool    *DataPoolService
	Vdisk       *VdiskService
	Iso         *IsoService
	Task        *TaskService
	Event       *EventService
	User        *UserService
	Vnet        *VnetService
	VMachineInf *VMachineInfService
}

type Error struct {
	Error string `json:"error,omitempty"`
	Field string `json:"field,omitempty"`
}

// NewClient Web client creating
func NewClient(apiUrl string, token string, insecure bool) *WebClient {
	if apiUrl == "" {
		apiUrl = GetEnvUrl()
	}
	if token == "" {
		token = GetEnvToken()
	}
	tlsConf := &tls.Config{InsecureSkipVerify: insecure}
	tr := &http.Transport{
		TLSClientConfig:    tlsConf,
		DisableCompression: true,
		Proxy:              nil,
	}
	hClient := &http.Client{Transport: tr}
	client := &WebClient{
		Token:      token,
		HTTPClient: hClient,
		BaseURL:    apiUrl,
	}

	// Passing client to all services for easy client mocking in future and not passing it to every function
	client.Domain = &DomainService{client}
	client.Node = &NodeService{client}
	client.Cluster = &ClusterService{client}
	client.DataCenter = &DataCenterService{client}
	client.DataPool = &DataPoolService{client}
	client.Vdisk = &VdiskService{client}
	client.Iso = &IsoService{client}
	client.Task = &TaskService{client}
	client.Event = &EventService{client}
	client.User = &UserService{client}
	client.Vnet = &VnetService{client}
	client.VMachineInf = &VMachineInfService{client}
	return client
}

// ExecuteRequest Executing HTTP Request (receiving info from API)
func (client *WebClient) ExecuteRequest(method string, url string, body []byte, object interface{}) (*http.Response, error) {
	req, err := http.NewRequest(method, fmt.Sprint(client.BaseURL, url), bytes.NewBuffer(body))
	if err != nil {
		return new(http.Response), err
	}

	req.Header.Set("Authorization", "jwt "+client.Token)
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("Accept-Language", "en")
	res, err := client.HTTPClient.Do(req)
	if err != nil {
		return res, err
	}
	defer res.Body.Close()

	// Cloning response body for future using
	buf, _ := ioutil.ReadAll(res.Body)
	reader := ioutil.NopCloser(bytes.NewBuffer(buf))
	res.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

	if !IsSuccess(res.StatusCode) {
		response := new(ErrorResponse)
		err := json.NewDecoder(reader).Decode(response)
		if err != nil && err != io.EOF {
			log.Println(err)
			return res, err
		} else {
			errMsg := fmt.Sprintf("status code: %d, detail: %s on url %s %s", res.StatusCode, res.Body, method, url)
			return res, errors.New(errMsg)
		}

	}
	if object != nil && (res.StatusCode == 200 || res.StatusCode == 202) {
		err := json.NewDecoder(reader).Decode(object)

		// EOF means empty response body, this error is not needed
		if err != nil && err != io.EOF {
			log.Println(err)
			return res, err
		}
	}

	return res, nil
}

// Execute user HTTP Request
func (client *WebClient) Execute(req *http.Request) (*http.Response, error) {
	return client.HTTPClient.Do(req)
}

func (client *WebClient) RetClient() *WebClient {
	return client
}
