package middleware

import (
	"github.com/ZhdanovichVlad/go_final_project/internal"
	"github.com/ZhdanovichVlad/go_final_project/internal/pkg/errorsx"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

//func ValidateToken() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		var header = c.Request.Header["Authorization"]
//		if header == nil {
//			c.AbortWithStatusJSON(http.StatusUnauthorized, entity.ErrorMsg{errorsx.Unauthorized.Error())
//			return
//		}
//		var token = header[0]
//		token = strings.Replace(token, "Bearer ", "", 1)
//
//		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//				// Проверяем алгоритм подписи
//				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
//				}
//				// Возвращаем секретный ключ для проверки подписи
//				return []byte(secretKey), nil
//			})
//
//			if err != nil {
//				return nil, err
//			}
//
//			// Проверяем, валиден ли токен
//			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
//				return claims, nil
//			} else {
//				return nil, fmt.Errorf("invalid token")
//			}
//
//		c.Next()
//	}
//}
//}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errorsx.AuthHeaderIsEmpty.Error()})
			return
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errorsx.UnexpSignedMetod
			}
			return []byte(internal.SecretKey), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

			exp, ok := claims["exp"].(float64)
			if !ok || time.Now().Unix() > int64(exp) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
				return
			}

			userUuid, ok := claims["userUuid"].(string)
			if !ok || userUuid == "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errorsx.InvUserUuidInToken.Error()})
				return
			}

			c.Set("userUuid", userUuid)
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errorsx.InvalidToken.Error()})
		}
	}
}
