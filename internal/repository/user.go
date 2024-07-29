package repository

import (
	"context"
	"ebank/internal/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user model.User) (model.User, error)
	GetUserByID(ctx context.Context, id int64) (*model.User, error)
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (model.User, error)
	UpdateUser(ctx context.Context, user model.User) error
	GetAllUsers(ctx context.Context, isDeleted *bool) ([]model.User, error)
}
