package repository

import (
	"database/sql"
	"github.com/ph4r5h4d/ask-spock/models"
	"time"
)

const dateTimeFormat = "2006-01-02 15:04:05.999999-07:00"

func GetOrCreateUser(db *sql.DB, username string) (models.User, error) {
	row := db.QueryRow("SELECT * FROM user WHERE username = ?", username)
	var cat, uat string
	user := models.User{}
	err := row.Scan(&user.Id, &user.Username, &user.Active, &cat, &uat)
	if err == sql.ErrNoRows {
		// create the user
		u := models.User{
			Username:  username,
			Active:    false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		stm, err := db.Prepare("INSERT INTO user(username, active, created_at, updated_at) VALUES (?,?,?,?)")
		if err != nil {
			return models.User{}, err
		}
		defer stm.Close()

		result, err := stm.Exec(u.Username, u.Active, u.CreatedAt, u.UpdatedAt)
		if err != nil {
			return models.User{}, err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return models.User{}, err
		}
		u.Id = id
		return u, nil
	}

	if err != nil {
		return models.User{}, err
	}

	user.CreatedAt, err = time.Parse(dateTimeFormat, cat)
	user.UpdatedAt, err = time.Parse(dateTimeFormat, uat)
	return user, nil
}
