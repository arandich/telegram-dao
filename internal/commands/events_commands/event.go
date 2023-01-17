package events_commands

import (
	"database/sql"
	"github.com/arandich/telegram-dao/internal/commands"
	"github.com/arandich/telegram-dao/internal/database/entity"
	"github.com/arandich/telegram-dao/internal/database/query"
	"github.com/arandich/telegram-dao/internal/database/query/events_query"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
	"time"
)

func EventsListAdmin(update *tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {
	eventsList, ok := query.SelectAllUsersEvents(db)
	if !ok {
		commands.ErrorMsg(update, bot, "Ошибка в запросе активностей или журнал пуст")
		return
	}
	for eventName, v := range eventsList.List {
		text := "Название: " + eventName + "\n" +
			"Имя пользователя: " + v.UserName + "\n" +
			"Дата: " + v.TimeToString() + "\n" +
			"Статус: " + v.Status
		res := tgbotapi.NewMessage(update.Message.Chat.ID, text)
		var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Подтвердить", "Подтвердить "+strconv.Itoa(v.Id)+" "+strconv.Itoa(v.Reward)),
				tgbotapi.NewInlineKeyboardButtonData("Отклонить", "Отклонить "+strconv.Itoa(v.Id)+" 0"),
			),
		)
		res.ReplyMarkup = numericKeyboard
		if _, err := bot.Send(res); err != nil {
			panic(err)
		}
	}
}

func UserEvents(update *tgbotapi.Update, bot *tgbotapi.BotAPI, user *entity.User, db *sql.DB) {
	events, ok := query.SelectAllUserEvents(user, db)
	if !ok {
		commands.ErrorMsg(update, bot, "У вас нет активностей ;(")
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

func EventsListUser(update *tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB, user *entity.User) {
	eventsList, ok := query.SelectAllEvents(db, user)
	if !ok {
		commands.ErrorMsg(update, bot, "Ошибка в запросе активностей или журнал пуст")
		return
	}
	for eventName, v := range eventsList.List {
		text := "Название: " + eventName + "\n" +
			"Дата: " + v.TimeToString() + "\n" +
			"Награда: " + strconv.Itoa(v.Reward)
		res := tgbotapi.NewMessage(update.Message.Chat.ID, text)
		var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Участвовать", "Участвовать "+strconv.Itoa(v.Id)+" "+strconv.Itoa(user.Id)),
			),
		)
		res.ReplyMarkup = numericKeyboard
		if _, err := bot.Send(res); err != nil {
			panic(err)
		}
	}
}

func CreateEvent(update *tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {
	s := strings.Split(update.Message.CommandArguments(), " ")

	if len(s) != 3 {
		commands.ErrorMsg(update, bot, "Неверные аргументы")
		return
	}

	name, date, reward := s[0], s[1], s[2]

	rewardV, err := strconv.Atoi(reward)
	if err != nil {
		commands.ErrorMsg(update, bot, "Неверный аргумент награды")
		return
	}
	if rewardV < 0 {
		commands.ErrorMsg(update, bot, "Неверный аргумент награды")
		return
	}

	dateV, err := time.Parse("2006-01-02", date)
	if err != nil {
		commands.ErrorMsg(update, bot, "Неверный аргумент даты")
		return
	}

	event := entity.Event{
		Id:     0,
		Name:   name,
		Date:   dateV,
		Reward: rewardV,
	}
	ok := events_query.AddEvent(db, &event)
	if !ok {
		commands.ErrorMsg(update, bot, "Ошибка при создании активности")
		return
	}
	commands.Msg(update, bot, "Активность: "+event.Name+"\n Успешно создана!")

}

func JoinEvent(update *tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB, eventId string, userId string) {
	eventIdV, err := strconv.Atoi(eventId)
	if err != nil {
		log.Println(err)
		return
	}
	userIdV, err := strconv.Atoi(userId)
	if err != nil {
		log.Println(err)
		return
	}
	ok := events_query.JoinEventDb(eventIdV, userIdV, db)
	if !ok {
		commands.ErrorMsg(update, bot, "Возникла ошибка при отправки заявки - "+eventId)
	}
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "Успешно!")
	if _, err := bot.Request(callback); err != nil {
		panic(err)
	}
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Заявка отправлена!")
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
	_, err = bot.Send(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID))
	if err != nil {
		return
	}
}

func AcceptEvent(update *tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB, eventId string, reward string) {
	id, err := strconv.Atoi(eventId)
	if err != nil {
		log.Println(err)
		return
	}
	karma, err := strconv.Atoi(reward)
	if err != nil {
		log.Println(err)
		return
	}
	ok := events_query.Accept(id, karma, db, update.CallbackQuery.From.UserName)
	if !ok {
		commands.ErrorMsg(update, bot, "Возникла ошибка при подтверждении заявки - "+eventId)
	}
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "Успешно!")
	if _, err := bot.Request(callback); err != nil {
		panic(err)
	}
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Заявка: "+eventId+"\n Подтверждена!")
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
	_, err = bot.Send(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID))
	if err != nil {
		return
	}
}

func DenyEvent(update *tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB, data string) {
	id, err := strconv.Atoi(data)
	if err != nil {
		log.Println(err)
		return
	}
	ok := events_query.Deny(id, db)
	if !ok {
		commands.ErrorMsg(update, bot, "Возникла ошибка при отклонении заявки - "+data)
	}
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "Успешно!")
	if _, err := bot.Request(callback); err != nil {
		panic(err)
	}
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Заявка: "+data+"\n Отклонена!")
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
	_, err = bot.Send(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID))
	if err != nil {
		return
	}
}
