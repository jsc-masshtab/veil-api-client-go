package veil

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

const baseDomainUrl = baseApiUrl + "domains/"

type DomainService struct {
	client Client
}

type GuestUtils struct {
	Version    string   `json:"version,omitempty"`
	Hostname   string   `json:"hostname,omitempty"`
	Ipv4       []string `json:"ipv4,omitempty"`
	Interfaces []string `json:"interfaces,omitempty"`
	QemuState  bool     `json:"qemu_state,omitempty"`
}

type DomainObjectsList struct {
	Id                 string           `json:"id,omitempty"`
	VerboseName        string           `json:"verbose_name,omitempty"`
	MemoryCount        int              `json:"memory_count,omitempty"`
	Status             string           `json:"name,omitempty"`
	Parent             NameDomain       `json:"parent,omitempty"`
	CpuCount           int              `json:"cpu_count,omitempty"`
	MemoryPool         string           `json:"memory_pool,omitempty"`
	VmachineInfsCount  int              `json:"vmachine_infs_count,omitempty"`
	VdisksCount        int              `json:"vdisks_count,omitempty"`
	VfunctionsCount    int              `json:"vfunctions_count,omitempty"`
	LunsCount          int              `json:"luns_count,omitempty"`
	UserPowerState     int              `json:"user_power_state,omitempty"`
	Node               NameNode         `json:"node,omitempty"`
	Template           bool             `json:"template,omitempty"`
	MdevsCount         int              `json:"mdevs_count,omitempty"`
	Tags               []Tags           `json:"tags,omitempty"`
	Hints              int              `json:"hints,omitempty"`
	ResourcePool       NameResourcePool `json:"resource_pool,omitempty"`
	GuestUtils         GuestUtils       `json:"guest_utils,omitempty"`
	Thin               bool             `json:"thin,omitempty"`
	Replication        bool             `json:"replication,omitempty"`
	CpuUsedPercentUser string           `json:"cpu_used_percent_user,omitempty"`
	MemUsedPercentUser string           `json:"mem_used_percent_user,omitempty"`
	Priority           int              `json:"priority,omitempty"`
}

type DomainObject struct {
	Id                 string           `json:"id,omitempty"`
	VerboseName        string           `json:"verbose_name,omitempty"`
	Description        string           `json:"description,omitempty"`
	LockedBy           string           `json:"locked_by,omitempty"`
	Permissions        []string         `json:"permissions,omitempty"`
	Created            string           `json:"created,omitempty"`
	Modified           string           `json:"modified,omitempty"`
	MemoryCount        int              `json:"memory_count,omitempty"`
	Status             string           `json:"name,omitempty"`
	Parent             NameDomain       `json:"parent,omitempty"`
	CpuCount           int              `json:"cpu_count,omitempty"`
	MemoryPool         string           `json:"memory_pool,omitempty"`
	VmachineInfsCount  int              `json:"vmachine_infs_count,omitempty"`
	VdisksCount        int              `json:"vdisks_count,omitempty"`
	VfunctionsCount    int              `json:"vfunctions_count,omitempty"`
	LunsCount          int              `json:"luns_count,omitempty"`
	UserPowerState     int              `json:"user_power_state,omitempty"`
	Node               NameNode         `json:"node,omitempty"`
	Template           bool             `json:"template,omitempty"`
	MdevsCount         int              `json:"mdevs_count,omitempty"`
	Tags               []Tags           `json:"tags,omitempty"`
	Hints              int              `json:"hints,omitempty"`
	ResourcePool       NameResourcePool `json:"resource_pool,omitempty"`
	GuestUtils         GuestUtils       `json:"guest_utils,omitempty"`
	Thin               bool             `json:"thin,omitempty"`
	Replication        bool             `json:"replication,omitempty"`
	CpuUsedPercentUser string           `json:"cpu_used_percent_user,omitempty"`
	MemUsedPercentUser string           `json:"mem_used_percent_user,omitempty"`
	Priority           int              `json:"priority,omitempty"`
}

type DomainsResponse struct {
	BaseListResponse
	Results []DomainObjectsList `json:"results,omitempty"`
}

const MachineTypes = `(pc|q35)`
const CpuModes = `(default|host-model|host-passthrough|custom)`
const CleanTypes = `(zero|urandom)`

