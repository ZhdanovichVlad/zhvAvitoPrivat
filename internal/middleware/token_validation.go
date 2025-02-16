package middleware

import (
	"fmt"
	"github.com/ZhdanovichVlad/go_final_project/internal/pkg/jwttoken"
	"net/http"
	"strings"
	"time"

	"github.com/ZhdanovichVlad/go_final_project/internal/pkg/errorsx"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errorsx.ErrAuthHeaderIsEmpty.Error()})
			return
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errorsx.ErrUnexpSignedMetod
			}
			return []byte(jwttoken.SecretKey), nil
		})
		if err != nil {
			fmt.Println("test 1")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errorsx.ErrInvalidToken.Error()})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			exp, ok := claims["exp"].(float64)
			if !ok || time.Now().Unix() > int64(exp) {
				fmt.Println("test 2")
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errorsx.ErrTokenExpired.Error()})
				return
			}

			userUUID, ok := claims["userUUID"].(string)
			if !ok || userUUID == "" {
				fmt.Println("test 3")
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errorsx.ErrInvUserUUIDInToken.Error()})
				return
			}

			c.Set("userUUID", userUUID)
			c.Next()
		} else {
			fmt.Println("test 4")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errorsx.ErrInvalidToken.Error()})
			return
		}
	}
}
