package service

import (
	"app/models"
	"database/sql"
)
import "fmt"

func ProcessMessage(requestModel models.Request) error {
	user, err := getUser(requestModel)
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("%+v", user))

	return fmt.Errorf("method is not implemented")
}

func getUser(requestModel models.Request) (models.UserDb, error) {
	// Подключаемся к БД
	db, err := sql.Open("sqlite3", "storage.db")
	if err != nil {
		return models.UserDb{}, err
	}

	row := db.QueryRow("SELECT * FROM user WHERE id=?", requestModel.Message.User.Id)
	if err != nil {
		return models.UserDb{}, err
	}

	// Записываем выбранного пользователя в структуру
	user := models.UserDb{}
	err = row.Scan(&user.Id, &user.IsBot, &user.FirstName, &user.LastName, &user.Username, &user.LanguageCode, &user.LastMessage, &user.GreatingSent)
	if err != nil {
		user, err = createUser(db, &requestModel)
		if err != nil {
			return models.UserDb{}, err
		}
	}

	// Закрываем указатель на соединение с БД
	err = db.Close()
	if err != nil {
		return models.UserDb{}, err
	}

	return user, err
}

func createUser(db *sql.DB, requestModel *models.Request) (models.UserDb, error) {
	prepare, err := db.Prepare("INSERT INTO user (id, is_bot, first_name, last_name, username, language_code, last_message, greating_sent) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return models.UserDb{}, err
	}

	user := models.UserDb{
		Id: requestModel.Message.User.Id,
		IsBot: requestModel.Message.User.IsBot,
		FirstName: requestModel.Message.User.FirstName,
		LastName: emptyStringToNull(requestModel.Message.User.LastName),
		Username: emptyStringToNull(requestModel.Message.User.Username),
		LanguageCode: requestModel.Message.User.LanguageCode,
		LastMessage: sql.NullInt64{},
		GreatingSent: sql.NullBool{Bool: false},
	}

	_, err = prepare.Exec(user.Id,
		user.IsBot,
		user.FirstName,
		user.LastName,
		user.Username,
		user.LanguageCode,
		user.LastMessage,
		user.GreatingSent)
	if err != nil {
		return models.UserDb{}, err
	}

	return user, err
}

func emptyStringToNull(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}

	return sql.NullString{
		String: s,
		Valid: true,
	}
}
