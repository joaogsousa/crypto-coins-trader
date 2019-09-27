package register

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
)

type User struct {
	name      string
	email     string
	password  string
	birthdate string
}

func NewUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(200, "entered user register")
		fmt.Println("c.PostForm", c.PostForm)
	}
}
