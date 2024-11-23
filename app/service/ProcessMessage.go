package service

import (
	"app/models"
	"database/sql"
)
import "fmt"

func ProcessMessage(requestModel models.Request) error {
	// Подключаемся к БД
	db, err := sql.Open("sqlite3", "storage.db")
	if err != nil {
		return err
	}

	row := db.QueryRow(fmt.Sprintf(
		"SELECT * FROM user WHERE id=%d", requestModel.Message.User.Id))
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Записываем выбранного пользователя в структуру
	user := models.User{}
	err = row.Scan(&user.Id, &user.IsBot, &user.FirstName, &user.LastName, &user.Username, &user.LanguageCode, &user.LastMessage, &user.GreatingSent)
	if err != nil {
		err := createUser(db, &requestModel)
		user = requestModel.Message.User
		if err != nil {
			return err
		}
	}

	// Закрываем указатель на соединение с БД
	err = db.Close()
	if err != nil {
		return err
	}

	fmt.Println(user)

	return fmt.Errorf("method is not implemented")
}

func createUser(db *sql.DB, requestModel *models.Request) error {
	prepare, err := db.Prepare("INSERT INTO user (id, is_bot, first_name, last_name, username, language_code, last_message, greating_sent) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	requestModel.Message.User.LastMessage = "00:00"
	requestModel.Message.User.GreatingSent = false

	_, err = prepare.Exec(requestModel.Message.User.Id,
		requestModel.Message.User.IsBot,
		requestModel.Message.User.FirstName,
		requestModel.Message.User.LastName,
		requestModel.Message.User.Username,
		requestModel.Message.User.LanguageCode,
		requestModel.Message.User.LastMessage,
		requestModel.Message.User.GreatingSent)
	if err != nil {
		return err
	}

	return nil
}
