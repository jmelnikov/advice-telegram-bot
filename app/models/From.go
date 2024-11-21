package models

import "strings"

type From struct {
	Id            int
	Is_bot        bool
	First_name    string
	Last_name     string
	Username      string
	Language_code string
}

func (f From) GetFullName() string {
	name := f.First_name + " " + f.Last_name
	return strings.TrimSpace(name)
}
