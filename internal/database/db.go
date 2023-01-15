package database

import (
	"database/sql"
	"fmt"
	"github.com/arandich/telegram-dao/internal/config"
	_ "github.com/lib/pq"
	"log"
)

func ConnectDb(c *config.Config) *sql.DB {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Password, c.Dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		fmt.Println("Ошибка подключения к бд")
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		fmt.Println("Ошибка пинга")
		log.Fatal(err)
	}

	fmt.Println("Успешное подключение к бд")
	return db
}
