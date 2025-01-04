package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"ebank/api/v1"
	"ebank/pkg/config"
	"ebank/services/transaction/repository"
	transactionService "ebank/services/transaction/service"
)

var (
	// command-line options:
	// gRPC server endpoint
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:9090", "gRPC server endpoint")
)

func main() {
	cfg := config.New()
	lis, err := net.Listen("tcp", cfg.Server.Port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	logrusEntry := logrus.NewEntry(logrus.StandardLogger())

	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_logrus.UnaryServerInterceptor(logrusEntry),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)

	transactionRepository, err := repository.NewTransactionFileRepository(cfg.DB.TransactionTablePath)
	if err != nil {
		log.Fatalf("failed to make transactionRepository: %v", err)
	}

	transactionService := transactionService.NewTransactionService(transactionRepository)

	ebank.RegisterTransactionServiceServer(s, transactionService)

	mux := runtime.NewServeMux()

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	http.ListenAndServe(":8081", mux)

	fmt.Println("Server is running on " + cfg.Server.Port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
