package service

import (
	"context"
	pb "ebank/internal/api/v1"
	"ebank/internal/model"
	"ebank/internal/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"sync"
)

type accountService struct {
	pb.UnimplementedAccountServiceServer
	userHelper            UserHelper
	accountRepository     repository.AccountRepository
	transactionRepository repository.TransactionRepository
	mutex                 sync.RWMutex
}

func NewAccountService(
	userHelper UserHelper,
	accountRepository repository.AccountRepository,
	transactionRepository repository.TransactionRepository,
) pb.AccountServiceServer {
	return &accountService{
		userHelper:            userHelper,
		accountRepository:     accountRepository,
		transactionRepository: transactionRepository,
	}
}

func (s *accountService) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.AccountResponse, error) {
	account, err := s.accountRepository.CreateAccount(ctx, model.Account{
		AccountNumber: req.AccountNumber,
		CustomerID:    req.UserId,
		CreatedAt:     timestamppb.Now().AsTime(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to save account data")
	}

	return &pb.AccountResponse{Account: &pb.Account{
		Id:            account.ID,
		AccountNumber: account.AccountNumber,
		CustomerId:    account.CustomerID,
		Balance:       account.Balance,
		CreatedAt:     timestamppb.New(account.CreatedAt),
	}}, nil
}

func (s *accountService) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.AccountResponse, error) {
	account, err := s.accountRepository.GetAccountByID(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Account not found")
	}

	if account == nil {
		return nil, status.Errorf(codes.NotFound, "Account not found")
	}

	if _, err := s.userHelper.ValidateUser(ctx, account.CustomerID); err != nil {
		return nil, err
	}

	return &pb.AccountResponse{Account: &pb.Account{
		Id:            account.ID,
		AccountNumber: account.AccountNumber,
		CustomerId:    account.CustomerID,
		Balance:       account.Balance,
		CreatedAt:     timestamppb.New(account.CreatedAt),
	}}, nil
}

func (s *accountService) UpdateAccount(ctx context.Context, req *pb.UpdateAccountRequest) (*pb.AccountResponse, error) {
	account, err := s.accountRepository.GetAccountByID(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Account not found")
	}

	if account == nil {
		return nil, status.Errorf(codes.NotFound, "Account not found")
	}

	if _, err := s.userHelper.ValidateUser(ctx, account.CustomerID); err != nil {
		return nil, err
	}

	account.AccountNumber = req.AccountNumber

	if err := s.accountRepository.UpdateAccount(ctx, *account); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to save account data")
	}

	return &pb.AccountResponse{Account: &pb.Account{
		Id:            account.ID,
		AccountNumber: account.AccountNumber,
		CustomerId:    account.CustomerID,
		Balance:       account.Balance,
		CreatedAt:     timestamppb.New(account.CreatedAt),
	}}, nil
}

func (s *accountService) DeleteAccount(ctx context.Context, req *pb.DeleteAccountRequest) (*emptypb.Empty, error) {
	account, err := s.accountRepository.GetAccountByID(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Account not found")
	}

	if account == nil {
		return nil, status.Errorf(codes.NotFound, "Account not found")
	}

	if _, err := s.userHelper.ValidateUser(ctx, account.CustomerID); err != nil {
		return nil, err
	}

	if err := s.accountRepository.DeleteAccount(ctx, req.GetId()); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *accountService) Deposit(ctx context.Context, req *pb.DepositRequest) (*pb.TransactionResponse, error) {
	s.accountRepository.LockAccountByID(ctx, req.GetAccountId())
	defer s.accountRepository.UnlockAccountByID(ctx, req.GetAccountId())
	//s.mutex.Lock()
	//defer s.mutex.Unlock()

	account, err := s.accountRepository.GetAccountByID(ctx, req.GetAccountId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Account not found")
	}

	if account == nil {
		return nil, status.Errorf(codes.NotFound, "Account not found")
	}

	if _, err := s.userHelper.ValidateUser(ctx, account.CustomerID); err != nil {
		return nil, err
	}

	account = account.AddBalance(req.GetAmount())

	if err := s.accountRepository.UpdateAccount(ctx, *account); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to save transaction data")
	}

	transaction, err := s.transactionRepository.CreateTransaction(ctx, model.Transaction{
		AccountID:       req.GetAccountId(),
		Amount:          req.GetAmount(),
		TransactionType: "DEPOSIT",
		CreatedAt:       timestamppb.Now().AsTime(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to save transaction data")
	}

	return &pb.TransactionResponse{
		Transaction: &pb.Transaction{
			Id:              transaction.ID,
			AccountId:       transaction.AccountID,
			Amount:          transaction.Amount,
			TransactionType: transaction.TransactionType,
			Timestamp:       timestamppb.New(transaction.CreatedAt),
		},
		NewBalance: account.Balance,
	}, nil
}

func (s *accountService) Withdraw(ctx context.Context, req *pb.WithdrawRequest) (*pb.TransactionResponse, error) {
	s.accountRepository.LockAccountByID(ctx, req.GetAccountId())
	defer s.accountRepository.UnlockAccountByID(ctx, req.GetAccountId())
	//s.mutex.Lock()
	//defer s.mutex.Unlock()

	account, err := s.accountRepository.GetAccountByID(ctx, req.AccountId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Account not found")
	}

	if account == nil {
		return nil, status.Errorf(codes.NotFound, "Account not found")
	}

	if _, err := s.userHelper.ValidateUser(ctx, account.CustomerID); err != nil {
		return nil, err
	}

	if account.Balance < req.Amount {
		return nil, status.Errorf(codes.FailedPrecondition, "Insufficient funds")
	}

	account = account.SubtractBalance(req.Amount)

	if err := s.accountRepository.UpdateAccount(ctx, *account); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to save transaction data")
	}

	transaction, err := s.transactionRepository.CreateTransaction(ctx, model.Transaction{
		AccountID:       req.AccountId,
		Amount:          req.Amount,
		TransactionType: "WITHDRAWAL",
		CreatedAt:       timestamppb.Now().AsTime(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to save transaction data")
	}

	return &pb.TransactionResponse{
		Transaction: &pb.Transaction{
			Id:              transaction.ID,
			AccountId:       transaction.AccountID,
			Amount:          transaction.Amount,
			TransactionType: transaction.TransactionType,
			Timestamp:       timestamppb.New(transaction.CreatedAt),
		},
		NewBalance: account.Balance,
	}, nil
}

func (s *accountService) GetTransactionHistory(ctx context.Context, req *pb.GetTransactionHistoryRequest) (*pb.GetTransactionHistoryResponse, error) {
	account, err := s.accountRepository.GetAccountByID(ctx, req.AccountId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Account not found")
	}

	if account == nil {
		return nil, status.Errorf(codes.NotFound, "Account not found")
	}

	if _, err := s.userHelper.ValidateUser(ctx, account.CustomerID); err != nil {
		return nil, err
	}

	transactions, err := s.transactionRepository.GetTransactionsByAccountID(ctx, req.GetAccountId())
	if err != nil {
		return nil, err
	}
	transactionDtos := make([]*pb.Transaction, 0)
	for _, transaction := range transactions {
		transactionDtos = append(transactionDtos, &pb.Transaction{
			Id:              transaction.ID,
			AccountId:       transaction.AccountID,
			Amount:          transaction.Amount,
			TransactionType: transaction.TransactionType,
			Timestamp:       timestamppb.New(transaction.CreatedAt),
		})
	}

	return &pb.GetTransactionHistoryResponse{Transactions: transactionDtos}, nil
}
