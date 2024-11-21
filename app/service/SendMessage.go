package service

import (
	"fmt"
	"io"
	"net/http"
)

func SendMessage() error {
	request, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/test", nil)
	if err != nil {
		fmt.Println("Ошибка при создании запроса")
		_ = fmt.Errorf("%s", err)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса")
		_ = fmt.Errorf("%s", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении тела ответа")
		_ = fmt.Errorf("%s", err)
	}

	fmt.Println("Ответ от сервера:")
	fmt.Println(string(body))

	err = response.Body.Close()
	if err != nil {
		return err
	}

	return nil
}
