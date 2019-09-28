package transactions

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TradeInfo struct {
	buyingUserId  string
	sellingUserId string
	coinsAmount   int
	operationCost float64
	date          string
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

func Operation(db *sql.DB, c *gin.Context, operationType string) {
	// TODO receber tambem a DATA da operacao!!!!!
	if c.PostForm("userId") == "" || c.PostForm("coinsAmount") == "" || c.PostForm("date") == "" {
		c.String(http.StatusBadRequest, "Bad request. Provide serId, coinsAmount and date as POST form values")
		return
	}

	operationEndUserId := c.PostForm("userId")
	loggedUserId, err := c.Cookie("userId")
	if err != nil {
		c.String(http.StatusInternalServerError, "Coundnt retrieve ID for the logged user")
		return
	}

	coinsAmount, _ := strconv.Atoi(c.PostForm("coinsAmount"))

	var buyingUserId, sellingUserId string
	if operationType == "buy" {
		buyingUserId = loggedUserId
		sellingUserId = operationEndUserId
	} else { // sell operation
		buyingUserId = operationEndUserId
		sellingUserId = loggedUserId
	}

	tradeInfo := TradeInfo{
		buyingUserId:  buyingUserId,
		sellingUserId: sellingUserId,
		coinsAmount:   coinsAmount,
		operationCost: float64(coinsAmount * CoinPrice),
		date:          c.PostForm("date"),
	}

	if tradeInfo.buyingUserId == tradeInfo.sellingUserId {
		c.String(http.StatusBadRequest, "Invalid operation. buyingUserId and sellingUserId are the same.")
		return
	}

	ok, feedback := checkAvailableCredit(tradeInfo.buyingUserId, tradeInfo.operationCost, db)
	if !ok {
		c.String(http.StatusMethodNotAllowed, feedback)
		return
	}

	ok, feedback = checkAvailableCoins(tradeInfo.sellingUserId, tradeInfo.coinsAmount, db)
	if !ok {
		c.String(http.StatusMethodNotAllowed, feedback)
		return
	}

	// operation permited, proceed...
	err = TradeOperation(tradeInfo, c, db)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, "Coin trade processed succesfully")
}
