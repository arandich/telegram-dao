package main

import (
	"github.com/arandich/telegram-dao/internal/commands"
	"github.com/arandich/telegram-dao/internal/commands/events_commands"
	"github.com/arandich/telegram-dao/internal/commands/tokens"
	"github.com/arandich/telegram-dao/internal/commands/votes_commands"
	"github.com/arandich/telegram-dao/internal/commands/wallet"
	"github.com/arandich/telegram-dao/internal/config"
	"github.com/arandich/telegram-dao/internal/database"
	"github.com/arandich/telegram-dao/internal/database/entity"
	"github.com/arandich/telegram-dao/pkg/convert"
	"github.com/arandich/telegram-dao/pkg/response/callback"
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
	db.SetMaxOpenConns(10)
	db.SetMaxOpenConns(10)
	defer db.Close()

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
		if update.CallbackQuery != nil {
			log.Println("КОЛЛБЕК")
			s := strings.Split(update.CallbackQuery.Data, " ")
			var command string
			if len(s) == 3 || len(s) == 5 {
				command = s[0]
			} else {
				callback.CallBackReq(&update, bot, "Ошибка")
				callback.CallBackMsg(&update, bot, "Неверные аргументы команды")
			}

			switch command {
			case "Вариант1":
				userVote := entity.UserVote{
					VoteId: convert.ToInt(s[1]),
					UserId: convert.ToInt(s[2]),
					Amount: convert.ToInt(s[3]),
					Choice: s[4],
				}

				ok := votes_commands.Yes(&update, bot, db, &userVote)
				if !ok {
					callback.CallBackReq(&update, bot, "Ошибка")
					callback.CallBackMsg(&update, bot, "Ошибка голосования!")
				}
			case "Вариант2":
				userVote := entity.UserVote{
					VoteId: convert.ToInt(s[1]),
					UserId: convert.ToInt(s[2]),
					Amount: convert.ToInt(s[3]),
					Choice: s[4],
				}
				ok := votes_commands.No(&update, bot, db, &userVote)
				if !ok {
					callback.CallBackReq(&update, bot, "Ошибка")
					callback.CallBackMsg(&update, bot, "Ошибка голосования!")
				}
			case "Вариант3":
				userVote := entity.UserVote{
					VoteId: convert.ToInt(s[1]),
					UserId: convert.ToInt(s[2]),
					Amount: convert.ToInt(s[3]),
					Choice: s[4],
				}
				ok := votes_commands.Neutral(&update, bot, db, &userVote)
				if !ok {
					callback.CallBackReq(&update, bot, "Ошибка")
					callback.CallBackMsg(&update, bot, "Ошибка голосования!")
				}
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
			log.Println("КОМАНДА")
			user := commands.Check(&update, db)
			if user == nil {
				commands.Msg(&update, bot, "Ошибка")
				return
			}
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
			case "transactions":
				if user.RoleId >= 3 {
					commands.GetAllTransactions(&update, bot, db)
				} else {
					commands.ErrorMsg(&update, bot, "Тебе не по силам вызвать эту команду")
				}
			case "createVoting":
				if user.RoleId >= 3 {
					votes_commands.CreateVote(&update, bot, db)
				} else {
					commands.ErrorMsg(&update, bot, "Тебе не по силам вызвать эту команду")
				}
			case "wallet":

				wallet.AddWallet(&update, bot, db, user)
			case "start":
				commands.Start(&update, bot)
			default:
				commands.ErrorMsg(&update, bot, "Команда отсутствует ;(")
			}
		} else {
			log.Println("СООБЩЕНИЕ")
			user := commands.Check(&update, db)
			switch strings.ToLower(update.Message.Text) {
			case "инфо":
				commands.Info(&update, bot, user)
			case "голосования":
				votes_commands.AllVotes(&update, bot, db, user)
			case "мои голосования":
				votes_commands.UserVotes(&update, bot, user, db)
			case "активности":
				events_commands.EventsListUser(&update, bot, db, user)
			case "кошелек":
				if user.TonWallet == "" {
					commands.Msg(&update, bot, "У вас еще нет кошелька, чтобы добавить кошелек введите\n /wallet *ваш кошелек*")
				} else {
					commands.Msg(&update, bot, user.TonWallet)
				}
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
