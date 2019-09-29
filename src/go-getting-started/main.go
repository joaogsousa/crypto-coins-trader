package main

import (
	_ "bytes"
	"database/sql"
	"log"
	"net/http"
	"os"
	_ "strconv"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/src/register"
	"github.com/heroku/go-getting-started/src/signin"
	"github.com/heroku/go-getting-started/src/transactions"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.Default()

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "RedCoins crypto-coins trader")
	})

	router.POST("/users/register", register.NewUser(db))
	router.POST("/users/signin", signin.SignIn(db))
	router.POST("/transactions/:operation", transactions.OperationHandler(db))

	router.Run(":" + port)
}
