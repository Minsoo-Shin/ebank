package service

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"ebank/pkg/jwt_manager"
)

type UserInterceptor struct {
	jwtManager jwt_manager.JWTManager
}

func NewUserInterceptor(jwtManager jwt_manager.JWTManager) *UserInterceptor {
	return &UserInterceptor{jwtManager}
}

func (interceptor *UserInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		var err error
		if !interceptor.skipper(info.FullMethod) {
			ctx, err = interceptor.authorize(ctx)
			if err != nil {
				return nil, err
			}
		}

		return handler(ctx, req)
	}
}

func (interceptor *UserInterceptor) authorize(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return ctx, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	if strings.Contains(accessToken, "Bearer ") {
		accessToken = strings.Split(accessToken, "Bearer ")[1]
	}
	claims, err := interceptor.jwtManager.Verify(accessToken)
	if err != nil {
		return ctx, status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}
	if claims == nil {
		return ctx, status.Error(codes.Unauthenticated, "access token is invalid")
	}

	return context.WithValue(ctx, "user", claims), nil
}

func (interceptor *UserInterceptor) skipper(method string) bool {
	switch method {
	case "/proto.UserService/Login":
		return true
	case "/proto.UserService/CreateUser":
		return true
	default:
		return false
	}
}
