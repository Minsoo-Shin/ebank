package main

import (
	pb "ebank/internal/api/v1"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"log"
	"net"

	"ebank/internal/repository/repository_impl"
	"ebank/internal/service"
	"ebank/pkg/config"
	"ebank/pkg/jwt_manager"
)

func main() {
	cfg := config.New()
	lis, err := net.Listen("tcp", cfg.Server.Port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	logrusEntry := logrus.NewEntry(logrus.StandardLogger())

	jwtManager := jwt_manager.NewJWTManager(cfg.Jwt.SecretKey, cfg.Jwt.Duration)
	interceptor := service.NewUserInterceptor(jwtManager)

	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_logrus.UnaryServerInterceptor(logrusEntry),
			grpc_recovery.UnaryServerInterceptor(),
			interceptor.Unary(),
		)),
	)

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

	userHelper := service.NewUserHelper(userFileRepository)
	userService := service.NewUserService(userHelper, userFileRepository, accountRepository, jwtManager)
	accountService := service.NewAccountService(userHelper, accountRepository, transactionRepository)

	pb.RegisterUserServiceServer(s, userService)
	pb.RegisterAccountServiceServer(s, accountService)

	fmt.Println("Server is running on " + cfg.Server.Port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
