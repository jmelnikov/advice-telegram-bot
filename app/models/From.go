package models

import "strings"

type From struct {
	Id            int		`json:"id"`
	IsBot         bool		`json:"is_bot"`
	FirstName     string	`json:"first_name"`
	LastName      string	`json:"last_name"`
	Username      string	`json:"username"`
	LanguageCode  string	`json:"language_code"`
}

func (f From) GetFullName() string {
	name := f.FirstName + " " + f.LastName
	return strings.TrimSpace(name)
}
