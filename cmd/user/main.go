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
	"ebank/pkg/jwt_manager"
	"ebank/services/user/repository"
	authService "ebank/services/user/service"
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

	jwtManager := jwt_manager.NewJWTManager(cfg.Jwt.SecretKey, cfg.Jwt.Duration)
	interceptor := authService.NewUserInterceptor(jwtManager)

	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_logrus.UnaryServerInterceptor(logrusEntry),
			grpc_recovery.UnaryServerInterceptor(),
			interceptor.Unary(),
		)),
	)

	userFileRepository, err := repository.NewUserFileRepository(cfg.DB.UserTablePath)
	if err != nil {
		log.Fatalf("failed to make userFileRepository: %v", err)
	}

	userHelper := authService.NewUserHelper(userFileRepository)
	userService := authService.NewUserService(userHelper, userFileRepository, jwtManager)
	authService := authService.NewAuthService(userFileRepository, jwtManager)

	ebank.RegisterUserServiceServer(s, userService)
	ebank.RegisterAuthServiceServer(s, authService)

	mux := runtime.NewServeMux()
	// opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	// err = ebank.RegisterAccountServiceHandlerFromEndpoint(context.TODO(), mux, *grpcServerEndpoint, opts)
	// if err != nil {
	// 	log.Fatalf("Failed to register handler: %v", err)
	// }

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	http.ListenAndServe(":8081", mux)

	fmt.Println("Server is running on " + cfg.Server.Port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
