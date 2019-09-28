package signin

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/src/jwtauth"
)

type Credentials struct {
	email    string
	password string
}

const MINUTE int = 60 * 1000

func checkValidCredentials(credentials Credentials, db *sql.DB) (int, string) {
	if credentials.email == "" || credentials.password == "" {
		return 0, "In order to sign in provide as form values email and password"
	}

	row := db.QueryRow(`
		SELECT id, password FROM users  
		WHERE email = $1;
	`, credentials.email)

	var userId int
	var expectedPassword string
	if err := row.Scan(&userId, &expectedPassword); err != nil {
		return 0, "There is no user with the specifyed email"
	}

	if expectedPassword == credentials.password {
		return userId, "User authorized!"
	} else {
		return 0, "Password does not match users password"
	}
}

func SignIn(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var credentials Credentials = Credentials{}

		credentials.email = c.PostForm("email")
		credentials.password = c.PostForm("password")

		if userId, feedback := checkValidCredentials(credentials, db); userId == 0 {
			c.String(http.StatusUnauthorized, feedback)
			return
		} else {
			generatedJwt, err := jwtauth.GetUserJwt(credentials.email)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			} else {
				c.SetCookie("userId", strconv.Itoa(userId), 10*MINUTE, "/", "", true, true)
				c.SetCookie("jwt", generatedJwt, 10*MINUTE, "/", "", true, true)
				c.JSON(http.StatusOK, gin.H{
					"message": "User successfully signed in! Use this jwt for requests authentication",
					"email":   credentials.email,
					"jwt":     generatedJwt,
				})
				return
			}
		}
	}
}
