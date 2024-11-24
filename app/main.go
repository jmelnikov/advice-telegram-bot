package main

import (
	"app/models"
	"app/service"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
)

func main() {
	fmt.Println("Запустились, слушаем запросы...")

	http.HandleFunc("/", ServeBot)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}

func ServeBot(response http.ResponseWriter, request *http.Request) {
	// Если метод запроса не POST, то возвращаем 405 Status Method Not Allowed
	if request.Method != http.MethodPost {
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Декодируем сообщение пользователя
	var requestModel models.Request
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&requestModel)

	// Проверяем не только на ошибки, но и на пустую структуру запроса
	if err != nil || requestModel == (models.Request{}) {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	// Запускаем основновной обработчик сообщения пользователя
	err = service.ProcessMessage(requestModel)

	// Ошибка обработки сообщения пользователя
	if err != nil {
		// Печатаем ошибку в консоль
		fmt.Println(err)

		// Отправляем сообщение об ошибке пользователю
		response.WriteHeader(http.StatusInternalServerError)

		return
	}

	// В случае успеха, возвращаем пустой ответ с кодом 204
	response.WriteHeader(http.StatusNoContent)
}
