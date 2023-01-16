package events_commands

import (
	"database/sql"
	"github.com/arandich/telegram-dao/internal/commands"
	"github.com/arandich/telegram-dao/internal/database/query"
	"github.com/arandich/telegram-dao/internal/database/query/events_query"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
)

func EventsList(update *tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {
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
				tgbotapi.NewInlineKeyboardButtonData("Подтвердить", "Подтвердить "+strconv.Itoa(v.Id)),
				tgbotapi.NewInlineKeyboardButtonData("Отклонить", "Отклонить "+strconv.Itoa(v.Id)),
			),
		)
		res.ReplyMarkup = numericKeyboard
		if _, err := bot.Send(res); err != nil {
			panic(err)
		}
	}
}

func AcceptEvent(update *tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB, data string) {
	id, err := strconv.Atoi(data)
	if err != nil {
		log.Println(err)
		return
	}
	ok := events_query.Accept(id, db)
	if !ok {
		commands.ErrorMsg(update, bot, "Возникла ошибка при подтверждении заявки - "+data)
	}
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "Успешно!")
	if _, err := bot.Request(callback); err != nil {
		panic(err)
	}
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Заявка: "+data+"\n Подтверждена!")
	if _, err := bot.Send(msg); err != nil {
		panic(err)
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
}
