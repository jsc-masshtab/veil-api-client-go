package veil

import (
	"fmt"
	"net/http"
	"net/url"
)

const baseUserUrl string = "/api/users/"

type UserService struct {
	client Client
}

type Logins struct {
	CurrentFailedAttempts int    `json:"current_failed_attempts,omitempty"`
	TotalFailedAttempts   int    `json:"total_failed_attempts,omitempty"`
	LastFailedAttempt     string `json:"last_failed_attempt,omitempty"`
	LastSuccessAttempt    string `json:"last_success_attempt,omitempty"`
}

type Settings struct {
	Timezone               string `json:"timezone,omitempty"`
	ExpirationDate         string `json:"expiration_date,omitempty"`
	PasswordExpirationDate string `json:"password_expiration_date,omitempty"`
	DailyPeriod            bool   `json:"daily_period,omitempty"`
	DailyPeriodStart       string `json:"daily_period_start,omitempty"`
	DailyperiodEnd         string `json:"daily_period_end,omitempty"`
	SendErrors             bool   `json:"send_errors,omitempty"`
	TwoFaEnabled           bool   `json:"two_fa_enabled,omitempty"`
	MaxSessionsPerUser     int    `json:"max_sessions_per_user,omitempty"`
	InactivityDays         int    `json:"inactivity_days,omitempty"`
}

type Groups struct {
	Id             int    `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	TranslatedName string `json:"translated_name,omitempty"`
}

type UserObjectsList struct {
	Id       int      `json:"id,omitempty"`
	UserName string   `json:"username,omitempty"`
	Groups   []Groups `json:"groups,omitempty"`
	IsActive bool     `json:"is_active,omitempty"`
	Tags     []Tags   `json:"tags,omitempty"`
}

type UserObject struct {
	Id          int      `json:"id,omitempty"`
	UserName    string   `json:"username,omitempty"`
	Groups      []Groups `json:"groups,omitempty"`
	IsActive    bool     `json:"is_active,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	Email       string   `json:"email,omitempty"`
	Settings    Settings `json:"settings,omitempty"`
	FirstName   string   `json:"first_name,omitempty"`
	LastName    string   `json:"last_name,omitempty"`
	Logins      Logins   `json:"logins,omitempty"`
}

type UsersResponse struct {
	BaseListResponse
	Results []UserObjectsList `json:"results,omitempty"`
}

func (d *UserService) List() (*UsersResponse, *http.Response, error) {

	response := new(UsersResponse)

	res, err := d.client.ExecuteRequest("GET", baseUserUrl, []byte{}, response)

	return response, res, err
}

func (d *UserService) ListParams(queryParams map[string]string) (*UsersResponse, *http.Response, error) {
	listUrl := baseUserUrl
	if len(queryParams) != 0 {
		params := url.Values{}
		for k, v := range queryParams {
			params.Add(k, v)
		}
		listUrl += "?"
		listUrl += params.Encode()
	}
	response := new(UsersResponse)
	res, err := d.client.ExecuteRequest("GET", listUrl, []byte{}, response)
	return response, res, err
}

func (d *UserService) Get(Id int) (*UserObject, *http.Response, error) {

	user := new(UserObject)

	res, err := d.client.ExecuteRequest("GET", fmt.Sprint(baseUserUrl, Id, "/"), []byte{}, user)

	return user, res, err
}
