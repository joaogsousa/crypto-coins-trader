package transactions

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/src/report"
	"github.com/heroku/go-getting-started/src/utils"
)

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
		c.String(http.StatusBadRequest, "Bad request. Provide userId, coinsAmount and date as POST form values")
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

	tradeInfo := report.TradeInfo{
		BuyingUserId:  buyingUserId,
		SellingUserId: sellingUserId,
		CoinsAmount:   coinsAmount,
		OperationCost: float64(coinsAmount * utils.CoinPrice),
		Date:          c.PostForm("date"),
	}

	if tradeInfo.BuyingUserId == tradeInfo.SellingUserId {
		c.String(http.StatusBadRequest, "Invalid operation. buyingUserId and sellingUserId are the same.")
		return
	}

	ok, feedback := checkAvailableCredit(tradeInfo.BuyingUserId, tradeInfo.OperationCost, db)
	if !ok {
		c.String(http.StatusMethodNotAllowed, feedback)
		return
	}

	ok, feedback = checkAvailableCoins(tradeInfo.SellingUserId, tradeInfo.CoinsAmount, db)
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

	fmt.Println("Coin trade registered in database (without report yet)")

	//insert operation report on database
	var reportInfo *report.ReportInfo = &report.ReportInfo{}
	reportInfo.Init(tradeInfo)
	ok = reportInfo.ReportOperation(db)

	if ok {
		c.String(http.StatusOK, "Coin trade operation registered succesfully")
	} else {
		c.String(http.StatusInternalServerError, "Server error, unable to register coin trade")
	}
}
