package transactions

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/src/report"
)

func CashOperation(userId string, operationValue float64, db *sql.DB) (bool, string) {
	row := db.QueryRow(`
		SELECT cash FROM users  
		WHERE id = $1;
	`, userId)

	var cash float64
	if err := row.Scan(&cash); err != nil {
		return false, "There is no user with the specifyed id, cash operation refused"
	}

	cash += operationValue

	_, err := db.Exec(`
		UPDATE users
		SET cash = $1  
		WHERE id = $2;
	`, fmt.Sprintf("%f", cash), userId)
	if err != nil {
		return false, "Unable to realize the cash operation"
	}
	return true, "Payment succesfull"
}

func CoinsOperation(userId string, coinsAmount int, db *sql.DB) (bool, string) {
	row := db.QueryRow(`
		SELECT coins FROM users  
		WHERE id = $1;
	`, userId)

	var coins int
	if err := row.Scan(&coins); err != nil {
		return false, "There is no user with the specifyed id, coins operation refused"
	}

	coins += coinsAmount

	_, err := db.Exec(`
		UPDATE users
		SET coins = $1  
		WHERE id = $2;
	`, strconv.Itoa(coins), userId)
	if err != nil {
		return false, "Unable to realize the coins trading operation"
	}
	return true, "Coins transfer succesfull"
}

func TradeOperation(tradeInfo report.TradeInfo, c *gin.Context, db *sql.DB) error {
	var ok bool
	var feedback string

	// charge money from the buying user
	ok, feedback = CashOperation(tradeInfo.BuyingUserId, -1*tradeInfo.OperationCost, db)
	if !ok {
		return errors.New(feedback)
	}

	// credit money from the selling user
	ok, feedback = CashOperation(tradeInfo.SellingUserId, tradeInfo.OperationCost, db)
	if !ok {
		return errors.New(feedback)
	}

	// charge coins from the selling user
	ok, feedback = CoinsOperation(tradeInfo.SellingUserId, -1*tradeInfo.CoinsAmount, db)
	if !ok {
		return errors.New(feedback)
	}

	// credit coins to the buying user
	ok, feedback = CoinsOperation(tradeInfo.BuyingUserId, tradeInfo.CoinsAmount, db)
	if !ok {
		return errors.New(feedback)
	}

	return nil
}
