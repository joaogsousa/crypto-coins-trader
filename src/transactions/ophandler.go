package transactions

import (
	"database/sql"
	"net/http"
	_ "strconv"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/src/jwtauth"
)

func OperationHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		operationType := c.Param("operation")
		if operationType != "buy" && operationType != "sell" {
			c.String(
				http.StatusBadRequest,
				"Bad request. Please indicate operation type. Either /transactions/buy or /transactions/sell",
			)
			return
		}

		if jwtauth.IsAuthorized(c) {
			Operation(db, c, operationType)
		}
	}
}
