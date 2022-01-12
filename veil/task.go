package veil

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const baseTaskUrl string = "/api/tasks/"

// TaskStatusCheckInterval - time between async checks in seconds
const TaskStatusCheckInterval = 1

type TaskService struct {
	client Client
}

type TaskStatusStruct struct {
	InProgress, Success, Failed, Canceled, Lost, Partial string
}

var TaskStatus = TaskStatusStruct{
	InProgress: "IN_PROGRESS",
	Success:    "SUCCESS",
	Failed:     "FAILED",
	Canceled:   "CANCELED",
	Lost:       "LOST",
	Partial:    "PARTIAL",
}

type TaskUser struct {
	Id       int    `json:"id,omitempty"`
	UserName string `json:"username,omitempty"`
}

type NodesUserResponses struct {
	NodeId       string `json:"node_id,omitempty"`
	NodeName     string `json:"node_name,omitempty"`
	NodeResponse string `json:"node_response,omitempty"`
}

type TaskObjectsList struct {
	Id                 string               `json:"id,omitempty"`
	Progress           int                  `json:"progress,omitempty"`
	Status             string               `json:"status,omitempty"`
	Name               string               `json:"name,omitempty"`
	Created            string               `json:"created,omitempty"`
	Executed           string               `json:"executed,omitempty"`
	NodesUserResponses []NodesUserResponses `json:"nodes_user_responses,omitempty"`
	IsMultitask        bool                 `json:"is_multitask,omitempty"`
	User               TaskUser             `json:"user,omitempty"`
	Parent             string               `json:"parent,omitempty"`
	ErrorMessage       string               `json:"error_message,omitempty"`
	IsCancellable      bool                 `json:"is_cancellable,omitempty"`
	FinishedTime       string               `json:"finished_time,omitempty"`
}

type TaskObject struct {
	Id                 string               `json:"id,omitempty"`
	Progress           int                  `json:"progress,omitempty"`
	Status             string               `json:"status,omitempty"`
	Name               string               `json:"name,omitempty"`
	VerboseName        string               `json:"verbose_name,omitempty"`
	Created            string               `json:"created,omitempty"`
	Executed           string               `json:"executed,omitempty"`
	FinishedTime       string               `json:"finished_time,omitempty"`
	NodesUserResponses []NodesUserResponses `json:"nodes_user_responses,omitempty"`
	IsMultitask        bool                 `json:"is_multitask,omitempty"`
	User               TaskUser             `json:"user,omitempty"`
	Parent             string               `json:"parent,omitempty"`
	ErrorMessage       string               `json:"error_message,omitempty"`
	IsCancellable      bool                 `json:"is_cancellable,omitempty"`
	Permissions        []string             `json:"permissions,omitempty"`
}

type TasksResponse struct {
	BaseListResponse
	Results []TaskObjectsList `json:"results,omitempty"`
}

func (d *TaskService) List() (*TasksResponse, *http.Response, error) {

	response := new(TasksResponse)

	res, err := d.client.ExecuteRequest("GET", baseTaskUrl, []byte{}, response)

	return response, res, err
}

func (d *TaskService) Get(Id string) (*TaskObject, *http.Response, error) {

	task := new(TaskObject)

	res, err := d.client.ExecuteRequest("GET", fmt.Sprint(baseTaskUrl, Id, "/"), []byte{}, task)

	return task, res, err
}

func WaitTaskReady(uuid string, blocked bool, timeout int, panicTimeout bool) *TaskObject {
	if timeout == 0 {
		timeout = 180
	}
	client := NewClient("", "", false)
	task, _, _ := client.Task.Get(uuid)
	if task.Status != TaskStatus.InProgress {
		return task
	} else if blocked {
		timeoutTime := time.Now().Second() + timeout
		for true {
			task, _, _ := client.Task.Get(uuid)
			if task.Status != TaskStatus.InProgress {
				time.Sleep(time.Second * TaskStatusCheckInterval)
				task, _, _ := client.Task.Get(uuid)
				return task
			}
			if time.Now().Second() > timeoutTime {
				if panicTimeout {
					errMsg := fmt.Sprintf("Task: %s wait %d timeout error. is_multitask: %s. progress: %d. status: %s", task.Name, timeout, strconv.FormatBool(task.IsMultitask), task.Progress, task.Status)
					panic(errMsg)
				}
			}
		}
	}
	//time.Sleep(time.Second)
	return task
}

type AsyncResponse struct {
	Entity string     `json:"entity,omitempty"`
	Task   TaskObject `json:"_task,omitempty"`
}
