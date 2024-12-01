package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func extractTokenFromHeader(header string) (string, error) {
	// The Authorization header is in the format `Bearer <token>`
	if header == "" {
		return "", errors.New("missing token")
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 {
		return "", errors.New("incorrectly formatted authorization header")
	}

	return jwtToken[1], nil
}
func MiddlewareCheckLoginJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// since the token s been set to cookie, we need to double check.
		// 1. Check the Authorization header
		tokenString, err := extractTokenFromHeader(c.GetHeader("Authorization"))
		if err != nil {
			if err.Error() != "missing token" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				c.Abort()
				return
			} else {
				cookie, err := c.Cookie("access_token")
				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
					c.Abort()
					return
				}
				tokenString = cookie
			}
		}
		// 3. Validate the token
		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 4. Add claims to the context for downstream handlers
		c.Set("claims", claims)
		c.Next()
	}
}
