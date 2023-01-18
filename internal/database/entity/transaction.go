package entity

import "time"

type Transaction struct {
	TrId       int
	Sender     string
	ToUsername string
	Amount     int
	Date       time.Time
}

type TransactionList struct {
	List []Transaction
}

func (tr *Transaction) TimeToString() string {
	return tr.Date.Format("2006-01-02-15-01")
}
