package main

import (
	"app/models"
	"app/service"
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", ServeBot)
	http.HandleFunc("/test", ServeTest)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}

func ServeTest(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		fprintf, err := fmt.Fprintf(response, "Method %s is not found!", request.Method)
		if err != nil {
			// Надо что-то сделать с этой переменной, в ней записано количество записанных байт
			fmt.Println(fprintf)
			return
		}
		return
	}

	message := "Hello, Yuriy"

	_, err := fmt.Fprint(response, message)
	if err != nil {
		return
	}
}

func ServeBot(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(response, "404", http.StatusBadRequest)
		return
	}
	err := service.SendMessage()
	if err != nil {
		return
	}

	var requestModel models.Request

	decoder := json.NewDecoder(request.Body)
	err = nil
	err = decoder.Decode(&requestModel)

	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	fprintf, err := fmt.Fprintf(response, "Received message: %+v, Fillname is %s", requestModel, requestModel.Message.From.GetFullName())
	if err != nil {
		// Надо что-то сделать с этой переменной, в ней записано количество записанных байт
		fmt.Println(fprintf)
		return
	}
}

// HelloName удалить потом, просто интересуюсь как работают тесты
func HelloName(name string, language string) (string, error) {
	if name == "" {
		name = "World"
	}

	prefix := ""

	switch language {
	case "english":
		prefix = "Hello"
	case "russian":
		prefix = "Привет"
	default:
		return "", fmt.Errorf("%s", "Не передан язык")
	}

	return prefix + ", " + name, nil
}