type SshInject struct {
	CreateUser bool   `json:"create_user,omitempty"`
	SshUser    string `json:"ssh_user,omitempty"`
	SshKey     string `json:"ssh_key,omitempty"`
}

type CloudConfig struct {
	UserData string `json:"user_data,omitempty"`
	MetaData string `json:"meta_data,omitempty"`
}

type CloudInitConf struct {
	CloudInitConfig CloudConfig `json:"cloud_init_config,omitempty"`
	CloudInit       bool        `json:"cloud_init,omitempty"`
}

type CpuTopology struct {
	CpuCount    int `json:"cpu_count,omitempty"`     // Group 1
	CpuCountMax int `json:"cpu_count_max,omitempty"` // Group 1

	CpuSockets int `json:"cpu_sockets,omitempty"` // Group 2
	CpuCores   int `json:"cpu_cores,omitempty"`   // Group 2
	CpuThreads int `json:"cpu_threads,omitempty"` // Group 2

	CpuMap map[string]string `json:"cpu_map,omitempty"` // Group 4

	CpuMode  string `json:"cpu_mode,omitempty"`  // Group 5
	CpuModel string `json:"cpu_model,omitempty"` // Group 5

	CpuPriority         int      `json:"cpu_priority,omitempty"`          // Group 6
	CpuShares           int      `json:"cpu_shares,omitempty"`            // Group 6
	CpuMinGuarantee     int      `json:"cpu_min_guarantee,omitempty"`     // Group 6
	CpuFeaturesRequired []string `json:"cpu_features_required,omitempty"` // Group 6
}

type DomainCreateConfig struct {
	IdempotencyKeyBase
	VerboseName  string `json:"verbose_name,omitempty"`
	DomainId     string `json:"domain_id,omitempty"`
	Description  string `json:"description,omitempty"`
	Node         string `json:"node,omitempty"`
	ResourcePool string `json:"resource_pool,omitempty"`
	MemoryCount  int    `json:"memory_count,omitempty"`
	BootType     string `json:"boot_type,omitempty"`
	CpuCount     int    `json:"cpu_count,omitempty"`
	CpuCountMax  int    `json:"cpu_count_max,omitempty"`
	CpuPriority  int    `json:"cpu_priority,omitempty"`
	CpuMode      string `json:"cpu_mode,omitempty"`
	CpuModel     string `json:"cpu_model,omitempty"`
	OsType       string `json:"os_type,omitempty"`
	OsVersion    string `json:"os_version,omitempty"`
	Machine      string `json:"machine,omitempty"`
}

type DomainMultiCreateConfig struct {
	DomainCreateConfig
	CloudInitConf
	Safety             bool          `json:"safety,omitempty"`
	StartOnBoot        bool          `json:"start_on_boot,omitempty"`
	CleanType          string        `json:"clean_type,omitempty"`
	CleanCount         int           `json:"clean_count,omitempty"`
	MemoryMinGuarantee int           `json:"memory_min_guarantee,omitempty"`
	MemoryShares       int           `json:"memory_shares,omitempty"`
	MemoryLimit        int           `json:"memory_limit,omitempty"`
	Vdisks             []VdiskAttach `json:"vdisks,omitempty"`
	Isos               []IsoAttach   `json:"isos,omitempty"`
	// luns
	// usb_devices
	// pci_devices
	// mdev_devices
	VmachineInfs []VMachineInfSoftCreate `json:"vmachine_infs,omitempty"`
	// vfunction_infs
	// cdroms
	// new_cdroms
	NewVdisks []VdiskCreateAttach `json:"new_vdisks,omitempty"`
	// new_luns
	NewIsos      []IsoSoftAttach `json:"new_isos,omitempty"`
	StartOn      bool            `json:"start_on,omitempty"`
	RemoteAccess bool            `json:"remote_access,omitempty"`
	Parent       string          `json:"parent,omitempty"`
	Thin         bool            `json:"thin,omitempty"`
	Clone        bool            `json:"clone,omitempty"`
	Template     string          `json:"template,omitempty"`
	CpuTopology  []CpuTopology   `json:"cpu_topology,omitempty"`
	SshInject    *SshInject      `json:"ssh_inject,omitempty"`
}

