package api

import (
	"net/http"
	"time"
	"upload-service/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtKey            = ("my_secret_key")
	DefaultExpireTime = time.Now().Add(time.Minute * 5)
)

// this is only temporary as we are not using a database for users
var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

type Credentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (h *handler) Login(c *gin.Context) {
	var creds Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		errors.HandleBadRequest(c, err)
	}

	expectedPassword, ok := users[creds.Username]
	if !ok || expectedPassword != creds.Password {
		errors.HandleUnauthorized(c, errors.ErrInvalidCredentials)
		return
	}
	accesToken, err := jwtGenerate(creds.Username, jwtKey, time.Until(DefaultExpireTime))
	if err != nil {
		errors.HandleInternalServerError(c, err)
		return
	}
	c.SetCookie("access_token", accesToken, int(time.Until(DefaultExpireTime).Seconds()), "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "login successful", "accessToken": accesToken})
}

func jwtGenerate(userName, secret string, expiresAt time.Duration) (string, error) {
	claims := &Claims{
		Username: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "upload-service",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresAt)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(secret))
}
