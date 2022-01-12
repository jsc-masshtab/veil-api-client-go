package veil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

const baseIsoUrl string = "/api/iso/"

type IsoService struct {
	client Client
}

type IsoObjectsList struct {
	Id       string           `json:"id,omitempty"`
	Status   string           `json:"name,omitempty"`
	FileName string           `json:"filename,omitempty"`
	Size     float64          `json:"size,omitempty"`
	DataPool NameTypeDataPool `json:"datapool,omitempty"`
	Domains  []NameDomain     `json:"domains,omitempty"`
	Created  string           `json:"created,omitempty"`
}

type IsoObject struct {
	Id          string           `json:"id,omitempty"`
	FileName    string           `json:"filename,omitempty"`
	Description string           `json:"description,omitempty"`
	LockedBy    string           `json:"locked_by,omitempty"`
	EntityType  string           `json:"entity_type,omitempty"`
	Status      string           `json:"name,omitempty"`
	Created     string           `json:"created,omitempty"`
	Modified    string           `json:"modified,omitempty"`
	DataPool    NameTypeDataPool `json:"datapool,omitempty"`
	Domains     []NameDomain     `json:"domains,omitempty"`
	Size        float64          `json:"size,omitempty"`
	Path        string           `json:"path,omitempty"`
	Permissions []string         `json:"permissions,omitempty"`
	UploadUrl   string           `json:"upload_url,omitempty"`
	DownloadUrl string           `json:"download_url,omitempty"`
}

type IsosResponse struct {
	BaseListResponse
	Results []IsoObjectsList `json:"results,omitempty"`
}

type IsoSoftAttach struct {
	Iso string `json:"iso"`
}

type IsoAttach struct {
	IsoSoftAttach
	Cdrom string `json:"cdrom"`
}

func (d *IsoService) List() (*IsosResponse, *http.Response, error) {

	response := new(IsosResponse)

	res, err := d.client.ExecuteRequest("GET", baseIsoUrl, []byte{}, response)

	return response, res, err
}

func (d *IsoService) ListParams(queryParams map[string]string) (*IsosResponse, *http.Response, error) {
	listUrl := baseIsoUrl
	if len(queryParams) != 0 {
		params := url.Values{}
		for k, v := range queryParams {
			params.Add(k, v)
		}
		listUrl += "?"
		listUrl += params.Encode()
	}
	response := new(IsosResponse)
	res, err := d.client.ExecuteRequest("GET", listUrl, []byte{}, response)
	return response, res, err
}

func (d *IsoService) Get(Id string) (*IsoObject, *http.Response, error) {

	node := new(IsoObject)

	res, err := d.client.ExecuteRequest("GET", fmt.Sprint(baseIsoUrl, Id, "/"), []byte{}, node)

	return node, res, err
}

func (d *IsoService) Create(DataPoolId string, FileName string) (*IsoObject, *http.Response, error) {
	// Part 1
	iso := new(IsoObject)

	body := struct {
		DataPoolId string `json:"datapool,omitempty"`
		FileName   string `json:"filename,omitempty"`
	}{DataPoolId, FileName}

	b, _ := json.Marshal(body)
	res, err := d.client.ExecuteRequest("PUT", baseIsoUrl, b, iso)
	if err != nil {
		return nil, res, err
	}

	// Part 2
	pwd, _ := os.Getwd()
	file, err := os.Open(pwd + "/file_data/" + FileName)
	if err != nil {
		return iso, res, err
	}
	defer file.Close()

	fileBody := &bytes.Buffer{}
	writer := multipart.NewWriter(fileBody)
	part, err := writer.CreateFormFile("file", filepath.Base(FileName))
	if err != nil {
		return nil, res, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, res, err
	}
	request, err := http.NewRequest("POST", fmt.Sprint(GetEnvUrl(), iso.UploadUrl), fileBody)
	err = writer.Close()
	if err != nil {
		return nil, res, err
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())
	response, err := d.client.Execute(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	return iso, response, err
}

func (d *IsoService) Download(iso *IsoObject) (*IsoObject, *http.Response, error) {
	// Get download_url
	res, err := d.client.ExecuteRequest("PUT", fmt.Sprint(baseIsoUrl, iso.Id, "/download/"), []byte{}, iso)
	if err != nil {
		return iso, res, err
	}
	// Create the file
	pwd, _ := os.Getwd()
	filePath := filepath.Dir(pwd) + "/file_data/downloaded_" + iso.FileName
	out, err := os.Create(filePath)
	defer out.Close()

	// Get the data
	resp, err := http.Get(fmt.Sprint(GetEnvUrl(), iso.DownloadUrl))
	if err != nil {
		return iso, res, fmt.Errorf("get the data error: %w", err)
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		errF := fmt.Errorf("bad status: %s", resp.Status)
		return iso, res, errF
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return iso, res, fmt.Errorf("write the body to file error: %w", err)
	}

	// Delete file
	err = os.Remove(filePath)
	if err != nil {
		return iso, res, fmt.Errorf("delete file error: %w", err)
	}
	return iso, res, nil
}
