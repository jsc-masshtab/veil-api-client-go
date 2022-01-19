package veil

type IdempotencyKeyBase struct {
	IdempotencyKey string `json:"idempotency_key,omitempty"`
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

// Base name structures

type NameDomain struct {
	Id          string `json:"id,omitempty"`
	VerboseName string `json:"verbose_name,omitempty"`
}

type NameNode struct {
	Id          string `json:"id,omitempty"`
	VerboseName string `json:"verbose_name,omitempty"`
}

type NameResourcePool struct {
	Id          string `json:"id,omitempty"`
	VerboseName string `json:"verbose_name,omitempty"`
}

type NameCluster struct {
	Id          string `json:"id,omitempty"`
	VerboseName string `json:"verbose_name,omitempty"`
}

type NameDatacenter struct {
	Id          string `json:"id,omitempty"`
	VerboseName string `json:"verbose_name,omitempty"`
}

// Name storages structures

type NameTypeDataPool struct {
	Id          string `json:"id,omitempty"`
	VerboseName string `json:"verbose_name,omitempty"`
	Type        string `json:"type,omitempty"`
}

type NameSharedStorage struct {
	Id          string `json:"id,omitempty"`
	VerboseName string `json:"verbose_name,omitempty"`
}

type NameLun struct {
	Id     string `json:"id,omitempty"`
	Device string `json:"device,omitempty"`
	Status string `json:"status,omitempty"`
}

type NameClusterStorage struct {
	Id          string `json:"id,omitempty"`
	VerboseName string `json:"verbose_name,omitempty"`
}
