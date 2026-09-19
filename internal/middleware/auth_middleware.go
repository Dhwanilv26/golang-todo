package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"todo-api/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if tokenString == "" || tokenString == authHeader { // checked for no token, or tokenString= authheader is used for authorization = something else instead of bearer
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims) // same as "token.claims as jwt.mapclaims" (type assertion like TS)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}

		userId, ok := claims["user_id"].(string) // jwt.mapclaims is a map, so key ki value extract karne ke liye claims["user_id"] which is parsed into a string

		if exp, ok := claims["exp"].(float64); ok {
			// pass second and miliseconds in time.unix to get time.Time value for expirationTime
			expirationTime := time.Unix(int64(exp), 0)

			if time.Now().After(expirationTime) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "token has expired"})
				c.Abort()
				return
			}
		}

		c.Set("user_id", userId) // settting the user_id here, in the context variable, so that other routes can use it
		c.Next()
	}
}
