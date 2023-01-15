package commands

import (
	"database/sql"
	"fmt"
	"github.com/arandich/telegram-dao/internal/database/entity"
	"github.com/arandich/telegram-dao/internal/database/query"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Commands struct {
	list []string
}

func GetList() *Commands {
	return &Commands{list: []string{
		"/start",
		"/инфо",
		"/участники",
		"/мои_активности",
	}}
}

func ContainsCommand(c *Commands, s string) (string, bool) {
	for _, v := range c.list {
		if s == v {
			return v, true
		}
	}
	return "", false
}

func ErrorMsg(update *tgbotapi.Update, bot *tgbotapi.BotAPI, text string) {
	res := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	res.ReplyToMessageID = update.Message.MessageID

	if _, err := bot.Send(res); err != nil {
		panic(err)
	}
}

func Check(update *tgbotapi.Update, db *sql.DB) *entity.User {
	user, ok := query.FindByUsername(update, db)
	if !ok {
		fmt.Println("Юзер не найден")
		fmt.Println("Добавляем юзера в бд...")
		adduser(update, db)
		return nil
	}

	return user
}

func Start(update *tgbotapi.Update, bot *tgbotapi.BotAPI, user *entity.User) {
	text := "Добро пожаловать, " + user.Username + "\nВ наше сообщество 'Bored Student Club'\n" +
		"Список доступных команд: \n" +
		"/check \n" +
		"/info \n" +
		"/test"
	res := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	res.ReplyToMessageID = update.Message.MessageID

	if _, err := bot.Send(res); err != nil {
		panic(err)
	}
}

func Info(update *tgbotapi.Update, bot *tgbotapi.BotAPI, user *entity.User) {

	text := "Информация об аккаунте: " + user.Username + "\n" +
		`Карма: ` + user.KarmaToString() + "\n" +
		"Токены: " + user.TokenToString() + "\n" +
		"Роль: " + entity.Roles.ListRoles[user.RoleId]
	res := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	res.ReplyToMessageID = update.Message.MessageID

	if _, err := bot.Send(res); err != nil {
		panic(err)
	}
}

func UserEvents(update *tgbotapi.Update, bot *tgbotapi.BotAPI, user *entity.User, db *sql.DB) {
	events, ok := query.SelectAllUserEvents(user, db)
	if !ok {
		ErrorMsg(update, bot, "У вас нет активностей ;(")
	}

	text := "Список твоих активностей: \n"
	for eventName, val := range events.List {
		text += "\n" + eventName + " - Статус: " + val.Status
	}

	res := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	res.ReplyToMessageID = update.Message.MessageID

	if _, err := bot.Send(res); err != nil {
		panic(err)
	}
}

func GetAllUsers(update *tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {
	users, ok := query.SelectAllUsers(db)
	if !ok {
		ErrorMsg(update, bot, "Пользователи не найдены")
	}

	text := "Список участников: \n"
	for username, val := range users.List {
		text += "\n" + username + " - Роль: " + entity.Roles.ListRoles[val.RoleId]
	}

	res := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	res.ReplyToMessageID = update.Message.MessageID

	if _, err := bot.Send(res); err != nil {
		panic(err)
	}

}

func adduser(update *tgbotapi.Update, db *sql.DB) {
	ok := query.AddUser(update, db)
	if !ok {
		fmt.Println("Ошибка добавления юзера")
		return
	}
}
