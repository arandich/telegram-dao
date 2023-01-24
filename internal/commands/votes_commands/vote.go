package votes_commands

import (
	"database/sql"
	"encoding/json"
	"github.com/arandich/telegram-dao/internal/commands"
	"github.com/arandich/telegram-dao/internal/database/entity"
	"github.com/arandich/telegram-dao/internal/database/query/votes_query"
	"github.com/arandich/telegram-dao/pkg/response/callback"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type UVote struct {
	Username string
	VoteId   int
	Choice   string
}

func CreateVote(update *tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {
	vote := entity.Vote{}
	s := strings.Split(update.Message.CommandArguments(), " ")
	if len(s) != 6 {
		commands.ErrorMsg(update, bot, "Неверные аргументы")
		return
	}
	timeEnd, err := time.Parse("2006-01-02", s[2])
	if err != nil {
		commands.ErrorMsg(update, bot, "Неверный аргумент даты")
		log.Println(err)
		return
	}
	vote.Name, vote.Url, vote.DateEnd, vote.Text1, vote.Text2, vote.Text3 = s[0], s[1], timeEnd, s[3], s[4], s[5]

	ok := votes_query.AddVote(db, &vote)
	if !ok {
		commands.ErrorMsg(update, bot, "Ошибка создания голосования")
		return
	}

	votes_query.SelectVote(db, &vote)

	data := []UVote{
		{Username: "System",
			VoteId: 0,
			Choice: "System",
		},
	}

	jsonString, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return
	}

	err = os.WriteFile("storage/votes/"+strconv.Itoa(vote.Id)+".json", jsonString, os.ModePerm)
	if err != nil {
		log.Println(err)
		return
	}
	commands.Msg(update, bot, "Голосование успешно создано")
}

func AllVotes(update *tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB, user *entity.User) {
	allVotes, ok := votes_query.SelectVotes(db, user)
	if !ok {
		commands.ErrorMsg(update, bot, "Ошибка получения списка голосований")
		return
	}

	var votes_data []entity.Vote

	for _, v := range allVotes.List {
		data := readJSONToken("storage/votes/" + strconv.Itoa(v.Id) + ".json")
		for i, v2 := range data {
			if v2.VoteId != v.Id && v2.Username != update.Message.From.UserName {
				votes_data = append(votes_data, v)

			} else {
				if len(votes_data) == 1 {
					votes_data[0] = entity.Vote{}
				} else {
					votes_data = append(votes_data[:i], votes_data[i+1:]...)
				}
			}
		}
	}
	if votes_data[0].Name == "" {
		commands.Msg(update, bot, "Сейчас нет активных голосований")
		return
	}
	for _, v := range votes_data {
		text := "*" + v.Name + "* \n" +
			"Дата старта: " + v.DateStart.Format("2006-01-02") + "\n" +
			"Дата Окончания: " + v.DateEnd.Format("2006-01-02") + "\n" +
			"Ссылка на документ: " + v.Url
		res := tgbotapi.NewMessage(update.Message.Chat.ID, text)
		if v.Text3 != "-" {
			var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(v.Text1+" +"+strconv.Itoa(user.Karma), "Вариант1 "+strconv.Itoa(v.Id)+" "+strconv.Itoa(user.Id)+" "+strconv.Itoa(user.Karma)+" "+v.Text1),
					tgbotapi.NewInlineKeyboardButtonData(v.Text2+" +"+strconv.Itoa(user.Karma), "Вариант2 "+strconv.Itoa(v.Id)+" "+strconv.Itoa(user.Id)+" "+strconv.Itoa(user.Karma)+" "+v.Text2),
					tgbotapi.NewInlineKeyboardButtonData(v.Text3+" +"+strconv.Itoa(user.Karma), "Вариант3 "+strconv.Itoa(v.Id)+" "+strconv.Itoa(user.Id)+" "+strconv.Itoa(user.Karma)+" "+v.Text3),
				),
			)
			res.ReplyMarkup = numericKeyboard
			if _, err := bot.Send(res); err != nil {
				panic(err)
			}
		} else {
			var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(v.Text1+" +"+strconv.Itoa(user.Karma), "Вариант1 "+strconv.Itoa(v.Id)+" "+strconv.Itoa(user.Id)+" "+strconv.Itoa(user.Karma)+" "+v.Text1),
					tgbotapi.NewInlineKeyboardButtonData(v.Text2+" -"+strconv.Itoa(user.Karma), "Вариант2 "+strconv.Itoa(v.Id)+" "+strconv.Itoa(user.Id)+" "+strconv.Itoa(user.Karma)+" "+v.Text1),
				),
			)
			res.ReplyMarkup = numericKeyboard
			if _, err := bot.Send(res); err != nil {
				panic(err)
			}
		}
	}
}

