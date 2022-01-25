package veil

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

const baseLibraryUrl = baseApiUrl + "library/"

var LibraryUrlUploadTimeout int64 = 360

type LibraryService struct {
	client Client
}

type LibraryObjectsList struct {
	Id       string           `json:"id,omitempty"`
	FileName string           `json:"filename,omitempty"`
	Status   string           `json:"status,omitempty"`
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
	Status     string           `json:"status,omitempty"`
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

type FileImportConfig struct {
	VerboseName  string `json:"verbose_name,omitempty"`
	WithDeletion bool   `json:"with_deletion,omitempty"`
	Preallocate  bool   `json:"preallocate,omitempty"`
	Datapool     string `json:"datapool,omitempty"`
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

func (d *LibraryService) Import(Id string, config FileImportConfig) (*VdiskObject, *http.Response, error) {
	entity := new(VdiskObject)
	b, _ := json.Marshal(config)
	asyncResp := new(AsyncEntityResponse)
	res, err := d.client.ExecuteRequest("POST", fmt.Sprint(baseLibraryUrl, Id, "/import-file/?async=1"), b, asyncResp)
	if err != nil {
		return entity, res, err
	}
	client := d.client.RetClient()
	taskObj := WaitTaskReady(client, asyncResp.Task.Id, true, 0, true)
	res, err = client.Task.Response(taskObj.Id, entity)
	return entity, res, err
}

func (d *LibraryService) Create(DataPoolId string, FilenameUrl string, timeout int64) (*LibraryObject, error) {
	if timeout == 0 {
		timeout = LibraryUrlUploadTimeout
	}
	// Part 1
	entity := new(LibraryObject)
	isUrl := isValidUrl(FilenameUrl)
	body := map[string]string{
		"datapool": DataPoolId,
	}
	if isUrl {
		body["url"] = FilenameUrl
	} else {
		body["filename"] = FilenameUrl
		if _, err := os.Stat("../file_data/" + FilenameUrl); errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("library file does not exists (check folder file_data)")
		}
	}
	b, _ := json.Marshal(body)
	_, err := d.client.ExecuteRequest("PUT", baseLibraryUrl, b, entity)
	if err != nil {
		return nil, err
	}

	// Part 2
	if isUrl {
		request, err := http.NewRequest("POST", fmt.Sprint(GetEnvUrl(), entity.UploadUrl), nil)
		if err != nil {
			return nil, err
		}
		_, err = d.client.Execute(request)
		if err != nil {
			return nil, err
		}
		timeoutTime := time.Now().Unix() + timeout
		for true {
			_, err := d.client.ExecuteRequest("GET", fmt.Sprint(baseLibraryUrl, entity.Id, "/"), []byte{}, entity)
			if entity.Status == Status.Active {
				return entity, err
			}
			if time.Now().Unix() > timeoutTime {
				return entity, fmt.Errorf("error uploading file by url: %w", err)
			}
			time.Sleep(time.Second * StatusCheckInterval)
		}
	} else {
		file, err := os.Open("../file_data/" + FilenameUrl)
		defer file.Close()
		if err != nil {
			return entity, err
		}
		fileBody := &bytes.Buffer{}
		writer := multipart.NewWriter(fileBody)
		part, err := writer.CreateFormFile("file", filepath.Base(FilenameUrl))
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(part, file)
		if err != nil {
			return nil, err
		}
		err = writer.Close()
		if err != nil {
			return nil, err
		}
		request, err := http.NewRequest("POST", fmt.Sprint(GetEnvUrl(), entity.UploadUrl), fileBody)
		request.Header.Add("Content-Type", writer.FormDataContentType())
		response, err := d.client.Execute(request)
		defer response.Body.Close()
		if err != nil {
			return nil, err
		}
	}
	return entity, nil
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

// Remove Эндпоинт удаления файла
func (d *LibraryService) Remove(Id string) (bool, *http.Response, error) {
	res, err := d.client.ExecuteRequest("POST", fmt.Sprint(baseLibraryUrl, Id, "/remove/"), []byte{}, nil)
	if err != nil {
		return false, res, err
	}
	return true, res, err
}
