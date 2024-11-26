package service

import (
    "fmt"
    "os"
)

func GetSendMessageUrl() string {
    return fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", os.Getenv("BOT_API_KEY"))
}

func GetSendChatActionUrl() string {
    return fmt.Sprintf("https://api.telegram.org/bot%s/sendChatAction", os.Getenv("BOT_API_KEY"))
}
