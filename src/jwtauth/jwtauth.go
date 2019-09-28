package jwtauth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	email string
	jwt.StandardClaims
}

var SecretKey = []byte("very_strong_and_roboust_secret_key")

func GetUserJwt(email string) (string, error) {
	expirationTime := time.Now().Add(10 * time.Minute)

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
