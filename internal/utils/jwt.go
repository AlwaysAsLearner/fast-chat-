package utils

import (
	"time"

	config "github.com/AlwaysAsLearner/fast-chat/backend/configs"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userId uint) (string, error) {
	jwtConfig := config.JWT
	expirationTime := time.Now().Add(time.Duration(jwtConfig.Expiration) * time.Minute)
	claims := &Claims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtConfig.Secret))
	return signedToken, err
}

func ValidateToken(tokenString string) (*Claims, error) {
	jwtConfig := config.JWT
	claims := &Claims{}
	// This function takes token, claims ( to write data ), function which returns secret key ( token signed with )
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtConfig.Secret), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	// you can also return token.Claims
	return claims, nil
}
