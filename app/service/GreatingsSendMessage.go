package service

import (
	"app/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
	"unicode/utf8"
)

func GreatingsSendMessage(user models.UserDb, requestModel models.Request) error {
	// Если пользователь писал сообщения в течении предыдущих 8 часов, то ничего не пишем
	now := time.Now()
	lastMessage := time.Unix(user.LastMessage.Int64, 0)
	if now.Sub(lastMessage) < time.Hour*8 {
		return nil
	}

	// Получаем случайное приветствие
	greating, err := getRandomGreating(user)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Отправляем совет пользователю в новом потоке
	go sendGreating(user, requestModel, greating)

	return nil
}

func getRandomGreating(user models.UserDb) (models.GreatingDb, error) {
	// Подключаемся к БД
	db, err := sql.Open("sqlite3", "storage.db")
	if err != nil {
		fmt.Println(err)
		return models.GreatingDb{}, err
	}

	var timeOfDay string
	hour := time.Now().Hour()
	if hour >= 6 && hour < 11 {
		timeOfDay = "morning"
	} else if hour >= 11 && hour < 18 {
		timeOfDay = "day"
	} else if hour >= 18 && hour < 23 {
		timeOfDay = "evening"
	} else {
		timeOfDay = "night"
	}

	// Выбираем случайное приветствие из БД с учётом пола пользователя и времени суток
	row := db.QueryRow("SELECT * FROM greating WHERE (gender=? OR gender IS NULL) AND (time_of_day = ? OR time_of_day IS NULL) ORDER BY RANDOM()", user.Gender, timeOfDay)

	// Записываем приветствие в структуру
	var greating models.GreatingDb
	err = row.Scan(&greating.Id, &greating.Text, &greating.Gender, &greating.TimeOfDay)
	if err != nil {
		fmt.Println(err)
		return models.GreatingDb{}, err
	}

	// Закрываем соединение с БД
	err = db.Close()
	if err != nil {
		fmt.Println(err)
		return models.GreatingDb{}, err
	}

	return greating, nil
}

func sendGreating(user models.UserDb, requestModel models.Request, greating models.GreatingDb) {
	message := models.SendMessage{
		ChatId: requestModel.Message.Chat.Id,
		Text:   greating.GetGreatingTextForUser(user),
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
			fmt.Println(err)
			continue
		}

		// Отправляем уведомление "Печатает..." пользователю
		sendRequest(encodedChatAction, getSendChatActionUrl())

		// После каждой отправки засыпаем на 5 секунд
		time.Sleep(5 * time.Second)
	}

	// Кодируем сообщение в JSON
	encodedMessage, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Отправляем пользователю
	sendRequest(encodedMessage, getSendMessageUrl())
}
