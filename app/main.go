package main

import (
	"app/models"
	"encoding/json"
	"fmt"
	"net/http"
	// "app/service"
)

func main() {
	http.HandleFunc("/", ServeBot)
	http.ListenAndServe(":8080", nil)
}

func ServeBot(response http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.Error(response, "404", http.StatusBadRequest)
		return
	}

	var requestModel models.Request

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&requestModel)

	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(response, "Received message: %+v, Fillname is %s", requestModel, requestModel.Message.From.GetFullName())
}

/*
Этот метод удалить, просто интересуюсь как работают тесты
*/
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
