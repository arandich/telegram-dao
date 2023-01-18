package token

import (
	"database/sql"
	"log"
)

func UpdateUserToken(usernameTo string, usernameFrom string, amount int, db *sql.DB) bool {
	tx, err := db.Begin()
	defer db.Close()
	if err != nil {
		log.Println(err)
		return false
	}

	_, err = tx.Exec("INSERT INTO transaction_journal (sender, to_username, amount) VALUES ($1,$2,$3)", usernameFrom, usernameTo, amount)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return false
	}
	_, err = tx.Exec("UPDATE users set tokens = tokens + $1 WHERE username = $2", amount, usernameTo)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return false
	}

	_, err = tx.Exec("UPDATE users set tokens = tokens - $1 WHERE username = $2", amount, usernameFrom)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return false
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
