package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"ebank/services/transaction/model"
	"ebank/services/transaction/service"
)

type transactionFileRepository struct {
	nextID                  int64
	transactions            map[int64]model.Transaction
	transactionsByAccountID map[int64][]int64
	mapMutex                sync.RWMutex
	fileMutex               sync.RWMutex
	filePath                string
}

func NewTransactionFileRepository(filePath string) (service.TransactionRepository, error) {
	repo := &transactionFileRepository{
		transactions:            make(map[int64]model.Transaction),
		transactionsByAccountID: make(map[int64][]int64),
		filePath:                filePath,
	}

	if err := repo.load(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *transactionFileRepository) load() error {
	data, err := ioutil.ReadFile(r.filePath)
	if os.IsNotExist(err) {
		return nil // 파일이 없으면 새로 시작
	} else if err != nil {
		return err
	}

	var transactions []model.Transaction
	if err := json.Unmarshal(data, &transactions); err != nil {
		return err
	}

	for i, transaction := range transactions {
		r.transactions[transaction.ID] = transaction
		r.transactionsByAccountID[transaction.AccountID] = append(r.transactionsByAccountID[transaction.AccountID], transaction.ID)
		if i == len(transactions)-1 {
			r.nextID = transaction.ID
		}
	}

	return nil
}

func (r *transactionFileRepository) save() error {
	r.fileMutex.RLock()
	defer r.fileMutex.RUnlock()

	transactions := make([]model.Transaction, 0, len(r.transactions))
	for _, transaction := range r.transactions {
		transactions = append(transactions, transaction)
	}

	data, err := json.Marshal(transactions)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(r.filePath, data, 0644)
}

func (r *transactionFileRepository) CreateTransaction(ctx context.Context, transaction model.Transaction) (model.Transaction, error) {
	r.mapMutex.Lock()
	defer r.mapMutex.Unlock()

	r.nextID++
	transaction.ID = r.nextID

	r.transactions[transaction.ID] = transaction
	r.transactionsByAccountID[transaction.AccountID] = append(r.transactionsByAccountID[transaction.AccountID], transaction.ID)

	if err := r.save(); err != nil {
		return model.Transaction{}, err
	}

	return transaction, nil
}

func (r *transactionFileRepository) GetTransactionByID(ctx context.Context, id int64) (model.Transaction, error) {
	r.mapMutex.RLock()
	defer r.mapMutex.RUnlock()

	transaction, exists := r.transactions[id]
	if !exists {
		return model.Transaction{}, fmt.Errorf("transaction with ID %d not found", id)
	}

	return transaction, nil
}

func (r *transactionFileRepository) UpdateTransaction(ctx context.Context, transaction model.Transaction) error {
	r.mapMutex.Lock()
	defer r.mapMutex.Unlock()

	if _, exists := r.transactions[transaction.ID]; !exists {
		return fmt.Errorf("transaction with ID %d not found", transaction.ID)
	}

	r.transactions[transaction.ID] = transaction

	return r.save()
}

func (r *transactionFileRepository) DeleteTransaction(ctx context.Context, id int64) error {
	r.mapMutex.Lock()
	defer r.mapMutex.Unlock()

	transaction, exists := r.transactions[id]
	if !exists {
		return fmt.Errorf("transaction with ID %d not found", id)
	}

	delete(r.transactions, id)

	// Remove transaction from transactionsByAccountID
	userTransactions := r.transactionsByAccountID[transaction.AccountID]
	for i, accID := range userTransactions {
		if accID == id {
			r.transactionsByAccountID[transaction.AccountID] = append(userTransactions[:i], userTransactions[i+1:]...)
			break
		}
	}

	return r.save()
}

func (r *transactionFileRepository) GetTransactionsByAccountID(ctx context.Context, accountID int64) ([]model.Transaction, error) {
	r.mapMutex.RLock()
	defer r.mapMutex.RUnlock()

	transactionIDs, exists := r.transactionsByAccountID[accountID]
	if !exists {
		return []model.Transaction{}, nil
	}

	transactions := make([]model.Transaction, len(transactionIDs))
	for i, id := range transactionIDs {
		transactions[i] = r.transactions[id]
	}

	return transactions, nil
}

func (r *transactionFileRepository) GetAllTransactions(ctx context.Context) ([]model.Transaction, error) {
	r.mapMutex.RLock()
	defer r.mapMutex.RUnlock()

	transactions := make([]model.Transaction, 0, len(r.transactions))
	for _, transaction := range r.transactions {
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