type DomainUpdateConfig struct {
	VerboseName string   `json:"verbose_name,omitempty"`
	Description string   `json:"description,omitempty"`
	SpiceStream bool     `json:"spice_stream,omitempty"`
	UserXml     string   `json:"user_xml,omitempty"`
	GpuOptimize bool     `json:"gpu_optimize,omitempty"`
	OsType      string   `json:"os_type,omitempty"`
	OsVersion   string   `json:"os_version,omitempty"`
	StartOnBoot bool     `json:"start_on_boot,omitempty"`
	Tablet      bool     `json:"tablet,omitempty"`
	Features    []string `json:"features,omitempty"`
	QemuArgs    []string `json:"qemu_args,omitempty"`
	Priority    int      `json:"priority,omitempty"`
}

type DomainCloneConfig struct {
	CloudInitConf
	Node         string `json:"node,omitempty"`
	ResourcePool string `json:"resource_pool,omitempty"`
	VerboseName  string `json:"verbose_name,omitempty"`
	DataPool     string `json:"datapool,omitempty"`
	// snapshot
	Count       int      `json:"count,omitempty"`
	DomainsIds  []string `json:"domains_ids,omitempty"`
	StartOn     bool     `json:"start_on,omitempty"`
	Template    bool     `json:"template"`
	Replication bool     `json:"replication,omitempty"`
}

func (entity *DomainObject) Refresh(client *WebClient) (*DomainObject, error) {
	_, err := client.ExecuteRequest("GET", fmt.Sprint(baseDomainUrl, entity.Id, "/"), []byte{}, entity)
	return entity, err
}

func (entity *DomainObject) WaitForGA(client *WebClient, timeout int64) (*DomainObject, error) {
	if timeout == 0 {
		timeout = 420
	}
	timeStart := time.Now().Unix()
	for true {
		_, err := client.ExecuteRequest("GET", fmt.Sprint(baseDomainUrl, entity.Id, "/"), []byte{}, entity)
		if err != nil {
			return entity, err
		}
		if entity.GuestUtils.QemuState == true {
			log.Printf("successfully waiting guest agent of domain %s", entity.VerboseName)
			return entity, nil
		}
		time.Sleep(time.Second * 5)
		timeNow := time.Now().Unix()
		if timeNow > timeStart+timeout {
			errMsg := fmt.Sprintf("waiting guest agent timeout error for domain %s.", entity.VerboseName)
			return entity, fmt.Errorf(errMsg)
		}
		log.Printf("waiting %d sec for guest agent of domain %s (max %dsec)", timeNow-timeStart, entity.VerboseName, timeout)
	}

	return entity, nil
}

func (d *DomainService) List() (*DomainsResponse, *http.Response, error) {
	response := new(DomainsResponse)
	res, err := d.client.ExecuteRequest("GET", baseDomainUrl, []byte{}, response)
	return response, res, err
}

func (d *DomainService) ListParams(queryParams map[string]string) (*DomainsResponse, *http.Response, error) {
	listUrl := baseDomainUrl
	if len(queryParams) != 0 {
		params := url.Values{}
		for k, v := range queryParams {
			params.Add(k, v)
		}
		listUrl += "?"
		listUrl += params.Encode()
	}
	response := new(DomainsResponse)
	res, err := d.client.ExecuteRequest("GET", listUrl, []byte{}, response)
	return response, res, err
}

func (d *DomainService) Create(config DomainCreateConfig) (*DomainObject, *http.Response, error) {
	domain := new(DomainObject)
	b, _ := json.Marshal(config)
	res, err := d.client.ExecuteRequest("POST", baseDomainUrl, b, domain)
	return domain, res, err
}

func (d *DomainService) MultiCreate(config DomainMultiCreateConfig) (*DomainObject, *http.Response, error) {
	domain := new(DomainObject)
	b, _ := json.Marshal(config)
	asyncResp := new(AsyncEntityResponse)
	res, err := d.client.ExecuteRequest("POST", fmt.Sprint(baseDomainUrl, "multi-create-domain/?async=1"), b, asyncResp)
	if err != nil {
		return domain, res, err
	}
	WaitTaskReady(d.client.RetClient(), asyncResp.Task.Id, true, 0, true)
	res, err = d.client.ExecuteRequest("GET", fmt.Sprint(baseDomainUrl, asyncResp.Entity, "/"), []byte{}, domain)
	return domain, res, err
}

