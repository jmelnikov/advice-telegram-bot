package service

import (
	"app/models"
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"
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

	// Отправляем совет пользователю в новом потоке
	go sendAdvice(user, requestModel, advice)

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

func sendAdvice(user models.UserDb, requestModel models.Request, advice models.AdviceDb) {
	// Подготавливаем структуру сообщения для пользователя
	message := models.SendMessage{
		ChatId: requestModel.Message.Chat.Id,
		Text:   advice.GetAdviceTextForUser(user),
		ReplyParameters: models.ReplyParameters{
			MessageId: requestModel.Message.MessageId,
		},
		ParseMode: "html",
	}

	// Получаем количество секунд, нужное на набор сообщения
	// Рассчитываем, что средняя скорость печати -- 8 символов в секунду
	needSecondsForWriteMessage := utf8.RuneCountInString(message.Text) / 8

	// При отравке уведомления "Печатает...", он держится на стороне клиента 5 секунд.
	// Чтобы светилось сообщение "Печатает..." весь срок формирования сообщения,
	// Отправляем это уведомление каждые 5 секунд
	actionCount := needSecondsForWriteMessage / 5

	// Не отправлять сообщение моментально
	if actionCount < 1 {
		actionCount = 1
	}

	// Создаём уведомление о том, что бот печатает
	chatAction := models.SendChatAction{
		ChatId: requestModel.Message.Chat.Id,
		Action: "typing",
	}

	// В цикле отправляем сообщение столько раз, сколько высчитали выше
	for i := 0; i < actionCount; i++ {
		encodedChatAction, err := json.Marshal(chatAction)
		if err != nil {
			continue
		}

		// Отправляем уведомление "Печатает..." пользователю
		sendRequest(encodedChatAction, GetSendChatActionUrl())

		// После каждой отправки засыпаем на 5 секунд
		time.Sleep(5 * time.Second)
	}

	// Кодируем сообщение в JSON
	encodedMessage, err := json.Marshal(message)
	if err != nil {
		return
	}

	// Отправляем пользователю
	sendRequest(encodedMessage, GetSendMessageUrl())
}

func sendRequest(encodedJson []byte, endpoint string) {
	request, err := http.NewRequest(
		http.MethodPost,
		endpoint,
		bytes.NewBuffer(encodedJson))
	request.Header.Set("Content-Type", "application/json")

	if err != nil {
		return
	}

	// Отправляем подготовленный запрос
	client := &http.Client{}
	_, _ = client.Do(request)
}
