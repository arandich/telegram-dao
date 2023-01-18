package callback

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func CallBackReq(update *tgbotapi.Update, bot *tgbotapi.BotAPI, text string) {
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, text)
	if _, err := bot.Request(callback); err != nil {
		log.Println(err)
	}
	_, err := bot.Send(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID))
	if err != nil {
		log.Println(err)
	}
}

func CallBackMsg(update *tgbotapi.Update, bot *tgbotapi.BotAPI, text string) {
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}
