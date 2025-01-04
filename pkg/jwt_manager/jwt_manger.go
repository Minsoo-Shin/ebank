package jwt_manager

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	"ebank/services/user/model"
)

type JWTManager interface {
	Generate(user model.User) (string, error)
	Verify(accessToken string) (*UserClaims, error)
}

type jwtManager struct {
	secretKey        string
	tokenDuration    time.Duration
	userJwtExpiryMap map[string]int64 // redis 서버를 이용하여 저장
}

type UserClaims struct {
	jwt.StandardClaims
	PhoneNumber string `json:"username"`
}

func NewJWTManager(secretKey string, tokenDuration time.Duration) JWTManager {
	return &jwtManager{
		secretKey:        secretKey,
		tokenDuration:    tokenDuration,
		userJwtExpiryMap: make(map[string]int64),
	}
}

func (manager *jwtManager) Generate(user model.User) (string, error) {
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(manager.tokenDuration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		PhoneNumber: user.PhoneNumber,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.secretKey))
}

func (manager *jwtManager) Verify(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return []byte(manager.secretKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	// 서버에서 강제로 토큰을 만료하고 싶을 때 활용
	// if jwtExpiryDate, ok := manager.userJwtExpiryMap[claims.PhoneNumber]; ok && claims.IssuedAt < jwtExpiryDate {
	//	// 발급한 토큰이 서버 지정 만료시간보다 더 이전 토큰이라면 만료된 토큰으로 취급
	//	return nil, fmt.Errorf("invalid token claims - expired token")
	// }

	return claims, nil
}
