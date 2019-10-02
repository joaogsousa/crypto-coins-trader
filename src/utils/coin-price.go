package utils

import (
	"log"
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

	var result map[string]interface{}
	json.Unmarshal(respBody, &result)

	data := result["data"].(map[string]interface{})
	ethereum := data["1027"].(map[string]interface{})
	quote := ethereum["quote"].(map[string]interface{})
	usd := quote["USD"].(map[string]interface{})

	return usd["price"].(float64), nil
}

func (coinInfo *CoinInfo) GetPrice() float64 {
	nowTimestamp := time.Now().Unix()

	if nowTimestamp <= coinInfo.expirationTime {
		fmt.Println("Used previously set coin price, i.e, price not expired yet")
		return coinInfo.price
	} else {
		fmt.Println("fetch new coin price, i.e, price expired or not set yet")
		coinPrice, err := coinPriceRequest()
		if err != nil {
			log.Fatal("Error getting the coin price")
		}

		coinInfo.price = coinPrice
		coinInfo.expirationTime = time.Now().Add(time.Hour).Unix()

		return coinPrice
	}
}
