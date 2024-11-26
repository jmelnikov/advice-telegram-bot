package service

import (
	"app/models"
	"database/sql"
	"fmt"
	"time"
)

func ProcessMessage(requestModel models.Request) error {
	// Получаем пользователя
	user, err := getUser(requestModel)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Передаём пользователя и его сообщение в функцию отправки приветствия
	err = GreatingsSendMessage(user, requestModel)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Передаём пользователя и его сообщение в функцию отправки совета
	err = AdviceSendMessage(user, requestModel)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Обновляем дату последнего сообщения пользователя
	user.LastMessage = sql.NullInt64{
		Int64: time.Now().Unix(),
		Valid: true,
	}

	// Обновляем данные пользователя в отдельном потоке
	go updateUser(user)

	return nil
}

func getUser(requestModel models.Request) (models.UserDb, error) {
	// Подключаемся к БД
	db, err := sql.Open("sqlite3", "storage.db")
	if err != nil {
		fmt.Println(err)
		return models.UserDb{}, err
	}

	// Ищем пользователя по его ID
	row := db.QueryRow("SELECT * FROM user WHERE id=?", requestModel.Message.User.Id)

	// Записываем выбранного пользователя в структуру
	user := models.UserDb{}
	err = row.Scan(&user.Id, &user.IsBot, &user.FirstName, &user.LastName, &user.Username, &user.LanguageCode, &user.LastMessage, &user.Gender)
	if err != nil {
		fmt.Println(err)
		// Если при заполнении пользователя прозошла ошибка, то создаём нового пользователя
		user, err = createUser(db, &requestModel)
		if err != nil {
			fmt.Println(err)
			return models.UserDb{}, err
		}
	}

	// Подставляем новое имя, фамилию и ник из запроса на случай, если пользователь поменял их
	user.FirstName = requestModel.Message.User.FirstName
	user.LastName = emptyStringToNull(requestModel.Message.User.LastName)
	user.Username = emptyStringToNull(requestModel.Message.User.Username)

	// Закрываем указатель на соединение с БД
	err = db.Close()
	if err != nil {
		fmt.Println(err)
		return models.UserDb{}, err
	}

	return user, err
}

func createUser(db *sql.DB, requestModel *models.Request) (models.UserDb, error) {
	// Подготавливаем запрпос
	prepare, err := db.Prepare("INSERT INTO user (id, is_bot, first_name, last_name, username, language_code, last_message, gender) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err)
		return models.UserDb{}, err
	}

	// Создаём структуру пользователя для БД
	user := models.UserDb{
		Id:           requestModel.Message.User.Id,
		IsBot:        requestModel.Message.User.IsBot,
		FirstName:    requestModel.Message.User.FirstName,
		LastName:     emptyStringToNull(requestModel.Message.User.LastName),
		Username:     emptyStringToNull(requestModel.Message.User.Username),
		LanguageCode: requestModel.Message.User.LanguageCode,
		LastMessage:  sql.NullInt64{},
		Gender:       sql.NullString{},
	}

	// Подставляем значения и выполняем запрос
	_, err = prepare.Exec(user.Id,
		user.IsBot,
		user.FirstName,
		user.LastName,
		user.Username,
		user.LanguageCode,
		user.LastMessage,
		user.Gender)
	if err != nil {
		fmt.Println(err)
		return models.UserDb{}, err
	}

	return user, err
}

func updateUser(user models.UserDb) {
	// Подключаемся к БД
	db, err := sql.Open("sqlite3", "storage.db")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Подготавливаем запрпос
	prepare, err := db.Prepare("UPDATE user SET first_name=?, last_name=?, username=?, last_message=? WHERE id=?")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Подставляем значения и выполняем запрос
	_, err = prepare.Exec(user.FirstName, user.LastName, user.Username, user.LastMessage, user.Id)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func emptyStringToNull(s string) sql.NullString {
	// Если длина строки нулевая, то возвращаем NullString
	if len(s) == 0 {
		return sql.NullString{}
	}

	// Заполняем NullString значением
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
