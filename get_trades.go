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
//		"state": 0,
//		"result": [
//		  {
//			"trade_id": "01J017Q6B3JGHZRP9D2NZHVKFX",
//			"price": "59498.63487492",
//			"base_volume": "94.00784310",
//			"quote_volume": "0.00158000",
//			"timestamp": 1718028573,
//			"type": "sell"
//		  },
//		  {
//			"trade_id": "01J017Q4CR4FEQ535GBWEBXZPR",
//			"price": "59480.45000999",
//			"base_volume": "11.68433959",
//			"quote_volume": "0.00019644",
//			"timestamp": 1718028571,
//			"type": "sell"
//		  }
//		]
//	  }
func GetTrades(currencyPair string) ([]Trade, error) {
	url := fmt.Sprintf(urlGetTrades, currencyPair)

	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error sending GET request: %w", err)
	}
	defer response.Body.Close()

	var responseStruct struct {
		State  int     `json:"state"`
		Result []Trade `json:"result"`
	}

	if err := json.NewDecoder(response.Body).Decode(&responseStruct); err != nil {
		return nil, fmt.Errorf("error decoding response payload: %w", err)
	}

	if response.StatusCode != http.StatusOK || responseStruct.State != 0 {
		return nil, fmt.Errorf("error fetching trades with status %s", response.Status)
	}

	return responseStruct.Result, nil
}
