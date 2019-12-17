package models

type Group struct {
	Id          int    `json:"id"`
	CreatorId   int    `json:"creator_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type JsonGroup struct {
	Group   Group   `json:"group"`
	Message string  `json:"message"`
	Status  string  `json:"status"`
}

type JsonGroups struct {
	Groups  []Group   `json:"groups"`
	Message string  `json:"message"`
	Status  string  `json:"status"`
}