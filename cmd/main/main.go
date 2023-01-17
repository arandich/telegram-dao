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
	log.Println("Конфиг создан")

	updates := bot.GetUpdatesChan(updateConfig)
	log.Println("Получение апдейтов")
	for update := range updates {
		log.Println("Update найден")
		log.Println(update)
		if update.CallbackQuery != nil {
			log.Println("ЭТО КОЛЛБЕК")
			s := strings.Split(update.CallbackQuery.Data, " ")
			var command string

			if len(s) == 3 {
				command = s[0]
			} else {
				commands.ErrorMsg(&update, bot, "Неверные аргументы команды")
				continue
			}

			switch command {
			case "Участвовать":
				eventId, userId := s[1], s[2]
				events_commands.JoinEvent(&update, bot, db, eventId, userId)
			case "Подтвердить":
				eventId, reward := s[1], s[2]
				events_commands.AcceptEvent(&update, bot, db, eventId, reward)
			case "Отклонить":
				eventId := s[1]
				events_commands.DenyEvent(&update, bot, db, eventId)
			}
		} else if update.Message.IsCommand() {
			log.Println("ЭТО КОМАНДА")
			user := commands.Check(&update, db)
			switch update.Message.Command() {
			case "send":
				tokens.SendTokensTo(&update, bot, db, user)
			case "events":
				if user.RoleId >= 3 {
					events_commands.EventsListAdmin(&update, bot, db)
				} else {
					commands.ErrorMsg(&update, bot, "Тебе не по силам вызвать эту команду")
				}
			case "addEvent":
				if user.RoleId >= 3 {
					events_commands.CreateEvent(&update, bot, db)
				} else {
					commands.ErrorMsg(&update, bot, "Тебе не по силам вызвать эту команду")
				}
			case "start":
				commands.Start(&update, bot, user)
			default:
				commands.ErrorMsg(&update, bot, "Команда отсутствует ;(")
			}
		} else {
			log.Println("ЭТО СООБЩЕНИЕ")
			user := commands.Check(&update, db)
			switch update.Message.Text {
			case "инфо":
				commands.Info(&update, bot, user)
			case "активности":
				events_commands.EventsListUser(&update, bot, db, user)
			case "участники":
				if user.RoleId >= 2 {
					commands.GetAllUsers(&update, bot, db)
				} else {
					commands.ErrorMsg(&update, bot, "Тебе не по силам вызвать эту команду")
				}
			case "мои активности":
				events_commands.UserEvents(&update, bot, user, db)
			default:
				commands.ErrorMsg(&update, bot, "Я не понимаю ;(")
			}
		}

	}
}
