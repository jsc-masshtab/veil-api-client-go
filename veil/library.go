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

const baseLibraryUrl = baseApiUrl + "library/"

type LibraryService struct {
	client Client
}

type LibraryObjectsList struct {
	Id       string           `json:"id,omitempty"`
	FileName string           `json:"filename,omitempty"`
	Status   string           `json:"name,omitempty"`
	DataPool NameTypeDataPool `json:"datapool,omitempty"`
	Domain   NameDomain       `json:"domain,omitempty"`
	Size     float64          `json:"size,omitempty"`
	// additional_flags
	Created        string `json:"created,omitempty"`
	Path           string `json:"path,omitempty"`
	AssignmentType string `json:"assignment_type,omitempty"`
	InvalidHash    bool   `json:"invalid_hash,omitempty"`
}

type LibraryObject struct {
	Id         string           `json:"id,omitempty"`
	FileName   string           `json:"filename,omitempty"`
	LockedBy   string           `json:"locked_by,omitempty"`
	Status     string           `json:"name,omitempty"`
	Created    string           `json:"created,omitempty"`
	Modified   string           `json:"modified,omitempty"`
	EntityType string           `json:"entity_type,omitempty"`
	DataPool   NameTypeDataPool `json:"datapool,omitempty"`
	Domain     NameDomain       `json:"domain,omitempty"`
	Size       float64          `json:"size,omitempty"`
	Path       string           `json:"path,omitempty"`
	// additional_flags
	IntegrityHash  string   `json:"integrity_hash,omitempty"`
	Description    string   `json:"description,omitempty"`
	AssignmentType string   `json:"assignment_type,omitempty"`
	Compressed     bool     `json:"compressed,omitempty"`
	InvalidHash    bool     `json:"invalid_hash,omitempty"`
	Permissions    []string `json:"permissions,omitempty"`
	UploadUrl      string   `json:"upload_url,omitempty"`
	DownloadUrl    string   `json:"download_url,omitempty"`
}

type LibraryResponse struct {
	BaseListResponse
	Results []LibraryObjectsList `json:"results,omitempty"`
}

func (d *LibraryService) List() (*LibraryResponse, *http.Response, error) {
	response := new(LibraryResponse)
	res, err := d.client.ExecuteRequest("GET", baseLibraryUrl, []byte{}, response)
	return response, res, err
}

func (d *LibraryService) ListParams(queryParams map[string]string) (*LibraryResponse, *http.Response, error) {
	listUrl := baseLibraryUrl
	if len(queryParams) != 0 {
		params := url.Values{}
		for k, v := range queryParams {
			params.Add(k, v)
		}
		listUrl += "?"
		listUrl += params.Encode()
	}
	response := new(LibraryResponse)
	res, err := d.client.ExecuteRequest("GET", listUrl, []byte{}, response)
	return response, res, err
}

func (d *LibraryService) Get(Id string) (*LibraryObject, *http.Response, error) {
	entity := new(LibraryObject)
	res, err := d.client.ExecuteRequest("GET", fmt.Sprint(baseLibraryUrl, Id, "/"), []byte{}, entity)
	return entity, res, err
}

func (d *LibraryService) Create(DataPoolId string, FileName string) (*LibraryObject, *http.Response, error) {
	// Part 1
	entity := new(LibraryObject)

	body := struct {
		DataPoolId string `json:"datapool,omitempty"`
		FileName   string `json:"filename,omitempty"`
	}{DataPoolId, FileName}

	b, _ := json.Marshal(body)
	res, err := d.client.ExecuteRequest("PUT", baseLibraryUrl, b, entity)
	if err != nil {
		return nil, res, err
	}

	// Part 2
	pwd, _ := os.Getwd()
	file, err := os.Open(pwd + "/file_data/" + FileName)
	if err != nil {
		return entity, res, err
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
	request, err := http.NewRequest("POST", fmt.Sprint(GetEnvUrl(), entity.UploadUrl), fileBody)
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
	return entity, response, err
}

func (d *LibraryService) Download(entity *LibraryObject) (*LibraryObject, *http.Response, error) {
	// Get download_url
	res, err := d.client.ExecuteRequest("PUT", fmt.Sprint(baseLibraryUrl, entity.Id, "/download/"), []byte{}, entity)
	if err != nil {
		return entity, res, err
	}
	// Create the file
	pwd, _ := os.Getwd()
	filePath := filepath.Dir(pwd) + "/file_data/downloaded_" + entity.FileName
	out, err := os.Create(filePath)
	defer out.Close()

	// Get the data
	resp, err := http.Get(fmt.Sprint(GetEnvUrl(), entity.DownloadUrl))
	if err != nil {
		return entity, res, fmt.Errorf("get the data error: %w", err)
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		errF := fmt.Errorf("bad status: %s", resp.Status)
		return entity, res, errF
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return entity, res, fmt.Errorf("write the body to file error: %w", err)
	}

	// Delete file
	err = os.Remove(filePath)
	if err != nil {
		return entity, res, fmt.Errorf("delete file error: %w", err)
	}
	return entity, res, nil
}
