package models

type Request struct {
	UpdateId int `json:"update_id"`
	Message  Message
}
