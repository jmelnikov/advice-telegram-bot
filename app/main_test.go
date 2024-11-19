package main

import (
	"testing"
)

func TestHelloNameEnglish(t *testing.T) {
	name := "Yuriy"
	language := "english"
	expected := "Hello, Yuriy"
	actual, err := HelloName(name, language)

	if err != nil {
		t.Error(err.Error())
		return
	}

	if expected != actual {
		t.Errorf("Тест провален! Got: %s, expected: %s.", actual, expected)
	}
}

func TestHelloNameRussian(t *testing.T) {
	name := "Yuriy"
	language := "russian"
	expected := "Привет, Yuriy"
	actual, err := HelloName(name, language)

	if err != nil {
		t.Error(err.Error())
		return
	}

	if expected != actual {
		t.Errorf("Тест провален! Got: %s, expected: %s.", actual, expected)
	}
}
