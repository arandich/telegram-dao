package entity

import (
	"time"
)

type Vote struct {
	Id        int
	Name      string
	Url       string
	DateStart time.Time
	DateEnd   time.Time
	Text1     string
	Text2     string
	Text3     string
	Var1      int
	Var2      int
	Var3      int
}

type VoteList struct {
	List map[string]Vote
}

type VoteArr struct {
	List []Vote
}

type UserVote struct {
	VoteId int
	UserId int
	Amount int
	Choice string
}
