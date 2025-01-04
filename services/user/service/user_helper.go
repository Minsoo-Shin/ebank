package service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"ebank/pkg/jwt_manager"
	"ebank/services/user/model"
)

type UserHelper interface {
	ValidateUser(ctx context.Context, userID int64) (model.User, error)
}

type userHelper struct {
	userRepository UserRepository
}

func NewUserHelper(userRepository UserRepository) UserHelper {
	return &userHelper{
		userRepository: userRepository,
	}
}

func (u userHelper) ValidateUser(ctx context.Context, userID int64) (model.User, error) {
	user, err := u.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return model.User{}, status.Errorf(codes.Internal, err.Error())
	}

	if user == nil {
		return model.User{}, status.Errorf(codes.NotFound, "User not found")
	}

	if user.IsDeleted {
		return model.User{}, status.Errorf(codes.NotFound, "User not found")
	}

	if claims, ok := ctx.Value("user").(*jwt_manager.UserClaims); ok &&
		claims.PhoneNumber != user.PhoneNumber {
		return model.User{}, status.Errorf(codes.PermissionDenied, "Not allowed")
	}

	return *user, nil
}