func Yes(update *tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB, vote *entity.UserVote) bool {

	ok := votes_query.VoteYes(db, vote)
	if !ok {
		return false
	}
	return writeAndRead(update, bot, vote)
}

func No(update *tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB, vote *entity.UserVote) bool {
	ok := votes_query.VoteNo(db, vote)
	if !ok {
		return false
	}
	return writeAndRead(update, bot, vote)
}

func Neutral(update *tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB, vote *entity.UserVote) bool {
	ok := votes_query.VoteElse(db, vote)
	if !ok {
		return false
	}
	return writeAndRead(update, bot, vote)
}

func writeAndRead(update *tgbotapi.Update, bot *tgbotapi.BotAPI, vote *entity.UserVote) bool {
	_, err := os.Stat("storage/votes/" + strconv.Itoa(vote.VoteId) + ".json")
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("File does not exist.")
			return false
		}
	}

	data := readJSONToken("storage/votes/" + strconv.Itoa(vote.VoteId) + ".json")
	data = append(data, UVote{
		Username: update.CallbackQuery.From.UserName,
		VoteId:   vote.VoteId,
		Choice:   vote.Choice,
	})
	log.Println(data)
	jsonString, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return false
	}

	err = os.WriteFile("storage/votes/"+strconv.Itoa(vote.VoteId)+".json", jsonString, os.ModePerm)
	if err != nil {
		log.Println(err)
		return false
	}
	callback.CallBackReq(update, bot, "Успешно")
	callback.CallBackMsg(update, bot, "Ваш голос засчитан!")
	return true
}

func UserVotes(update *tgbotapi.Update, bot *tgbotapi.BotAPI, user *entity.User, db *sql.DB) {
	votes, ok := votes_query.SelectVotesArr(db, user)
	if !ok {
		commands.ErrorMsg(update, bot, "Список пуст ;(")
		return
	}
	type vote struct {
		Name      string
		DateStart time.Time
		DateEnd   time.Time
		Url       string
		Choice    string
	}
	var votes_data []vote
	for _, v := range votes.List {
		data := readJSONToken("storage/votes/" + strconv.Itoa(v.Id) + ".json")

		for _, v2 := range data {

			if v2.VoteId == v.Id && v2.Username == update.Message.From.UserName {
				votes_data = append(votes_data, vote{
					Name:      v.Name,
					DateStart: v.DateStart,
					DateEnd:   v.DateEnd,
					Url:       v.Url,
					Choice:    v2.Choice,
				})

			} else {

			}

		}
	}
	for _, val := range votes_data {
		text := val.Name + "\n" + "Дата начала: " + val.DateStart.Format("2006-01-02") + "\n" + "Дата окончания: " + val.DateEnd.Format("2006-01-02") + "\n" + "Ссылка: " + val.Url + "\nВы выбрали: " + val.Choice
		commands.MsgWithoutReply(update, bot, text)
	}
}

func readJSONToken(fileName string) []UVote {
	file, _ := os.Open(fileName)
	defer file.Close()

	decoder := json.NewDecoder(file)

	var filteredData []UVote

	// Read the array open bracket
	_, err := decoder.Token()
	if err != nil {
		return nil
	}

	data := UVote{}
	for decoder.More() {
		err := decoder.Decode(&data)
		if err != nil {
			return nil
		}
		filteredData = append(filteredData, data)

	}

	return filteredData
}
