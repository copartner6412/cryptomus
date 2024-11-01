package cryptomus

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ExchangeRate struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Course string `json:"course"`
}

// See "List" https://doc.cryptomus.com/business/exchange-rates/list
//
// # Response example
//
//	{
//		"state": 0,
//		"result": [
//		  {
//			"from": "ETH",
//			"to": "USD",
//			"course": "1228.45000000"
//		  }
//		]
//	  }
func GetExchangeRate(currency string) ([]ExchangeRate, error) {
	url := fmt.Sprintf(urlGetExchangeRate, currency)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var exchangeRateResponse struct {
		State  int            `json:"state"`
		Result []ExchangeRate `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&exchangeRateResponse); err != nil {
		return nil, err
	}

	return exchangeRateResponse.Result, nil
}
