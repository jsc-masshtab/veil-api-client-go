package veil

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func GetEnvToken() string {
	token := os.Getenv("VEIL_API_TOKEN")
	if token == "" {
		token := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjo1NiwidXNlcm5hbWUiOiJidXIiLCJleHAiOjE5NTU0Mjc5OTEsInNzbyI6ZmFsc2UsIm9yaWdfaWF0IjoxNjQwOTMxOTkxfQ.BCPJi1hE_uvlv_sCjLYwGGq2qKJU8dbR9UUC5Cy79AA"
		//return token, errors.New("Token is empty")
		return token
	}

	return token
}

func GetEnvUrl() string {
	url := os.Getenv("VEIL_API_URL")
	if url == "" {
		url := "http://192.168.11.105"
		//return url, errors.New("Url is empty")
		return url
	}

	return url
}

func IsSuccess(code int) bool {

	successCodes := []int{200, 202}

	for _, i := range successCodes {
		if i == code {
			return true
		}
	}
	return false
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func NameGenerator(class string) string {
	return fmt.Sprint(class, "_", StringWithCharset(3, charset), "__TEST")
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

type ErrorDict struct {
	Detail string `json:"detail,omitempty"`
	Code   string `json:"code,omitempty"`
	MsgKey string `json:"msg_key,omitempty"`
}

type ErrorResponse struct {
	Errors []ErrorDict `json:"errors,omitempty"`
}

type BaseListResponse struct {
	Count    int    `json:"count,omitempty"`
	Next     string `json:"next,omitempty"`
	Previous string `json:"previous,omitempty"`
}

type Tags struct {
	Colour      string `json:"colour,omitempty"`
	VerboseName string `json:"verbose_name,omitempty"`
}

type NameResourcePool struct {
	Id          string `json:"id,omitempty"`
	VerboseName string `json:"verbose_name,omitempty"`
}

type NameDomain struct {
	Id          string `json:"id,omitempty"`
	VerboseName string `json:"verbose_name,omitempty"`
}

type IdempotencyKeyBase struct {
	IdempotencyKey string `json:"idempotency_key,omitempty"`
}
