package models

import (
	"database/sql"
	"strings"
)

type UserDb struct {
	Id           int            `json:"id"`
	IsBot        bool           `json:"is_bot"`
	FirstName    string         `json:"first_name"`
	LastName     sql.NullString `json:"last_name,omitempty"`
	Username     sql.NullString `json:"username,omitempty"`
	LanguageCode string         `json:"language_code"`
	LastMessage  sql.NullInt64  `json:"-"`
	GreatingSent sql.NullBool   `json:"-"`
	Gender       sql.NullString `json:"-"`
}

func (u UserDb) GetUserDbFullName() string {
	name := u.FirstName + " " + u.LastName.String

	return strings.TrimSpace(name)
}
