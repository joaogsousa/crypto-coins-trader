package transactions

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BuyOperation struct {
	sellingUserId string
	coinsAmount   int
}

func checkAvailableCredit(userId string, operationCost float64, db *sql.DB) (bool, string) {
	row := db.QueryRow(`
		SELECT cash FROM users  
		WHERE id = $1;
	`, userId)

	var cash float64
	if err := row.Scan(&cash); err != nil {
		return false, "There is no user with the specifyed id"
	}

	if operationCost > cash {
		return false, "Operation denied, not enough cash for the operation"
	} else {
		return true, ""
	}
}

func checkAvailableCoins(userId string, coinsAmount int, db *sql.DB) (bool, string) {
	row := db.QueryRow(`
		SELECT coins FROM users  
		WHERE id = $1;
	`, userId)

	var coins int
	if err := row.Scan(&coins); err != nil {
		return false, "There is no user with the specifyed id"
	}

	if coinsAmount > coins {
		return false, fmt.Sprintf("Operation denied, user %v does not have the requested amount of coins", userId)
	} else {
		return true, ""
	}
}

func Buy(db *sql.DB, c *gin.Context) {
	// TODO receber tambem a DATA da operacao!!!!!
	var buyOperation BuyOperation = BuyOperation{}

	if c.PostForm("sellingUserId") == "" || c.PostForm("coinsAmount") == "" {
		c.String(http.StatusBadRequest, "Bad request. Provide sellingUserId and coinsAmount as POST form values")
	}

	buyOperation.sellingUserId = c.PostForm("sellingUserId")

	coinsAmount, _ := strconv.ParseInt(c.PostForm("coinsAmount"), 10, 64)
	buyOperation.coinsAmount = int(coinsAmount)

	userId, err := c.Cookie("userId")
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal server error... Coundnt retrieve user ID")
		return
	}

	ok, feedback := checkAvailableCredit(userId, float64(buyOperation.coinsAmount*CoinPrice), db)
	if !ok {
		c.String(http.StatusMethodNotAllowed, feedback)
		return
	}

	ok, feedback = checkAvailableCoins(buyOperation.sellingUserId, buyOperation.coinsAmount, db)
	if !ok {
		c.String(http.StatusMethodNotAllowed, feedback)
		return
	}

	c.String(http.StatusOK, "Operation permitted")
}
