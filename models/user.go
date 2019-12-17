package models

import "net/http"

type User struct {
	Id int            `json:"id"`
	Email string      `json:"email"`
	Login string      `json:"login"`
	Fullname string   `json:"fullname"`
	Password string   `json:"password"`
	AccVerified bool  `json:"acc_verified"`
}

type JsonUserBody struct {
	Message string `json:"message"`
	Status string `json:"status"`
	User User `json:"user"`
}

type JsonUser struct {
	Email string      `json:"email"`
	Login string      `json:"login"`
	Fullname string   `json:"fullname"`
}

type Tokens struct {
	Access *http.Cookie
	Refresh *http.Cookie
}