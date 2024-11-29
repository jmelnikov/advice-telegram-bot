package service

import (
    "bytes"
    "fmt"
    "net/http"
    "os"
)

func getSendMessageUrl() string {
    return fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", os.Getenv("BOT_API_KEY"))
}

func getSendChatActionUrl() string {
    return fmt.Sprintf("https://api.telegram.org/bot%s/sendChatAction", os.Getenv("BOT_API_KEY"))
}

func sendRequest(encodedJson []byte, endpoint string) {
	request, err := http.NewRequest(
		http.MethodPost,
		endpoint,
		bytes.NewBuffer(encodedJson))
	request.Header.Set("Content-Type", "application/json")

	if err != nil {
		fmt.Println(err)
		return
	}

	// Отправляем подготовленный запрос
	client := &http.Client{}
	_, _ = client.Do(request)
}
