package cryptomus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// See "Get order book" https://doc.cryptomus.com/personal/market-cap/orderbook
type Order struct {
	Price    string `json:"price"`
	Quantity string `json:"quantity"`
}

// Available options for level of volume: 0, 1, 2, 3, 4, 5
//
// See "Get order book" https://doc.cryptomus.com/personal/market-cap/orderbook
//
// # Response example
//
//	{
//		"state": 0,
//		"result": {
//		  "timestamp": "1724069797.1308",
//		  "bids": [
//			{
//			  "price": "0.04548320",
//			  "quantity": "12462000"
//			},
//			{
//			  "price": "3.00000000",
//			  "quantity": "12457000"
//			}
//		  ],
//		  "asks": [
//			{
//			  "price": "2.73042000",
//			  "quantity": "12506000"
//			},
//			{
//			  "price": "0.33660000",
//			  "quantity": "12508000"
//			}
//		  ]
//		}
//	  }
func GetOrderBook(currencyPair string, level int) (timestamp time.Time, bids, asks []Order, err error) {
	url := fmt.Sprintf(urlGetOrderBook+"?level=%d", currencyPair, level)

	response, err := http.Get(url)
	if err != nil {
		return time.Time{}, nil, nil, fmt.Errorf("error sending GET request: %w", err)
	}
	defer response.Body.Close()

	var responseStruct struct {
		State  int `json:"state"`
		Result struct {
			Timestamp string  `json:"timestamp"`
			Bids      []Order `json:"bids"`
			Asks      []Order `json:"asks"`
		} `json:"result"`
	}
	if err := json.NewDecoder(response.Body).Decode(&responseStruct); err != nil {
		return time.Time{}, nil, nil, fmt.Errorf("error decoding response payload: %w", err)
	}

	if response.StatusCode != http.StatusOK || responseStruct.State != 0 {
		return time.Time{}, nil, nil, fmt.Errorf("error fetching order book with status %s", response.Status)
	}

	timestamp, err = parseUnixTimeString(responseStruct.Result.Timestamp)
	if err != nil {
		return time.Time{}, nil, nil, fmt.Errorf("error converting timestamp: %w", err)
	}

	return timestamp, responseStruct.Result.Bids, responseStruct.Result.Asks, nil
}

func parseUnixTimeString(unixDecimal string) (time.Time, error) {
	parts := strings.Split(unixDecimal, ".")

	seconds, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing seconds: %w", err)
	}

	var nanoseconds int64
	if len(parts) == 2 {
		fractionalPart := parts[1] + strings.Repeat("0", 9-len(parts[1]))
		nanoseconds, err = strconv.ParseInt(fractionalPart, 10, 64)
		if err != nil {
			return time.Time{}, fmt.Errorf("error parsing nanoseconds: %w", err)
		}
	}

	return time.Unix(seconds, nanoseconds), nil
}
