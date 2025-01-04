package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"ebank/api/v1"
	"ebank/internal/model"
	"ebank/mocks"
)

func TestUserUsecaseSuite(t *testing.T) {
	suite.Run(t, new(UserUsecaseTestSuite))
}

var (
	_ suite.SetupTestSuite = &UserUsecaseTestSuite{}
)

type UserUsecaseTestSuite struct {
	suite.Suite
	userRepository    *mocks.UserRepository
	accountRepository *mocks.AccountRepository
	userHelper        *mocks.UserHelper
	jwtManager        *mocks.JWTManager
	usecase           ebank.UserServiceServer
}

func (ts *UserUsecaseTestSuite) SetupTest() {
	ts.userRepository = new(mocks.UserRepository)
	ts.accountRepository = new(mocks.AccountRepository)
	ts.userHelper = new(mocks.UserHelper)
	ts.jwtManager = new(mocks.JWTManager)
	ts.usecase = NewUserService(ts.userHelper, ts.userRepository, ts.accountRepository, ts.jwtManager)
}

func (ts *UserUsecaseTestSuite) Test_userService_GetUser() {
	testUser := &ebank.User{
		Id:          1,
		Name:        "Test User",
		Birth:       "2000-01-01",
		PhoneNumber: "010-1234-5678",
		// Accounts:    []*ebank.Account{},
	}

	ts.userRepository.EXPECT().GetUserByID(mock.Anything, testUser.GetId()).
		Return(&model.User{
			ID:          testUser.Id,
			Name:        testUser.Name,
			Birth:       testUser.Birth,
			PhoneNumber: testUser.PhoneNumber,
		}, nil)

	ts.accountRepository.EXPECT().GetAccountsByUserID(mock.Anything, testUser.GetId()).
		Return([]model.Account{}, nil)

	req := &ebank.GetUserRequest{
		Id: testUser.Id,
	}

	resp, err := ts.usecase.GetUser(context.Background(), req)

	ts.Equal(testUser, resp.User)
	ts.NoErrorf(err, "error should be nil")
}
