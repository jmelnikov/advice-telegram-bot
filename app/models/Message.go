package models

type Message struct {
	MessageId  int		`json:"message_id"`
	From       From
	Chat       Chat
	Date       int		`json:"date"`
	Text       string	`json:"text"`
}
