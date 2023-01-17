package events_query

import (
	"database/sql"
	"github.com/arandich/telegram-dao/internal/database/entity"
	"log"
)

func Accept(id int, reward int, db *sql.DB, username string) bool {
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		return false
	}
	_, err = tx.Exec(`UPDATE events_journal set status = 'Подтверждена' WHERE id = $1`, id)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return false
	}
	_, err = tx.Exec(`UPDATE users set karma = karma + $1 WHERE username = $2`, reward, username)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return false
	}
	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return false
	} else {
		return true
	}
}

func JoinEventDb(eventId int, userId int, db *sql.DB) bool {
	result, err := db.Exec(`INSERT INTO events_journal (event_id, user_id) VALUES ($1,$2)`, eventId, userId)
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

func AddEvent(db *sql.DB, event *entity.Event) bool {
	result, err := db.Exec(`INSERT INTO event (name, date, reward) VALUES ($1,$2,$3)`, event.Name, event.Date, event.Reward)
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
