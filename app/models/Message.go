package models

type Message struct {
	MessageId  int		`json:"message_id"`
	User       User     `json:"from"`
	Chat       Chat     `json:"chat"`
	Date       int		`json:"date"`
	Text       string	`json:"text"`
}
