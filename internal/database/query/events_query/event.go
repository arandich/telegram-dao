package events_query

import (
	"database/sql"
	"log"
)

func Accept(id int, db *sql.DB) bool {
	result, err := db.Exec(`UPDATE events_journal set status = 'Подтверждена' WHERE id = $1`, id)
	if err != nil {
		log.Println(err)
		return false
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return false
	}
	return true
}

func Deny(id int, db *sql.DB) bool {
	result, err := db.Exec(`UPDATE events_journal set status = 'Отклонена' WHERE id = $1`, id)
	if err != nil {
		log.Println(err)
		return false
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return false
	}
	return true
}