func (d *DomainService) Get(Id string) (*DomainObject, *http.Response, error) {
	entity := new(DomainObject)
	res, err := d.client.ExecuteRequest("GET", fmt.Sprint(baseDomainUrl, Id, "/"), []byte{}, entity)
	return entity, res, err
}

func (d *DomainService) Update(Id string, config DomainUpdateConfig) (*DomainObject, *http.Response, error) {
	entity := new(DomainObject)
	b, _ := json.Marshal(config)
	res, err := d.client.ExecuteRequest("PUT", fmt.Sprint(baseDomainUrl, Id, "/"), b, entity)
	return entity, res, err
}

func (d *DomainService) Start(domain *DomainObject) (*DomainObject, *http.Response, error) {
	res, err := d.client.ExecuteRequest("POST", fmt.Sprint(baseDomainUrl, domain.Id, "/start/"), []byte{}, domain)
	return domain, res, err
}

func (d *DomainService) Suspend(domain *DomainObject) (*DomainObject, *http.Response, error) {
	res, err := d.client.ExecuteRequest("POST", fmt.Sprint(baseDomainUrl, domain.Id, "/suspend/"), []byte{}, domain)
	return domain, res, err
}

func (d *DomainService) Resume(domain *DomainObject) (*DomainObject, *http.Response, error) {
	res, err := d.client.ExecuteRequest("POST", fmt.Sprint(baseDomainUrl, domain.Id, "/resume/"), []byte{}, domain)
	return domain, res, err
}

func (d *DomainService) Shutdown(domain *DomainObject, force bool) (*DomainObject, *http.Response, error) {
	body := struct {
		Force bool `json:"force,omitempty"`
	}{force}
	b, _ := json.Marshal(body)
	res, err := d.client.ExecuteRequest("POST", fmt.Sprint(baseDomainUrl, domain.Id, "/shutdown/"), b, domain)
	return domain, res, err
}

func (d *DomainService) Reboot(domain *DomainObject, force bool) (*DomainObject, *http.Response, error) {
	body := struct {
		Force bool `json:"force,omitempty"`
	}{force}
	b, _ := json.Marshal(body)
	res, err := d.client.ExecuteRequest("POST", fmt.Sprint(baseDomainUrl, domain.Id, "/reboot/"), b, domain)
	return domain, res, err
}

func (d *DomainService) Template(domain *DomainObject, template bool) (*DomainObject, *http.Response, error) {
	body := struct {
		Template bool `json:"template"`
	}{template}
	b, _ := json.Marshal(body)
	res, err := d.client.ExecuteRequest("PUT", fmt.Sprint(baseDomainUrl, domain.Id, "/template/"), b, domain)
	return domain, res, err
}

func (d *DomainService) Clone(Id string, config DomainCloneConfig) (*DomainObject, *http.Response, error) {
	domain := new(DomainObject)
	b, _ := json.Marshal(config)
	asyncResp := new(AsyncEntityResponse)
	res, err := d.client.ExecuteRequest("POST", fmt.Sprint(baseDomainUrl, Id, "/clone/?async=1"), b, asyncResp)
	if err != nil {
		return domain, res, err
	}
	client := d.client.RetClient()
	taskObj := WaitTaskReady(client, asyncResp.Task.Id, true, 0, true)
	res, err = client.Task.Response(taskObj.Id, domain)
	return domain, res, err
}

func (d *DomainService) CloudInit(domain *DomainObject, config CloudInitConf) (*DomainObject, *http.Response, error) {
	b, _ := json.Marshal(config)
	res, err := d.client.ExecuteRequest("PUT", fmt.Sprint(baseDomainUrl, domain.Id, "/cloud-init/"), b, domain)
	return domain, res, err
}

func (d *DomainService) Remove(domainID string, full bool, force bool) (bool, *http.Response, error) {
	body := struct {
		Force bool `json:"force"`
		Full  bool `json:"full"`
	}{force, full}
	b, _ := json.Marshal(body)
	res, err := d.client.ExecuteRequest("POST", fmt.Sprint(baseDomainUrl, domainID, "/remove/"), b, nil)
	if err != nil {
		return false, res, err
	}
	return true, res, err
}
