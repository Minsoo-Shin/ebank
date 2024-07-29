package repository

import (
	"context"
	"ebank/internal/model"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, account model.Account) (model.Account, error)
	GetAccountByID(ctx context.Context, id int64) (*model.Account, error)
	UpdateAccount(ctx context.Context, account model.Account) error
	DeleteAccount(ctx context.Context, id int64) error
	GetAccountsByUserID(ctx context.Context, userID int64) ([]model.Account, error)
	GetAllAccounts(ctx context.Context) ([]model.Account, error)
	LockAccountByID(ctx context.Context, id int64) error
	UnlockAccountByID(ctx context.Context, id int64) error
}
