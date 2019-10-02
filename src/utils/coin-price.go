package utils

import (
	"os"
	"time"

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
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest", nil)
	if err != nil {
		fmt.Println("Error requesting coin price", err)
		return 0, err
	}

	q := url.Values{}
	q.Add("start", "1")
	q.Add("limit", "5000")
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

	return 1, nil
}

func (coinInfo *CoinInfo) GetPrice() (float64, error) {
	nowTimestamp := time.Now().Unix()

	if nowTimestamp <= coinInfo.expirationTime {
		return coinInfo.price, nil
	} else {
		return coinPriceRequest()
	}
}
