package service

import (
	"context"

	"ebank/services/transaction/model"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction model.Transaction) (model.Transaction, error)
	GetTransactionByID(ctx context.Context, id int64) (model.Transaction, error)
	UpdateTransaction(ctx context.Context, transaction model.Transaction) error
	DeleteTransaction(ctx context.Context, id int64) error
	GetTransactionsByAccountID(ctx context.Context, accountID int64) ([]model.Transaction, error)
	GetAllTransactions(ctx context.Context) ([]model.Transaction, error)
}
