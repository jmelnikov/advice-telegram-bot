package models

import "strings"

type User struct {
	Id            int		`json:"id"`
	IsBot         bool		`json:"is_bot"`
	FirstName     string	`json:"first_name"`
	LastName      string	`json:"last_name"`
	Username      string	`json:"username"`
	LanguageCode  string	`json:"language_code"`
}

func (u User) GetFullName() string {
	name := u.FirstName + " " + u.LastName
	return strings.TrimSpace(name)
}
