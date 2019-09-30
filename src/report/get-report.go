package report

import (
	"database/sql"
	"fmt"
	"log"

	"net/http"

	"github.com/gin-gonic/gin"
)

func getQuery(userId string, date string) string {
	selectStatement := "SELECT * FROM transactions "
	userCondition := fmt.Sprintf("(user_buying_id = '%v' or user_selling_id = '%v') ", userId, userId)
	dateCondition := fmt.Sprintf("date = '%v' ", date)

	var query string
	if userId != "" && date != "" {
		query = selectStatement + "WHERE " + userCondition + "AND " + dateCondition
	} else if userId != "" {
		query = selectStatement + "WHERE " + userCondition
	} else if date != "" {
		query = selectStatement + "WHERE " + dateCondition
	} else {
		query = selectStatement
	}

	return query
}

func GetReport(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := getQuery(c.Query("userId"), c.Query("date"))

		rows, err := db.Query(query)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error querying the transactions report on database")
			return
		}
		defer rows.Close()

		reportInfo := ReportInfo{}

		for rows.Next() {
			if err := rows.Scan(
				&reportInfo.id,
				&reportInfo.user_buying_id,
				&reportInfo.user_selling_id,
				&reportInfo.coins_amount,
				&reportInfo.coin_unitary_value,
				&reportInfo.total_value,
				&reportInfo.date,
			); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Report found \n id: %v \t coins_amount: %v \t date: %v \n",
				reportInfo.id, reportInfo.coins_amount, reportInfo.date)
		}

		c.String(http.StatusOK, "Report sucssesfull")
	}
}
