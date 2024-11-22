package repository

import(
	"app/models"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func GetUserById(user_id int) (models.User, error) {
	db, err := sql.Open("sqlite3", "storage.db")
	if err != nil {
		return models.User{}, fmt.Errorf("Не удалось прочитать файл БД")
	}

	defer db.Close()

	row := db.QueryRow(fmt.Sprintf("SELECT * FROM user WHERE id=%d", user_id))

	user := models.User{}

	err = row.Scan(&user.Id, &user.IsBot, &user.FirstName, &user.LastName, &user.Username, &user.LanguageCode)
	if err != nil {
		return models.User{}, fmt.Errorf("Не удалось конвертировать пользователя в структуру или пользователь с id %d не найден", user_id)
	}

	return user, nil
}
