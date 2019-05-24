package models

type Users struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Fullname string `json:"fullname"`
}
