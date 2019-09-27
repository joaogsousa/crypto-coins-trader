package register

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	name      string
	email     string
	password  string
	birthdate string
	cash      float64
	coins     int
}

func checkValidUser(user User) bool {
	if user.name != "" && user.email != "" && user.password != "" && user.birthdate != "" {
		return true
	} else {
		return false
	}
}

func NewUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user User = User{}
		//c.String(200, "entered user register")
		user.name = c.PostForm("name")
		user.email = c.PostForm("email")
		user.password = c.PostForm("password")
		user.birthdate = c.PostForm("birthdate")
		user.cash = c.PostForm("cash")
		user.coins = c.PostForm("coins")

		if !checkValidUser(user) {
			c.String(http.StatusBadRequest, "User data was not correctly provided. Send: name,email,password and bithdate on POST form")
		}

		_, err := db.Exec(`
			INSERT INTO users (name, email, password, birthdate, cash, coins)  
			VALUES ($1, $2, $3, $4, $5, $6);
		`, user.name, user.email, user.password, user.birthdate, user.cash, user.coins)
		if err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error inserting new user: %v", err))
			return
		}

		c.String(http.StatusOK, "User registered!")
		fmt.Println("User registered: ", user)
	}
}
