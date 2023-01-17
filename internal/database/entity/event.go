package entity

import "time"

type Event struct {
	Id     int
	Name   string
	Date   time.Time
	Reward int
}

type UserEvent struct {
	Id       int
	Name     string
	UserName string
	Date     time.Time
	Status   string
	Reward   int
}

func (ue *UserEvent) TimeToString() string {
	return ue.Date.Format("2006-01-02")
}

func (ue *Event) TimeToString() string {
	return ue.Date.Format("2006-01-02")
}
