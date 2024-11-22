package main

import (
	"app/models"
	"app/service"
	"database/sql"
	"encoding/json"
	"fmt" 
	"net/http"
	"strconv"
	"io"
	_ "github.com/mattn/go-sqlite3"
	"app/repository"
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

	message := "Hello, " + request.URL.Query().Get("name")

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

	var requestModel models.Request

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&requestModel)

	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	err = service.SendMessage(requestModel)
	if err != nil {
		return
	}

	fprintf, err := fmt.Fprintf(response, "Received message: %+v, Fillname is %s", requestModel, requestModel.Message.User.GetFullName())
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

// ServeEcho выводит в консоль тело запроса
func ServeEcho(response http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	fmt.Println(string(body))

	message := "ok"

	_, err = fmt.Fprint(response, message)
	if err != nil {
		return
	}
}

func ServeDb(response http.ResponseWriter, request *http.Request) {
	db, err := sql.Open("sqlite3", "storage.db")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()

	rows, err := db.Query("SELECT * FROM user")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	users := []models.User{}
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

	// message := "ok"

	// _, err = fmt.Fprint(response, message)
	// if err != nil {
	// 	return
	// }
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
