package report

import (
	"database/sql"

	"github.com/heroku/go-getting-started/src/utils"
)

type TradeInfo struct {
	BuyingUserId  string
	SellingUserId string
	CoinsAmount   int
	OperationCost float64
	Date          string
}

type ReportInfo struct {
	id                 string
	user_buying_id     string
	user_selling_id    string
	coins_amount       int
	coin_unitary_value float64
	total_value        float64
	date               string
}

func (reportInfo *ReportInfo) ReportOperation(db *sql.DB) bool {
	_, err := db.Exec(`
		INSERT INTO transactions (user_buying_id, user_selling_id, coins_amount, coin_unitary_value, total_value, date)  
		VALUES ($1, $2, $3, $4, $5, $6);
	`, reportInfo.user_buying_id, reportInfo.user_selling_id,
		reportInfo.coins_amount, reportInfo.coin_unitary_value,
		reportInfo.total_value, reportInfo.date)

	if err == nil {
		return true
	} else {
		return false
	}
}

func (reportInfo *ReportInfo) Init(tradeInfo TradeInfo) {
	reportInfo.user_buying_id = tradeInfo.BuyingUserId
	reportInfo.user_selling_id = tradeInfo.SellingUserId
	reportInfo.coins_amount = tradeInfo.CoinsAmount
	reportInfo.coin_unitary_value = float64(utils.CoinPrice)
	reportInfo.total_value = float64(tradeInfo.CoinsAmount * utils.CoinPrice)
	reportInfo.date = tradeInfo.Date
}
