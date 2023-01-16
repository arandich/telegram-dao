package query

import (
	"database/sql"
	"github.com/arandich/telegram-dao/internal/database/entity"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func FindByUsername(username string, db *sql.DB) (*entity.User, bool) {
	rows, err := db.Query(`SELECT * FROM users WHERE username = $1`, username)
	if err != nil {
		return nil, false
	}

	defer rows.Close()
	user := entity.User{}

	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Username, &user.RoleId, &user.Karma, &user.Tokens, &user.CreatedAt)
		if err != nil {
			log.Println(err)
			continue
		}
	}

	if user.Username == "" {
		return nil, false
	}
	return &user, true
}

func SelectAllUsers(db *sql.DB) (*entity.AllUsers, bool) {
	rows, err := db.Query(`SELECT * FROM users`)
	if err != nil {
		return nil, false
	}

	defer rows.Close()
	users := entity.AllUsers{List: map[string]entity.User{}}

	for rows.Next() {
		user := entity.User{}
		err := rows.Scan(&user.Id, &user.Username, &user.RoleId, &user.Karma, &user.Tokens, &user.CreatedAt)
		if err != nil {
			log.Println(err)
			continue
		}

		users.List[user.Username] = user

	}

	if len(users.List) == 0 {
		return nil, false
	}
	return &users, true
}

func SelectAllUserEvents(user *entity.User, db *sql.DB) (*entity.EventsJournal, bool) {
	rows, err := db.Query(`SELECT event.id,event.name,users.username, event.date,events_journal.status FROM events_journal,event,users WHERE users.id = events_journal.user_id AND event.id = events_journal.event_id AND events_journal.user_id = $1;`, user.Id)
	if err != nil {
		log.Println(err)
		return nil, false
	}

	defer rows.Close()

	events := entity.EventsJournal{List: map[string]entity.UserEvent{}}

	for rows.Next() {
		event := entity.UserEvent{}
		err := rows.Scan(&event.Id, &event.Name, &event.UserName, &event.Date, &event.Status)
		if err != nil {
			log.Println(err)
			continue
		}

		events.List[event.Name] = event
	}

	if len(events.List) == 0 {
		return nil, false
	}

	return &events, true
}

func SelectAllUsersEvents(db *sql.DB) (*entity.EventsJournal, bool) {
	rows, err := db.Query(`SELECT events_journal.id,event.name,users.username, event.date,events_journal.status FROM events_journal,event,users WHERE users.id = events_journal.user_id AND event.id = events_journal.event_id AND events_journal.status = 'Ожидание';`)
	if err != nil {
		log.Println(err)
		return nil, false
	}

	defer rows.Close()

	events := entity.EventsJournal{List: map[string]entity.UserEvent{}}

	for rows.Next() {
		event := entity.UserEvent{}
		err := rows.Scan(&event.Id, &event.Name, &event.UserName, &event.Date, &event.Status)
		if err != nil {
			log.Println(err)
			continue
		}

		events.List[event.Name] = event
	}

	if len(events.List) == 0 {
		return nil, false
	}

	return &events, true
}

func AddUser(update *tgbotapi.Update, db *sql.DB) bool {
	result, err := db.Exec(`INSERT INTO users (username, karma, tokens,role_id) VALUES ($1,0,0,1)`, update.Message.From.UserName)
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
