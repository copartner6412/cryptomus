package cryptomus

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// See "Get trades" https://doc.cryptomus.com/personal/market-cap/trades
type Trade struct {
	// Trade id
	TradeID string `json:"trade_id"`
	// Trade price
	Price string `json:"price"`
	// Trade base volume
	BaseVolume string `json:"base_volume"`
	// Trade quote volume
	QuoteVolume string `json:"quote_volume"`
	// Time
	Timestamp int `json:"timestamp"`
	// Direction type
	//
	// Available options:
	//  - sell
	//  - buy
	Type string `json:"type"`
}

// See "Get trades" https://doc.cryptomus.com/personal/market-cap/trades
//
// # Example response
//
//	{
//	  "data": [
//	    {
//	      "trade_id": "01JBP2KQ3VMKX8JSV3R0DCJ71Q",
//	      "price": "68928.9500000000000000",
//	      "base_volume": "64.9999998500000000",
//	      "quote_volume": "0.0009430000000000",
//	      "timestamp": 1730539019,
//	      "type": "sell"
//	    },
//	    {
//	      "trade_id": "01JBP00NCHWPH51YSS6TRK26H7",
//	      "price": null,
//	      "base_volume": null,
//	      "quote_volume": "0.2909080000000000",
//	      "timestamp": 1730536297,
//	      "type": "sell"
//	    },
//		...
//	  ]
//	}
func GetTrades(currencyPair string) ([]Trade, error) {
	url := fmt.Sprintf(urlGetTrades, currencyPair)

	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error sending GET request: %w", err)
	}
	defer response.Body.Close()

	var responseStruct struct {
		Data    []Trade `json:"data"`
		Code    int     `json:"code"`
		Message string  `json:"message"`
	}

	if err := json.NewDecoder(response.Body).Decode(&responseStruct); err != nil {
		return nil, fmt.Errorf("error decoding response payload: %w", err)
	}

	if response.StatusCode != http.StatusOK || responseStruct.Message != "" {
		return nil, fmt.Errorf("error with status %s: %s", response.Status, responseStruct.Message)
	}

	return responseStruct.Data, nil
}
