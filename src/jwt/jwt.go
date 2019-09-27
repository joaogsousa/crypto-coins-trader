package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	email string
	jwt.StandardClaims
}

func GetUserJwt(email string) (string, error) {

	var jwtKey = []byte("very_strong_and_roboust_secret_key")

	expirationTime := time.Now().Add(10 * time.Minute)

	claims := &Claims{
		email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", errors.New("Error generating signed jwt")
	}

	return tokenString, nil
}
