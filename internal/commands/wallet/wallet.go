package wallet

import (
	"database/sql"
	"github.com/arandich/telegram-dao/internal/commands"
	"github.com/arandich/telegram-dao/internal/database/entity"
	"github.com/arandich/telegram-dao/internal/database/query/wallet_query"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func AddWallet(update *tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB, user *entity.User) bool {
	if update.Message.CommandArguments() != "" && len(update.Message.CommandArguments()) == 64 {
		user.TonWallet = update.Message.CommandArguments()
	} else {
		commands.ErrorMsg(update, bot, "Invalid wallet")
		return false
	}

	ok := wallet_query.CreateWallet(db, user)
	if !ok {
		log.Println("Ошибка добавления кошелька")
		return false
	}
	commands.Msg(update, bot, "Кошелек успешно добавлен")
	return true
}
