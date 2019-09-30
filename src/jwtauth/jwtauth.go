package jwtauth

import (
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Claims struct {
	email string
	jwt.StandardClaims
}

var SecretKey = []byte("very_strong_and_roboust_secret_key")

func GetUserJwt(email string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &Claims{
		email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", errors.New("Error generating signed jwt")
	}

	return tokenString, nil
}

func IsAuthorized(c *gin.Context) bool {
	tokenStr, err := c.Cookie("jwt")

	if err != nil {
		c.String(http.StatusUnauthorized, "Unautorized. You must provide a JWT token to access this route. -> Sign in first...")
		return false
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if err != nil {
		c.String(http.StatusUnauthorized, "Invalid signature for jwt token")
		return false
	}
	if !token.Valid {
		c.String(http.StatusUnauthorized, "Invalid jwt token")
		return false
	}

	return true
}
