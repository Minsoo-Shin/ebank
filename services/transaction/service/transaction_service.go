package service

import (
	ebank "ebank/api/v1"
)

type transactionService struct {
	ebank.UnimplementedTransactionServiceServer
	transactionRepository TransactionRepository
}

func NewTransactionService(
	transactionRepository TransactionRepository,
) ebank.TransactionServiceServer {
	return &transactionService{
		transactionRepository: transactionRepository,
	}
}
