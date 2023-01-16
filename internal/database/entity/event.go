package entity

import "time"

type Event struct {
	Id   int
	Name string
	Date time.Time
}

type UserEvent struct {
	Id       int
	Name     string
	UserName string
	Date     time.Time
	Status   string
}

func (ue *UserEvent) TimeToString() string {
	return ue.Date.Format("2006-01-02")
}
