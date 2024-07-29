package service

import (
	"context"
	pb "ebank/internal/api/v1"
	"ebank/internal/model"
	"ebank/internal/repository"
	"ebank/pkg/jwt_manager"
	"ebank/pkg/zero"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type userService struct {
	pb.UnimplementedUserServiceServer
	userHelper        UserHelper
	userRepository    repository.UserRepository
	accountRepository repository.AccountRepository
	jwtManager        jwt_manager.JWTManager
}

func NewUserService(
	userHelper UserHelper,
	userRepository repository.UserRepository,
	accountRepository repository.AccountRepository,
	jwtManager jwt_manager.JWTManager,
) pb.UserServiceServer {
	return &userService{
		userHelper:        userHelper,
		userRepository:    userRepository,
		accountRepository: accountRepository,
		jwtManager:        jwtManager,
	}
}

func (s *userService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to hash password")
	}

	user := model.User{
		Name:        req.Name,
		Birth:       req.Birth,
		PhoneNumber: req.PhoneNumber,
		Password:    string(hashedPassword),
	}

	user, err = s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{User: &pb.User{
		Id:          user.ID,
		Name:        user.Name,
		Birth:       user.Birth,
		PhoneNumber: user.PhoneNumber,
	}}, nil
}

func (s *userService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	user, err := s.userRepository.GetUserByID(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "User not found")
	}

	if user == nil {
		return nil, status.Errorf(codes.NotFound, "User not found")
	}

	accounts, err := s.accountRepository.GetAccountsByUserID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	accountDtos := make([]*pb.Account, 0)
	for _, account := range accounts {
		accountDtos = append(accountDtos, &pb.Account{
			Id:            account.ID,
			AccountNumber: account.AccountNumber,
			CustomerId:    account.CustomerID,
			Balance:       account.Balance,
			CreatedAt:     timestamppb.New(account.CreatedAt),
		})
	}

	return &pb.UserResponse{User: &pb.User{
		Id:          user.ID,
		Name:        user.Name,
		Birth:       user.Birth,
		PhoneNumber: user.PhoneNumber,
		Accounts:    accountDtos,
	}}, nil
}

func (s *userService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	validateUser, err := s.userHelper.ValidateUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	validateUser.Name = req.Name
	validateUser.Birth = req.Birth
	validateUser.PhoneNumber = req.PhoneNumber

	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to hash password")
		}
		validateUser.Password = string(hashedPassword)
	}

	if err := s.userRepository.UpdateUser(ctx, validateUser); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to save user data")
	}

	return &pb.UserResponse{User: &pb.User{
		Id:          validateUser.ID,
		Name:        validateUser.Name,
		Birth:       validateUser.Birth,
		PhoneNumber: validateUser.PhoneNumber,
	}}, nil
}

func (s *userService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	validateUser, err := s.userHelper.ValidateUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	validateUser.PhoneNumber = validateUser.MaskPhoneNumber()
	validateUser.IsDeleted = true

	if err = s.userRepository.UpdateUser(ctx, validateUser); err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "Failed to save user data")
	}

	return &emptypb.Empty{}, nil
}

func (s *userService) GetAllUsers(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.UserListResponse, error) {
	users, err := s.userRepository.GetAllUsers(ctx, req.IsDeleted)
	if err != nil {
		return nil, err
	}

	userDtos := &pb.UserListResponse{Users: make([]*pb.User, 0, len(users))}
	for _, user := range users {
		userDtos.Users = append(userDtos.Users, &pb.User{
			Id:          user.ID,
			Name:        user.Name,
			Birth:       user.Birth,
			PhoneNumber: user.PhoneNumber,
		})
	}

	return userDtos, nil
}

func (s *userService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := s.userRepository.GetUserByPhoneNumber(ctx, req.GetPhoneNumber())
	if err != nil {
		return nil, err
	}
	if zero.IsStructZero(user) {
		return nil, status.Errorf(codes.NotFound, "User not found")
	}

	if !user.IsCorrectPassword(req.Password) {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid password")
	}

	tokenString, err := s.jwtManager.Generate(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to generate token")
	}

	return &pb.LoginResponse{Token: tokenString}, nil
}
