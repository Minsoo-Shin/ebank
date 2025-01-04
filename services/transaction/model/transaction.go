package model

import "time"

type Transaction struct {
	ID              int64
	AccountID       int64
	Amount          float64
	TransactionType string
	CreatedAt       time.Time
}
