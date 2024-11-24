package models

type SendChatAction struct {
	ChatId          int             `json:"chat_id"`
	Action       string             `json:"action"`
}
