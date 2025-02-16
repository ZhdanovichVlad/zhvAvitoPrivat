package jwttoken

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	SecretKey            = "jwtsecret"
	TokenLifetimeInHours = 12
)

type JwtTokenGenerator struct {
}

func NewJwtTokenGenerator() *JwtTokenGenerator {
	return &JwtTokenGenerator{}
}

func (j *JwtTokenGenerator) GenerateToken(userUUID string) (string, error) {
	now := time.Now().UTC()
	claims := jwt.MapClaims{
		"userUUID": userUUID,
		"exp":      now.Add(TokenLifetimeInHours * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accesstoken, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	return accesstoken, nil
}
