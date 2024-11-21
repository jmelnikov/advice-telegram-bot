package service

import (
	"fmt"
	"io"
	"net/http"
)

func SendMessage() {
	request, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/test", nil)
	if err != nil {
		fmt.Println("Ошибка при создании запроса")
		fmt.Println(err)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса")
		fmt.Println(err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении тела ответа")
		fmt.Println(err)
	}

	fmt.Println("Ответ от сервера:")
	fmt.Println(string(body))

	response.Body.Close()
}
