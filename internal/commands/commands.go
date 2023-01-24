package commands

import (
	"database/sql"
	"github.com/arandich/telegram-dao/internal/database/entity"
	"github.com/arandich/telegram-dao/internal/database/query"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
)

type Commands struct {
	list []string
}

func ErrorMsg(update *tgbotapi.Update, bot *tgbotapi.BotAPI, text string) {
	text = "! " + text + " !"
	res := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	res.ReplyToMessageID = update.Message.MessageID

	if _, err := bot.Send(res); err != nil {
		panic(err)
	}
}

func Msg(update *tgbotapi.Update, bot *tgbotapi.BotAPI, text string) {
	res := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	res.ReplyToMessageID = update.Message.MessageID

	if _, err := bot.Send(res); err != nil {
		panic(err)
	}
}

func MsgWithoutReply(update *tgbotapi.Update, bot *tgbotapi.BotAPI, text string) {
	res := tgbotapi.NewMessage(update.Message.Chat.ID, text)

	if _, err := bot.Send(res); err != nil {
		panic(err)
	}
}

func Check(update *tgbotapi.Update, db *sql.DB) *entity.User {
	user, ok := query.FindByUsername(update.Message.From.UserName, db)
	if !ok {
		log.Println("Юзер не найден")
		log.Println("Добавляем юзера в бд...")
		ok = adduser(update, db)
		if !ok {
			log.Println("Ошибка добавления")
			return nil
		}
		return Check(update, db)
	}

	return user
}

func Start(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	photo := tgbotapi.NewInputMediaPhoto(tgbotapi.FilePath("storage/images/logo.jpg"))
	mediaGroup := tgbotapi.NewMediaGroup(update.Message.Chat.ID, []interface{}{photo})
	var numericKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Инфо"),
			tgbotapi.NewKeyboardButton("Кошелек"),
			tgbotapi.NewKeyboardButton("Активности"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Голосования"),
			tgbotapi.NewKeyboardButton("Мои голосования"),
			tgbotapi.NewKeyboardButton("Мои активности"),
		),
	)
	res := tgbotapi.NewMessage(update.Message.Chat.ID, "Добро пожаловать! \nСписок доступных команд:")
	res.ReplyMarkup = numericKeyboard
	_, err := bot.SendMediaGroup(mediaGroup)
	if err != nil {
		panic(err)
	}

	if _, err = bot.Send(res); err != nil {
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

func GetAllTransactions(update *tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {
	trList, ok := query.SelectAllTransactions(db)
	if !ok {
		ErrorMsg(update, bot, "Ошибка запроса к журналу транзакицй")
	}

	text := "Список последних транзакций: \n"
	for _, val := range trList.List {
		text += "\n" + strconv.Itoa(val.TrId) + " - От: " + val.Sender + " - Кому: " + val.ToUsername + "\nКоличество токенов: " + strconv.Itoa(val.Amount) + "\nДата: " + val.TimeToString() + "\n"
	}

	res := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	res.ReplyToMessageID = update.Message.MessageID
	if _, err := bot.Send(res); err != nil {
		panic(err)
	}

}

func adduser(update *tgbotapi.Update, db *sql.DB) bool {
	ok := query.AddUser(update, db)
	if !ok {
		log.Println("Ошибка добавления юзера")
		return false
	}
	return true
}
