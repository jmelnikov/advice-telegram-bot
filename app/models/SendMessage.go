package models

type SendMessage struct {
	ChatId          int             `json:"chat_id"`
	Text            string          `json:"text"`
	ReplyParameters ReplyParameters `json:"reply_parameters"`
	ParseMode       string          `json:"parse_mode"`
}
