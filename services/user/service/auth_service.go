package service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ebank "ebank/api/v1"
	"ebank/pkg/jwt_manager"
	"ebank/pkg/zero"
)

type authService struct {
	ebank.UnimplementedAuthServiceServer
	userRepository UserRepository
	jwtManager     jwt_manager.JWTManager
}

func (a *authService) Login(ctx context.Context, req *ebank.LoginRequest) (*ebank.LoginResponse, error) {
	user, err := a.userRepository.GetUserByPhoneNumber(ctx, req.GetPhoneNumber())
	if err != nil {
		return nil, err
	}
	if zero.IsStructZero(user) {
		return nil, status.Errorf(codes.NotFound, "User not found")
	}

	if !user.IsCorrectPassword(req.Password) {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid password")
	}

	tokenString, err := a.jwtManager.Generate(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to generate token")
	}

	return &ebank.LoginResponse{Token: tokenString}, nil
}

func NewAuthService(
	userRepository UserRepository,
	jwtManager jwt_manager.JWTManager,
) ebank.AuthServiceServer {
	return &authService{
		userRepository: userRepository,
		jwtManager:     jwtManager,
	}
}
