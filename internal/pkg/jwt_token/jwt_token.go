package jwt_token

import (
	"github.com/ZhdanovichVlad/go_final_project/internal"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func GenerateToken(userUuid string) (string, error) {
	now := time.Now().UTC()
	claims := jwt.MapClaims{
		"userUuid": userUuid,
		"exp":      now.Add(internal.TokenLifetimeInHours * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accesstoken, err := token.SignedString([]byte(internal.SecretKey))
	if err != nil {
		return "", err
	}

	return accesstoken, nil
}
