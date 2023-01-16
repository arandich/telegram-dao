package tokens

import (
	"database/sql"
	"github.com/arandich/telegram-dao/internal/commands"
	"github.com/arandich/telegram-dao/internal/database/entity"
	"github.com/arandich/telegram-dao/internal/database/query"
	"github.com/arandich/telegram-dao/internal/database/query/token"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

func SendTokensTo(update *tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB, user *entity.User) {
	s := strings.Split(update.Message.CommandArguments(), " ")
	if len(s) != 2 {
		commands.ErrorMsg(update, bot, "Неправильные аргументы команды")
		return
	}
	usernameTo, amount := s[0], s[1]
	if len(usernameTo) > 32 || len(usernameTo) < 5 {
		commands.ErrorMsg(update, bot, "Имя пользователя указано неверно")
		return
	} else if _, err := strconv.Atoi(amount); err != nil {
		commands.ErrorMsg(update, bot, "Указано неверное количество токенов")
		return
	}
	am, _ := strconv.Atoi(amount)

	if user.Tokens < am || am == 0 || am < 0 {
		commands.ErrorMsg(update, bot, "У вас недостаточно токенов или значение неверно ;(")
		return
	}

	_, ok := query.FindByUsername(usernameTo, db)
	if !ok {
		commands.ErrorMsg(update, bot, "Пользователь не найден")
		return
	}

	ok = token.UpdateUserToken(usernameTo, update.Message.From.UserName, am, db)
	if ok {
		commands.Msg(update, bot, "Токены успешно переведены")
	} else {
		commands.ErrorMsg(update, bot, "Возникла ошибка при отправке токенов")
	}
}
