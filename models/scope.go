package models

type Scope struct {
	Id            int   `json:"id"`
	CreatorId     int   `json:"creator_id"`
	GroupId       int   `json:"group_id"`
	BeginInterval int64 `json:"begin_interval"`
	EndInterval   int64 `json:"end_interval"`
}

type JsonScope struct {
	Message string `json:"message"`
	Status string `json:"status"`
	Scope Scope `json:"scope"`
}

type JsonScopes struct {
	Message string `json:"message"`
	Status string `json:"status"`
	Scopes []Scope `json:"scopes"`
}