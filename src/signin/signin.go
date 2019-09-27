package signin

import (
	"database/sql"
	"net/http"
	_ "strconv"

	"github.com/gin-gonic/gin"
)

type Credentials struct {
	email    string
	password string
}

func checkValidCredentials(credentials Credentials, db *sql.DB) (bool, string) {
	if credentials.email == "" || credentials.password == "" {
		return false
		return "In order to sign in provide as form values email and password"
	}

	row := db.QueryRow(`
		SELECT password FROM users  
		WHERE email = $1;
	`, credentials.email)

	var expectedPassword string
	if err := row.Scan(&expectedPassword); err != nil {
		return false
		return "There is no user with the specifyed email"
	}

	if expectedPassword == credentials.password {
		return true
		return "User authorized!"
	} else {
		return false
		return "Password does not match users password"
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
			c.String(http.StatusOK, feedback)
		}
	}
}
