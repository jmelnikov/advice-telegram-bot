package models

import (
	"database/sql"
	"strings"
)

type GreatingDb struct {
	Id        int            `json:"id"`
	Text      sql.NullString `json:"text"`
	Gender    sql.NullString `json:"gender"`
	TimeOfDay sql.NullString `json:"time_of_day"`
}

func (greating GreatingDb) GetGreatingTextForUser(user UserDb) string {
	message := strings.Replace(greating.Text.String, "{{FIRST_NAME}}", user.FirstName, -1)

	return strings.TrimSpace(message)
}
