package wallet_query

import (
	"database/sql"
	"github.com/arandich/telegram-dao/internal/database/entity"
	"log"
)

func CreateWallet(db *sql.DB, user *entity.User) bool {
	result, err := db.Exec(`UPDATE users set ton_wallet = $1 WHERE username = $2`, user.TonWallet, user.Username)
	if err != nil {
		log.Println(err)
		return false
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return false
	}
	return true
}
