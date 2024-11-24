package models

type Chat struct {
	Id          int    `json:"id"`
	FirstName   string `json:"first_name"`
	MessageType string `json:"type"`
}
