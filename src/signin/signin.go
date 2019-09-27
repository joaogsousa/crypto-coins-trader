package signin

import (
	"database/sql"
	"fmt"
	"net/http"
	_ "strconv"

	"github.com/gin-gonic/gin"
)

type Credentials struct {
	email    string
	password string
}

func checkValidCredentials(credentials Credentials, db *sql.DB) bool {
	row := db.QueryRow(`
		SELECT password FROM users  
		WHERE email = $1;
	`, credentials.email)

	var expectedPassword string
	if err := row.Scan(&expectedPassword); err != nil {
		fmt.Println("Error occured while fetching password from database")
		return false
	}

	if expectedPassword == credentials.password {
		return true
	} else {
		fmt.Println("Password does not match users password")
		return false
	}
}

func SignIn(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var credentials Credentials = Credentials{}

		credentials.email = c.PostForm("email")
		credentials.password = c.PostForm("password")

		if !checkValidCredentials(credentials, db) {
			c.String(http.StatusUnauthorized, "Unautorized. Check if you sent correct credentials on POST form: email, password.")
			return
		}
		c.String(http.StatusOK, "User authorized!")
	}
}
