package service

import (
	"app/models"
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
)

func AdviceSendMessage(user models.UserDb, requestModel models.Request) error {
	// Если в сообщении пользователя не встречаются слова "совет" и "advice", то прерываем выполнение метода
	if strings.Contains(requestModel.Message.Text, "advice") == false &&
	   strings.Contains(requestModel.Message.Text, "совет") == false {
		return nil
	}

	// Получаем случайный совет
	advice, err := getRandomAdvice(user)
	if err != nil {
		return err
	}

	// Отправляем совет пользователю
	sendAdvice(user, requestModel, advice)

	return nil
}

func getRandomAdvice(user models.UserDb) (models.AdviceDb, error) {
	// Подключаемся к БД
	db, err := sql.Open("sqlite3", "storage.db")
	if err != nil {
		return models.AdviceDb{}, err
	}

	// Выбираем случайный совет из БД с учётом пола пользователя
	row := db.QueryRow("SELECT * FROM advice WHERE gender=? OR gender IS NULL ORDER BY RANDOM()", user.Gender)

	// Записываем совет в структуру
	var advice models.AdviceDb
	err = row.Scan(&advice.Id, &advice.Text, &advice.Gender)
	if err != nil {
		return models.AdviceDb{}, err
	}

	// Закрываем соединение с БД
	err = db.Close()
	if err != nil {
		return models.AdviceDb{}, err
	}

	return advice, nil
}

func sendAdvice(user models.UserDb, requestModel models.Request, advice models.AdviceDb)  {
	// Подготавливаем структуру сообщения для пользователя
	message := models.SendMessage{
		ChatId: requestModel.Message.Chat.Id,
		Text: advice.GetAdviceTextForUser(user),
		ReplyParameters: models.ReplyParameters{
			MessageId: requestModel.Message.MessageId,
		},
		ParseMode: "html",
	}

	encodedJson, err := json.Marshal(message)

	// Подготавляиваем запрос для отправки
	request, err := http.NewRequest(
		http.MethodPost,
		GetSendMessageUrl(),
		bytes.NewBuffer(encodedJson))
	request.Header.Set("Content-Type", "application/json")

	if err != nil {
		return
	}

	// Отправляем подготовленный запрос
	client := &http.Client{}
	_, _ = client.Do(request)
}
