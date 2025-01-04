package model

import "time"

type Account struct {
	ID            int64
	AccountNumber string
	CustomerID    int64
	Balance       float64
	CreatedAt     time.Time
}

func (r *Account) AddBalance(amount float64) *Account {
	r.Balance += amount
	return r
}

func (r *Account) SubtractBalance(amount float64) *Account {
	r.Balance -= amount
	return r
}
