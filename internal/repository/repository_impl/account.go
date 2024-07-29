package repository_impl

import (
	"context"
	"ebank/internal/repository"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"ebank/internal/model"
)

type accountFileRepository struct {
	nextID           int64
	accounts         map[int64]model.Account
	accountsByUserID map[int64][]int64
	accountMutex     map[int64]*sync.RWMutex
	mapMutex         sync.RWMutex
	fileMutex        sync.RWMutex
	filePath         string
}

func NewAccountFileRepository(filePath string) (repository.AccountRepository, error) {
	repo := &accountFileRepository{
		accounts:         make(map[int64]model.Account),
		accountsByUserID: make(map[int64][]int64),
		accountMutex:     make(map[int64]*sync.RWMutex),
		filePath:         filePath,
	}

	if err := repo.load(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *accountFileRepository) load() error {
	data, err := ioutil.ReadFile(r.filePath)
	if os.IsNotExist(err) {
		return nil // 파일이 없으면 새로 시작
	} else if err != nil {
		return err
	}

	var accounts []model.Account
	if err := json.Unmarshal(data, &accounts); err != nil {
		return err
	}

	for i, account := range accounts {
		r.accounts[account.ID] = account
		r.accountsByUserID[account.CustomerID] = append(r.accountsByUserID[account.CustomerID], account.ID)
		if i == len(accounts)-1 {
			r.nextID = account.ID
		}
	}

	return nil
}

func (r *accountFileRepository) save() error {
	r.fileMutex.Lock()
	defer r.fileMutex.Unlock()

	accounts := make([]model.Account, 0, len(r.accounts))
	for _, account := range r.accounts {
		accounts = append(accounts, account)
	}

	data, err := json.Marshal(accounts)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(r.filePath, data, 0644)
}

func (r *accountFileRepository) LockAccountByID(ctx context.Context, id int64) error {
	r.mapMutex.Lock()
	mutex, ok := r.accountMutex[id]
	if !ok {
		r.accountMutex[id] = &sync.RWMutex{}
		mutex = r.accountMutex[id]
	}
	r.mapMutex.Unlock()

	mutex.Lock()

	return nil
}

func (r *accountFileRepository) UnlockAccountByID(ctx context.Context, id int64) error {
	r.mapMutex.Lock()
	mutex := r.accountMutex[id]
	r.mapMutex.Unlock()

	mutex.Unlock()
	return nil
}

func (r *accountFileRepository) CreateAccount(ctx context.Context, account model.Account) (model.Account, error) {
	r.nextID++
	account.ID = r.nextID

	r.accounts[account.ID] = account
	r.accountsByUserID[account.CustomerID] = append(r.accountsByUserID[account.CustomerID], account.ID)

	if err := r.save(); err != nil {
		return model.Account{}, err
	}

	return account, nil
}

func (r *accountFileRepository) GetAccountByID(ctx context.Context, id int64) (*model.Account, error) {
	account, exists := r.accounts[id]
	if !exists {
		return nil, fmt.Errorf("account with ID %d not found", id)
	}

	return &account, nil
}

func (r *accountFileRepository) UpdateAccount(ctx context.Context, account model.Account) error {
	if _, exists := r.accounts[account.ID]; !exists {
		return fmt.Errorf("account with ID %d not found", account.ID)
	}

	r.accounts[account.ID] = account

	return r.save()
}

func (r *accountFileRepository) DeleteAccount(ctx context.Context, id int64) error {
	account, exists := r.accounts[id]
	if !exists {
		return fmt.Errorf("account with ID %d not found", id)
	}

	delete(r.accounts, id)

	// Remove account from accountsByUserID
	userAccounts := r.accountsByUserID[account.CustomerID]
	for i, accID := range userAccounts {
		if accID == id {
			r.accountsByUserID[account.CustomerID] = append(userAccounts[:i], userAccounts[i+1:]...)
			break
		}
	}

	return r.save()
}

func (r *accountFileRepository) GetAccountsByUserID(ctx context.Context, userID int64) ([]model.Account, error) {
	accountIDs, exists := r.accountsByUserID[userID]
	if !exists {
		return []model.Account{}, nil
	}

	accounts := make([]model.Account, len(accountIDs))
	for i, id := range accountIDs {
		accounts[i] = r.accounts[id]
	}

	return accounts, nil
}

func (r *accountFileRepository) GetAllAccounts(ctx context.Context) ([]model.Account, error) {
	accounts := make([]model.Account, 0, len(r.accounts))
	for _, account := range r.accounts {
		accounts = append(accounts, account)
	}

	return accounts, nil
}
