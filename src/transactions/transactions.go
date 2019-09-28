package transactions

import (
	"database/sql"
	"net/http"
	_ "strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/src/jwtauth"
)

func IsAuthorized(c *gin.Context) bool {
	tokenStr, err := c.Cookie("jwt")

	if err != nil {
		c.String(http.StatusUnauthorized, "Unautorized. You must provide a JWT token to access this route. -> Sign in first...")
		return false
	}

	claims := &jwtauth.Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtauth.SecretKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.String(http.StatusUnauthorized, "Invalid signature for jwt token")
			return false
		}
		c.String(http.StatusBadRequest, "Bad request...")
		return false
	}
	if !token.Valid {
		c.String(http.StatusUnauthorized, "Invalid jwt token")
		return false
	}

	return true
}

func Operation(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		//operationType := c.Param("operation")

		if IsAuthorized(c) {
			c.String(http.StatusOK, "Authorized route!")
		}
	}
}
