package main

import (
	"app/models"
	"app/repository"
	"app/service"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"net/http"
	"strconv"
)

func main() {
	fmt.Println("Запустились, слушаем запросы...")

	http.HandleFunc("/", ServeBot)
	http.HandleFunc("/test", ServeTest)
	http.HandleFunc("/echo", ServeEcho)
	http.HandleFunc("/db", ServeDb)
	http.HandleFunc("/user", ServeUser)

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

	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
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

// TODO: всё, что ниже -- под удаление!
// Пока оставлю, чтобы смотреть на примеры.

func ServeTest(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		_, err := fmt.Fprintf(response, "Method %s is not found!", request.Method)
		if err != nil {
			return
		}
		return
	}

	message := "Hello, " + request.URL.Query().Get("name")

	_, err := fmt.Fprint(response, message)
	if err != nil {
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

// ServeEcho выводит в консоль тело запроса
func ServeEcho(response http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	fmt.Println("Echo request body:")
	fmt.Println(string(body))

	message := "ok"

	_, err = fmt.Fprint(response, message)
	if err != nil {
		return
	}
}

func ServeDb(response http.ResponseWriter, _ *http.Request) {
	db, err := sql.Open("sqlite3", "storage.db")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			return
		}
	}(db)

	rows, err := db.Query("SELECT * FROM user")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var users []models.User
	for rows.Next() {
		u := models.User{}
		err := rows.Scan(&u.Id, &u.IsBot, &u.FirstName, &u.LastName, &u.Username, &u.LanguageCode)
		if err != nil {
			fmt.Println(err)
			continue
		}

		users = append(users, u)
	}

	for _, user := range users {
		fmt.Println(fmt.Sprintf("Id: %d,\nFirst name: %s,\nLastname: %s,\nUsername: %s",
			user.Id, user.FirstName, user.LastName, user.Username))
	}

	response.WriteHeader(http.StatusNoContent)
}

func ServeUser(response http.ResponseWriter, request *http.Request) {
	userId, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil {
		fmt.Println(err)
		return
	}
	message := fmt.Sprintf("Hello, %d", userId)

	user, err := repository.GetUserById(userId)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(user)
	}

	_, err = fmt.Fprint(response, message)
}
