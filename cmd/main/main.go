package main

import (
	"github.com/arandich/telegram-dao/internal/commands"
	"github.com/arandich/telegram-dao/internal/config"
	"github.com/arandich/telegram-dao/internal/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	token, _ := os.LookupEnv("TOKEN")

	newConfig := config.GetConfig()

	db := database.ConnectDb(newConfig)

	comList := commands.GetList()

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		command, ok := commands.ContainsCommand(comList, update.Message.Text)
		if !ok {
			commands.ErrorMsg(&update, bot, "Команда отсутствует ;(")
		}

		switch {
		case command == "/инфо":
			user := commands.Check(&update, db)
			commands.Info(&update, bot, user)
		case command == "/start":
			user := commands.Check(&update, db)
			commands.Start(&update, bot, user)
		case command == "/участники":
			user := commands.Check(&update, db)
			if user.RoleId >= 2 {
				commands.GetAllUsers(&update, bot, db)
			} else {
				commands.ErrorMsg(&update, bot, "Тебе не по силам вызвать эту команду")
			}
		case command == "/мои_активности":
			user := commands.Check(&update, db)
			commands.UserEvents(&update, bot, user, db)
		}

	}
}
