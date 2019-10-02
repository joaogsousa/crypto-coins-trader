package report

import (
	"database/sql"
	"fmt"
	"log"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/src/jwtauth"
)

type TransactionInfo struct {
	Id           string
	Coins_amount int
	Date         string
	User_b_id    int
	User_b_email string
	User_s_id    int
	User_s_email string
}

func getQuery(userId string, date string) string {
	selectStatement := `
	SELECT 
	transactions.id, transactions.coins_amount, transactions.date,
	user_b.id AS userb_id, user_b.email AS userb_email, user_s.id AS users_id, user_s.email AS users_email
	FROM transactions 
	INNER JOIN users AS user_b on transactions.user_buying_id = user_b.id
	INNER JOIN users AS user_s on transactions.user_selling_id = user_s.id
	`
	userCondition := fmt.Sprintf(`(transactions.user_buying_id = %v OR transactions.user_selling_id = %v) `, userId, userId)
	dateCondition := fmt.Sprintf(`transactions.date = '%v' `, date)

	var query string
	if userId != "" && date != "" {
		query = selectStatement + ` WHERE ` + userCondition + ` AND ` + dateCondition
	} else if userId != "" {
		query = selectStatement + ` WHERE ` + userCondition
	} else if date != "" {
		query = selectStatement + ` WHERE ` + dateCondition
	} else {
		query = selectStatement
	}

	return query
}

func GetReport(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !jwtauth.IsAuthorized(c) {
			c.String(http.StatusUnauthorized, "Unautorized, please sign in first.")
			return
		}

		query := getQuery(c.Query("userId"), c.Query("date"))

		rows, err := db.Query(query)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error querying the transactions report on database")
			fmt.Println("Error querying the transactions report on database: ", err.Error())
			return
		}
		defer rows.Close()

		var transactionInfo TransactionInfo
		transactionRows := make([]TransactionInfo, 0)

		for rows.Next() {
			transactionInfo = TransactionInfo{}
			if err := rows.Scan(
				&transactionInfo.Id,
				&transactionInfo.Coins_amount,
				&transactionInfo.Date,
				&transactionInfo.User_b_id,
				&transactionInfo.User_b_email,
				&transactionInfo.User_s_id,
				&transactionInfo.User_s_email,
			); err != nil {
				log.Fatal(err)
			}
			transactionRows = append(transactionRows, transactionInfo)
		}

		if len(transactionRows) == 0 {
			c.String(http.StatusOK, "There is no transactions for the given query.")
			return
		}

		c.JSON(http.StatusOK, transactionRows)
	}
}
