package main

import (
	"github.com/arandich/telegram-dao/internal/commands"
	"github.com/arandich/telegram-dao/internal/commands/events_commands"
	"github.com/arandich/telegram-dao/internal/commands/tokens"
	"github.com/arandich/telegram-dao/internal/config"
	"github.com/arandich/telegram-dao/internal/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
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

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.CallbackQuery != nil {

			log.Println(update.CallbackQuery.Data)

			s := strings.Split(update.CallbackQuery.Data, " ")
			command, data := s[0], s[1]
			switch command {
			case "Подтвердить":
				events_commands.AcceptEvent(&update, bot, db, data)
			case "Отклонить":
				events_commands.DenyEvent(&update, bot, db, data)
			}
		} else if update.Message.IsCommand() {
			user := commands.Check(&update, db)
			switch update.Message.Command() {
			case "send":
				tokens.SendTokensTo(&update, bot, db, user)
			case "events":
				if user.RoleId >= 3 {
					events_commands.EventsList(&update, bot, db)
				} else {
					commands.ErrorMsg(&update, bot, "Тебе не по силам вызвать эту команду")
				}
			case "start":
				commands.Start(&update, bot, user)
			default:
				commands.ErrorMsg(&update, bot, "Команда отсутствует ;(")
			}
		} else {
			user := commands.Check(&update, db)
			switch update.Message.Text {
			case "инфо":
				commands.Info(&update, bot, user)
			case "участники":
				if user.RoleId >= 2 {
					commands.GetAllUsers(&update, bot, db)
				} else {
					commands.ErrorMsg(&update, bot, "Тебе не по силам вызвать эту команду")
				}
			case "мои активности":
				commands.UserEvents(&update, bot, user, db)
			default:
				commands.ErrorMsg(&update, bot, "Я не понимаю ;(")
			}
		}

	}
}
