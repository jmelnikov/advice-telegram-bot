package models

type Message struct {
	Message_id int
	From       From
	Chat       Chat
	Date       int
	Text       string
}
