package veil

import (
	"net/http"
)

const SwaggerUrl string = "/api/swagger/?format=openapi"

type SwaggerService struct {
	client Client
}

type Contact struct {
	Email string `json:"email,omitempty"`
}

type Info struct {
	Title          string    `json:"title,omitempty"`
	Description    string    `json:"description,omitempty"`
	TermsOfService string    `json:"termsOfService,omitempty"`
	Contact        []Contact `json:"contact,omitempty"`
	Version        string    `json:"version,omitempty"`
}

type Swagger struct {
	Swagger  string   `json:"swagger,omitempty"`
	Info     Info     `json:"info,omitempty"`
	Host     string   `json:"host,omitempty"`
	Schemes  []string `json:"schemes,omitempty"`
	BasePath string   `json:"basePath,omitempty"`
	Consumes []string `json:"consumes,omitempty"`
	Produces []string `json:"produces,omitempty"`
	Paths    string   `json:"paths,omitempty"`
	// Definitions string   `json:"definitions,omitempty"`
}

func (d *SwaggerService) Get() (*Swagger, *http.Response, error) {

	response := new(Swagger)

	res, err := d.client.ExecuteRequest("GET", SwaggerUrl, []byte{}, response)

	return response, res, err
}
