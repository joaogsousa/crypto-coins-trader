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
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
	"github.com/russross/blackfriday"
)

/*
func dbFunc(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := db.Exec("CREATE TABLE IF NOT EXISTS ticks (tick timestamp)"); err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error creating database table: %q", err))
			return
		}

		if _, err := db.Exec("INSERT INTO ticks VALUES (now())"); err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error incrementing tick: %q", err))
			return
		}

		rows, err := db.Query("SELECT tick FROM ticks")
		if err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error reading ticks: %q", err))
			return
		}

		defer rows.Close()
		for rows.Next() {
			var tick time.Time
			if err := rows.Scan(&tick); err != nil {
				c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error scanning ticks: %q", err))
				return
			}
			c.String(http.StatusOK, fmt.Sprintf("Read from DB: %s\n", tick.String()))
		}
	}
} */

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
		c.String(http.StatusOK, string(blackfriday.Run([]byte("**hi!**"))))
	})

	router.POST("/users/register", register.NewUser(db))

	router.Run(":" + port)
}
