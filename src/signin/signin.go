package signin

import (
	"database/sql"
	"net/http"
	_ "strconv"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/src/jwtauth"
)

type Credentials struct {
	email    string
	password string
}

const HOUR int = 60 * 60 * 1000

func checkValidCredentials(credentials Credentials, db *sql.DB) (bool, string) {
	if credentials.email == "" || credentials.password == "" {
		return false, "In order to sign in provide as form values email and password"
	}

	row := db.QueryRow(`
		SELECT password FROM users  
		WHERE email = $1;
	`, credentials.email)

	var expectedPassword string
	if err := row.Scan(&expectedPassword); err != nil {
		return false, "There is no user with the specifyed email"
	}

	if expectedPassword == credentials.password {
		return true, "User authorized!"
	} else {
		return false, "Password does not match users password"
	}
}

func SignIn(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var credentials Credentials = Credentials{}

		credentials.email = c.PostForm("email")
		credentials.password = c.PostForm("password")

		if ok, feedback := checkValidCredentials(credentials, db); !ok {
			c.String(http.StatusUnauthorized, feedback)
			return
		} else {
			generatedJwt, err := jwtauth.GetUserJwt(credentials.email)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			} else {
				c.SetCookie("jwt", generatedJwt, 1*HOUR, "/", "", true, true)
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
