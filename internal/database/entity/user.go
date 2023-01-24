package entity

import (
	"strconv"
	"time"
)

type User struct {
	Id        int
	Username  string
	RoleId    int
	Karma     int
	Tokens    int
	CreatedAt time.Time
	TonWallet string
}

func (u *User) KarmaToString() string {
	return strconv.Itoa(u.Karma)
}

func (u *User) TokenToString() string {
	return strconv.Itoa(u.Tokens)
}
