package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"ebank/services/user/model"
	"ebank/services/user/service"
)

type userFileRepository struct {
	users              map[int64]model.User
	usersByPhoneNumber map[string]int64
	mapMutex           sync.RWMutex
	fileMutex          sync.RWMutex
	nextID             int64
	filePath           string
}

func NewUserFileRepository(filePath string) (service.UserRepository, error) {
	repo := &userFileRepository{
		users:              make(map[int64]model.User),
		usersByPhoneNumber: make(map[string]int64),
		filePath:           filePath,
		mapMutex:           sync.RWMutex{},
		fileMutex:          sync.RWMutex{},
	}

	if err := repo.load(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *userFileRepository) load() error {
	data, err := ioutil.ReadFile(r.filePath)
	if os.IsNotExist(err) {
		return nil // 파일이 없으면 새로 시작
	} else if err != nil {
		return err
	}

	var users []model.User
	if err := json.Unmarshal(data, &users); err != nil {
		return err
	}

	for i, user := range users {
		r.users[user.ID] = user
		r.usersByPhoneNumber[user.PhoneNumber] = user.ID
		if i == len(users)-1 {
			r.nextID = user.ID
		}
	}

	return nil
}

func (r *userFileRepository) save() error {
	users := make([]model.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	data, err := json.Marshal(users)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(r.filePath, data, 0644)
}

func (r *userFileRepository) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	r.mapMutex.Lock()
	defer r.mapMutex.Unlock()

	if _, exists := r.usersByPhoneNumber[user.PhoneNumber]; exists {
		return model.User{}, fmt.Errorf("user with phone number %s already exists", user.PhoneNumber)
	}
	r.nextID++
	user.ID = r.nextID
	r.users[user.ID] = user
	r.usersByPhoneNumber[user.PhoneNumber] = user.ID

	if err := r.save(); err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (r *userFileRepository) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	r.mapMutex.RLock()
	defer r.mapMutex.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, fmt.Errorf("user with ID %d not found", id)
	}
	if user.IsDeleted {
		return nil, fmt.Errorf("user with ID %d not found", id)
	}

	return &user, nil
}

func (r *userFileRepository) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (model.User, error) {
	r.mapMutex.RLock()
	defer r.mapMutex.RUnlock()

	id, exists := r.usersByPhoneNumber[phoneNumber]
	if !exists {
		return model.User{}, fmt.Errorf("user with phone number %s not found", phoneNumber)
	}

	return r.users[id], nil
}

func (r *userFileRepository) UpdateUser(ctx context.Context, user model.User) error {
	r.mapMutex.Lock()
	defer r.mapMutex.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return fmt.Errorf("user with ID %d not found", user.ID)
	}

	oldUser := r.users[user.ID]
	if oldUser.PhoneNumber != user.PhoneNumber {
		delete(r.usersByPhoneNumber, oldUser.PhoneNumber)
		r.usersByPhoneNumber[user.PhoneNumber] = user.ID
	}

	r.users[user.ID] = user

	return r.save()
}

func (r *userFileRepository) GetAllUsers(ctx context.Context, isDeleted *bool) ([]model.User, error) {
	r.mapMutex.RLock()
	defer r.mapMutex.RUnlock()

	users := make([]model.User, 0, len(r.users))
	for _, user := range r.users {
		if isDeleted != nil && *isDeleted != user.IsDeleted {
			continue
		}
		users = append(users, user)
	}

	return users, nil
}
