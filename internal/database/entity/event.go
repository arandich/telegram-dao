package entity

import "time"

type Event struct {
	Id     int
	Name   string
	Date   time.Time
	UserId int
	Status string
}
