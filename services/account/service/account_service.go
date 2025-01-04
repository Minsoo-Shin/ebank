package service

import (
	"context"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"ebank/api/v1"
	"ebank/services/account/model"
)

type accountService struct {
	ebank.UnimplementedAccountServiceServer
	// userHelper        UserHelper
	accountRepository AccountRepository
	mutex             sync.RWMutex
}

func NewAccountService(
	// userHelper UserHelper,
	accountRepository AccountRepository,
) ebank.AccountServiceServer {
	return &accountService{
		// userHelper:            userHelper,
		accountRepository: accountRepository,
	}
}

func (s *accountService) CreateAccount(ctx context.Context, req *ebank.CreateAccountRequest) (*ebank.AccountResponse, error) {
	account, err := s.accountRepository.CreateAccount(ctx, model.Account{
		AccountNumber: req.AccountNumber,
		CustomerID:    req.UserId,
		CreatedAt:     timestamppb.Now().AsTime(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to save account data")
	}

	return &ebank.AccountResponse{Account: &ebank.Account{
		Id:            account.ID,
		AccountNumber: account.AccountNumber,
		CustomerId:    account.CustomerID,
		Balance:       account.Balance,
		CreatedAt:     timestamppb.New(account.CreatedAt),
	}}, nil
}

func (s *accountService) GetAccount(ctx context.Context, req *ebank.GetAccountRequest) (*ebank.AccountResponse, error) {
	account, err := s.accountRepository.GetAccountByID(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Account not found")
	}

	if account == nil {
		return nil, status.Errorf(codes.NotFound, "Account not found")
	}

	// if _, err := s.userHelper.ValidateUser(ctx, account.CustomerID); err != nil {
	// 	return nil, err
	// }

	return &ebank.AccountResponse{Account: &ebank.Account{
		Id:            account.ID,
		AccountNumber: account.AccountNumber,
		CustomerId:    account.CustomerID,
		Balance:       account.Balance,
		CreatedAt:     timestamppb.New(account.CreatedAt),
	}}, nil
}

func (s *accountService) UpdateAccount(ctx context.Context, req *ebank.UpdateAccountRequest) (*ebank.AccountResponse, error) {
	account, err := s.accountRepository.GetAccountByID(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Account not found")
	}

	if account == nil {
		return nil, status.Errorf(codes.NotFound, "Account not found")
	}

	// if _, err := s.userHelper.ValidateUser(ctx, account.CustomerID); err != nil {
	// 	return nil, err
	// }

	account.AccountNumber = req.AccountNumber

	if err := s.accountRepository.UpdateAccount(ctx, *account); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to save account data")
	}

	return &ebank.AccountResponse{Account: &ebank.Account{
		Id:            account.ID,
		AccountNumber: account.AccountNumber,
		CustomerId:    account.CustomerID,
		Balance:       account.Balance,
		CreatedAt:     timestamppb.New(account.CreatedAt),
	}}, nil
}

func (s *accountService) DeleteAccount(ctx context.Context, req *ebank.DeleteAccountRequest) (*emptypb.Empty, error) {
	account, err := s.accountRepository.GetAccountByID(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Account not found")
	}

	if account == nil {
		return nil, status.Errorf(codes.NotFound, "Account not found")
	}

	// if _, err := s.userHelper.ValidateUser(ctx, account.CustomerID); err != nil {
	// 	return nil, err
	// }

	if err := s.accountRepository.DeleteAccount(ctx, req.GetId()); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
