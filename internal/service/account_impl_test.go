package service

import (
	"context"
	pb "ebank/internal/api/v1"
	"ebank/internal/model"
	"ebank/internal/repository/repository_impl"
	"ebank/mocks"
	"ebank/pkg/config"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestAccountUsecaseSuite(t *testing.T) {
	suite.Run(t, new(AccountUsecaseTestSuite))
}

var (
	_ suite.SetupTestSuite = &AccountUsecaseTestSuite{}
)

type AccountUsecaseTestSuite struct {
	suite.Suite
	userRepository        *mocks.UserRepository
	accountRepository     *mocks.AccountRepository
	transactionRepository *mocks.TransactionRepository
	userHelper            *mocks.UserHelper
	usecase               pb.AccountServiceServer
}

func (ts *AccountUsecaseTestSuite) SetupTest() {
	ts.userRepository = new(mocks.UserRepository)
	ts.accountRepository = new(mocks.AccountRepository)
	ts.transactionRepository = new(mocks.TransactionRepository)
	ts.userHelper = new(mocks.UserHelper)
	ts.usecase = NewAccountService(ts.userHelper, ts.accountRepository, ts.transactionRepository)
}

func (ts *AccountUsecaseTestSuite) Test_accountService_CreateAccount() {
	testAccount := model.Account{
		ID:            1,
		AccountNumber: "1234567890",
		CustomerID:    123,
		Balance:       1000,
	}

	ts.accountRepository.EXPECT().CreateAccount(mock.Anything, mock.Anything).Return(testAccount, nil)

	account, err := ts.usecase.(*accountService).CreateAccount(context.Background(), &pb.CreateAccountRequest{
		AccountNumber: testAccount.AccountNumber,
		UserId:        testAccount.CustomerID,
	})

	ts.NoError(err)
	ts.Equal(testAccount.ID, account.Account.Id)
	ts.Equal(testAccount.AccountNumber, account.Account.AccountNumber)
	ts.Equal(testAccount.CustomerID, account.Account.CustomerId)
	ts.Equal(testAccount.Balance, account.Account.Balance)
}

func (ts *AccountUsecaseTestSuite) Test_accountService_GetAccount() {
	testAccount := &model.Account{
		ID:            1,
		AccountNumber: "1234567890",
		CustomerID:    123,
		Balance:       1000,
	}

	ts.accountRepository.EXPECT().GetAccountByID(mock.Anything, mock.Anything).Return(testAccount, nil)
	ts.userHelper.EXPECT().ValidateUser(mock.Anything, mock.Anything).Return(model.User{ID: 1}, nil)

	account, err := ts.usecase.(*accountService).GetAccount(context.Background(), &pb.GetAccountRequest{
		Id: testAccount.ID,
	})

	ts.NoError(err)
	ts.Equal(testAccount.ID, account.Account.Id)
	ts.Equal(testAccount.AccountNumber, account.Account.AccountNumber)
	ts.Equal(testAccount.CustomerID, account.Account.CustomerId)
	ts.Equal(testAccount.Balance, account.Account.Balance)
}

func (ts *AccountUsecaseTestSuite) Test_accountService_UpdateAccount() {
	testAccount := model.Account{
		ID:            1,
		AccountNumber: "1234567890",
		CustomerID:    123,
		Balance:       1000,
	}

	newAccount := model.Account{
		ID:            1,
		AccountNumber: "234",
		CustomerID:    123,
		Balance:       1000,
	}

	ts.accountRepository.EXPECT().GetAccountByID(mock.Anything, mock.Anything).Return(&testAccount, nil)
	ts.userHelper.EXPECT().ValidateUser(mock.Anything, mock.Anything).Return(model.User{ID: 1}, nil)

	ts.accountRepository.EXPECT().UpdateAccount(mock.Anything, mock.Anything).Return(nil)

	account, err := ts.usecase.(*accountService).UpdateAccount(context.Background(), &pb.UpdateAccountRequest{
		Id:            testAccount.ID,
		AccountNumber: newAccount.AccountNumber,
	})

	ts.NoError(err)
	ts.Equal(newAccount.ID, account.Account.Id)
	ts.Equal(newAccount.AccountNumber, account.Account.AccountNumber)
	ts.Equal(newAccount.CustomerID, account.Account.CustomerId)
	ts.Equal(testAccount.Balance, account.Account.Balance)
}

func (ts *AccountUsecaseTestSuite) Test_accountService_DeleteAccount() {
	testAccount := model.Account{
		ID:            1,
		AccountNumber: "234",
		CustomerID:    123,
		Balance:       1000,
	}

	ts.accountRepository.EXPECT().GetAccountByID(mock.Anything, mock.Anything).Return(&testAccount, nil)
	ts.userHelper.EXPECT().ValidateUser(mock.Anything, mock.Anything).Return(model.User{ID: 1}, nil)
	ts.accountRepository.EXPECT().DeleteAccount(mock.Anything, mock.Anything).Return(nil)

	_, err := ts.usecase.(*accountService).DeleteAccount(context.Background(), &pb.DeleteAccountRequest{
		Id: 1,
	})

	ts.NoError(err)
}

func (ts *AccountUsecaseTestSuite) Test_accountService_GetAllAccounts() {
	testAccounts := []model.Account{
		{
			ID:            1,
			AccountNumber: "1234567890",
			CustomerID:    123,
			Balance:       1000,
		},
		{
			ID:            2,
			AccountNumber: "1234567891",
			CustomerID:    123,
			Balance:       2000,
		},
	}

	ts.accountRepository.EXPECT().GetAllAccounts(mock.Anything).Return(testAccounts, nil)

}

type TestServices struct {
	accountService pb.AccountServiceServer
	userService    pb.UserServiceServer
}

func testServiceGenerator(t *testing.T, dir string) TestServices {
	_ = os.Mkdir(dir, os.ModePerm)
	cfg := config.New()
	cfg.DB.AccountTablePath = dir + "/account_test.json"
	cfg.DB.UserTablePath = dir + "/user_test.json"
	cfg.DB.TransactionTablePath = dir + "/transaction_test.json"

	userFileRepository, err := repository_impl.NewUserFileRepository(cfg.DB.UserTablePath)
	if err != nil {
		log.Fatalf("failed to make userFileRepository: %v", err)
	}

	accountRepository, err := repository_impl.NewAccountFileRepository(cfg.DB.AccountTablePath)
	if err != nil {
		log.Fatalf("failed to make accountRepository: %v", err)
	}

	transactionRepository, err := repository_impl.NewTransactionFileRepository(cfg.DB.TransactionTablePath)
	if err != nil {
		log.Fatalf("failed to make transactionRepository: %v", err)
	}

	userHelper := NewUserHelper(userFileRepository)
	userService := NewUserService(userHelper, userFileRepository, accountRepository, nil)
	accountService := NewAccountService(userHelper, accountRepository, transactionRepository)

	return TestServices{
		accountService: accountService,
		userService:    userService,
	}
}

func initDataset(t *testing.T, services TestServices) {
	_, err := services.userService.CreateUser(context.Background(), &pb.CreateUserRequest{
		PhoneNumber: "123123123",
		Name:        "test",
		Birth:       "1990-01-01",
		Password:    "123",
	})
	if err != nil {
		t.Fatalf("Failed to create account: %v", err)
	}

	testAccount := model.Account{
		ID:            1,
		AccountNumber: "1234567890",
		CustomerID:    1,
		Balance:       0,
	}

	_, err = services.accountService.CreateAccount(context.Background(), &pb.CreateAccountRequest{
		AccountNumber: testAccount.AccountNumber,
		UserId:        testAccount.CustomerID,
	})
	if err != nil {
		t.Fatalf("Failed to create account: %v", err)
	}
}

// 함수는 지정된 디렉토리의 모든 내용을 삭제합니다.
func clearDirectory(dir string) error {
	// 디렉토리 읽기
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()

	// 디렉토리 내용 얻기
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}

	// 각 항목 삭제
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}

	err = os.Remove(dir)
	if err != nil {
		return err
	}

	return nil
}

