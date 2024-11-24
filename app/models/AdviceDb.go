package models

import (
	"database/sql"
	"strings"
)

type AdviceDb struct {
	Id     int            `json:"id"`
	Text   sql.NullString `json:"text"`
	Gender sql.NullString `json:"gender"`
}

func (advice AdviceDb) GetAdviceTextForUser(user UserDb) string {
	message := strings.Replace(advice.Text.String, "{{FIRST_NAME}}", user.FirstName, -1)

	return strings.TrimSpace(message)
}
