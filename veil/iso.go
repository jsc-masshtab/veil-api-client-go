package veil

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

const baseIsoUrl = baseApiUrl + "iso/"

var IsoUrlUploadTimeout int64 = 360

type IsoService struct {
	client Client
}

type IsoObjectsList struct {
	Id       string           `json:"id,omitempty"`
	Status   string           `json:"status,omitempty"`
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
	Status      string           `json:"status,omitempty"`
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
	entity := new(IsoObject)
	res, err := d.client.ExecuteRequest("GET", fmt.Sprint(baseIsoUrl, Id, "/"), []byte{}, entity)
	return entity, res, err
}

func (d *IsoService) Create(DataPoolId string, FilenameUrl string, timeout int64) (*IsoObject, error) {
	if timeout == 0 {
		timeout = IsoUrlUploadTimeout
	}
	// Part 1
	entity := new(IsoObject)
	isUrl := isValidUrl(FilenameUrl)
	body := map[string]string{
		"datapool": DataPoolId,
	}
	if isUrl {
		body["url"] = FilenameUrl
	} else {
		body["filename"] = FilenameUrl
		if _, err := os.Stat("../file_data/" + FilenameUrl); errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("iso file does not exists (check folder file_data)")
		}
	}
	b, _ := json.Marshal(body)
	_, err := d.client.ExecuteRequest("PUT", baseIsoUrl, b, entity)
	if err != nil {
		return nil, err
	}

	// Part 2
	if isUrl {
		timeoutTime := time.Now().Unix() + timeout
		for true {
			_, err := d.client.ExecuteRequest("GET", fmt.Sprint(baseIsoUrl, entity.Id, "/"), []byte{}, entity)
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
			log.Fatal(err)
		}
	}

	return entity, err
}

func (d *IsoService) Download(entity *IsoObject) (*IsoObject, *http.Response, error) {
	// Get download_url
	res, err := d.client.ExecuteRequest("PUT", fmt.Sprint(baseIsoUrl, entity.Id, "/download/"), []byte{}, entity)
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

// Remove Эндпоинт удаления образа
func (d *IsoService) Remove(Id string) (bool, *http.Response, error) {
	res, err := d.client.ExecuteRequest("POST", fmt.Sprint(baseIsoUrl, Id, "/remove/"), []byte{}, nil)
	if err != nil {
		return false, res, err
	}
	return true, res, err
}
