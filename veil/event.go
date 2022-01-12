package veil

import (
	"fmt"
	"net/http"
)

const baseEventUrl string = "/api/events/"

type EventService struct {
	client Client
}

type EventToEntity struct {
	EntityUuid  string `json:"entity_uuid,omitempty"`
	EntityClass string `json:"entity_class,omitempty"`
}

type EventObjectsList struct {
	Id            string          `json:"id,omitempty"`
	Message       string          `json:"message,omitempty"`
	DetailMessage string          `json:"detail_message,omitempty"`
	User          string          `json:"user,omitempty"`
	Type          string          `json:"type,omitempty"`
	Created       string          `json:"created,omitempty"`
	Task          string          `json:"task,omitempty"`
	Entities      []EventToEntity `json:"entities,omitempty"`
	Readed        []string        `json:"readed,omitempty"`
}

type EventObject struct {
	Id            string          `json:"id,omitempty"`
	Message       string          `json:"message,omitempty"`
	User          string          `json:"user,omitempty"`
	Created       string          `json:"created,omitempty"`
	Task          string          `json:"task,omitempty"`
	Readed        []string        `json:"readed,omitempty"`
	Entities      []EventToEntity `json:"entities,omitempty"`
	DetailMessage string          `json:"detail_message,omitempty"`
	Type          string          `json:"type,omitempty"`
	Permissions   []string        `json:"permissions,omitempty"`
}

type EventsResponse struct {
	BaseListResponse
	Results []EventObjectsList `json:"results,omitempty"`
}

func (d *EventService) List() (*EventsResponse, *http.Response, error) {

	response := new(EventsResponse)

	res, err := d.client.ExecuteRequest("GET", baseEventUrl, []byte{}, response)

	return response, res, err
}

func (d *EventService) Get(Id string) (*EventObject, *http.Response, error) {

	Event := new(EventObject)

	res, err := d.client.ExecuteRequest("GET", fmt.Sprint(baseEventUrl, Id, "/"), []byte{}, Event)

	return Event, res, err
}
