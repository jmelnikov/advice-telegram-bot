package main

import (
	"app/models"
	"app/service"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"os"
)

func main() {
	// Загружаем настройки исполнения
	loadConfig()

	fmt.Println("Запустились, слушаем запросы...")

	// Устанавливаем функцию ServeBot по маршруту /
	http.HandleFunc("/", ServeBot)

	// Начинаем слушать запросы на порту, указанном в файле config.json
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func loadConfig() {
	// Открвыаем файл config.json
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Printf("Ошибка при открытии конфигурационного файла: %v\n", err)
		return
	}

	// Читаем файл в переменную конфиг
	var config map[string]string
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Printf("Не удалось декодировать конфигурационный файл: %v", err)
		return
	}

	// Устанавливаем переменные окружения
	for key, value := range config {
		err := os.Setenv(key, value)
		if err != nil {
			fmt.Printf("Не удалось установить переменную окружения %s: %v", key, err)
			continue
		}
	}

	// Закрываем файл config.json
	err = file.Close()
	if err != nil {
		fmt.Printf("Не удалось закрыть файл: %v", err)
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
