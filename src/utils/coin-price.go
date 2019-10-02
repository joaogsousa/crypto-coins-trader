package utils

import (
	"os"
	"time"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type CoinInfo struct {
	price          float64
	expirationTime int64
}

// exported object that must be used for pricing information
var CoinObj *CoinInfo = &CoinInfo{}

func coinPriceRequest() (float64, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET",
		"https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest", nil)
	if err != nil {
		fmt.Println("Error requesting coin price", err)
		return 0, err
	}

	q := url.Values{}
	q.Add("slug", "ethereum")
	q.Add("convert", "USD")

	apiKey := os.Getenv("COIN_MCAP_API_KEY")

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", apiKey)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to server", err)
		return 0, err
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Request for price done, see response!")
	fmt.Println(string(respBody))

	var result map[string]interface{}
	json.Unmarshal(respBody, &result)

	data := result["data"].(map[string]interface{})
	ethereum := data["ethereum"].(map[string]interface{})
	quote := ethereum["quote"].(map[string]interface{})
	usd := quote["USD"].(map[string]float64)
	var price float64 = usd["price"]

	return price, nil
}

func (coinInfo *CoinInfo) GetPrice() (float64, error) {
	nowTimestamp := time.Now().Unix()

	if nowTimestamp <= coinInfo.expirationTime {
		return coinInfo.price, nil
	} else {
		return coinPriceRequest()
	}
}