func Test_accountService_DepositAndWithdrawConcurrently(t *testing.T) {
	path := "tmp"
	services := testServiceGenerator(t, path)
	defer clearDirectory(path)

	_, err := services.userService.CreateUser(context.Background(), &pb.CreateUserRequest{
		PhoneNumber: "01011112222",
		Name:        "test",
		Birth:       "1990-01-01",
		Password:    "123",
	})
	if err != nil {
		t.Fatalf("Failed to create account: %v", err)
	}

	testAccount := model.Account{
		ID:            1,
		AccountNumber: "1234567890",
		CustomerID:    1,
		Balance:       0,
	}

	_, err = services.accountService.CreateAccount(context.Background(), &pb.CreateAccountRequest{
		AccountNumber: testAccount.AccountNumber,
		UserId:        testAccount.CustomerID,
	})
	if err != nil {
		t.Fatalf("Failed to create account: %v", err)
	}

	const amt = 10
	const c = 1000
	var negBal int32
	var start, g sync.WaitGroup
	start.Add(1)
	g.Add(3 * c)
	for i := 0; i < c; i++ {
		go func() { // deposit
			start.Wait()
			services.accountService.Deposit(context.TODO(), &pb.DepositRequest{
				AccountId: testAccount.ID,
				Amount:    amt,
			}) // ignore return values
			g.Done()
		}()
		go func() { // withdraw
			start.Wait()
			for {
				_, err := services.accountService.Withdraw(context.TODO(), &pb.WithdrawRequest{
					AccountId: testAccount.ID,
					Amount:    amt,
				})
				if err == nil {
					break
				}
				//log.Println(err)
				time.Sleep(1 * time.Millisecond)
			}

			g.Done()
		}()
		go func() { // watch that balance stays >= 0
			start.Wait()
			if p, _ := services.accountService.GetAccount(context.TODO(), &pb.GetAccountRequest{
				Id: testAccount.ID,
			}); p.GetAccount().GetBalance() < 0 {
				atomic.StoreInt32(&negBal, 1)
			}
			g.Done()
		}()
	}
	start.Done()
	g.Wait()
	if negBal == 1 {
		t.Fatal("Balance went negative with concurrent deposits and " +
			"withdrawals.  Want balance always >= 0.")
	}
	if p, err := services.accountService.GetAccount(context.TODO(), &pb.GetAccountRequest{Id: testAccount.ID}); err != nil || p.GetAccount().Balance != 0 {
		t.Fatalf("After equal concurrent deposits and withdrawals, a.Balance = %v, %v.  Want 0, true", strconv.Itoa(int(p.GetAccount().GetBalance())), err)
	}

}
