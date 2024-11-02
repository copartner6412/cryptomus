package cryptomus

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// See "Get assets" https://doc.cryptomus.com/personal/market-cap/assets
type Asset struct {
	// Code of currency
	CurrencyCode string `json:"currency_code"`
	// Code of crypto network
	NetworkCode string `json:"network_code"`
	// Is possible withdraw
	CanWithdraw bool `json:"can_withdraw"`
	// Is possible payment
	CanDeposit bool `json:"can_deposit"`
	// Minimum withdraw value
	MinWithdraw string `json:"min_withdraw"`
	// Maximum withdraw value
	MaxWithdraw string `json:"max_withdraw"`
	// Maximum deposit value
	MaxDeposit string `json:"max_deposit"`
	// Minimum deposit value
	MinDeposit string `json:"min_deposit"`
}

// See "Get assets" https://doc.cryptomus.com/personal/market-cap/assets
//
//	{
//		"state": 0,
//		  "result": [
//			{
//			  "currency_code": "CRMS",
//			  "network_code": "polygon",
//			  "can_withdraw": true,
//			  "can_deposit": false,
//			  "min_withdraw": "1.00000000",
//			  "max_withdraw": "10000000.00000000",
//			  "max_deposit": null,
//			  "min_deposit": null
//			},
//			{
//			  "currency_code": "DASH",
//			  "network_code": "dash",
//			  "can_withdraw": true,
//			  "can_deposit": true,
//			  "min_withdraw": "0.01000000",
//			  "max_withdraw": "1000000.00000000",
//			  "max_deposit": "1000000.00000000",
//			  "min_deposit": "0.02000000"
//			}
//		  ]
//	}
func GetAssets() ([]Asset, error) {
	response, err := http.Get(urlGetAssets)
	if err != nil {
		return nil, fmt.Errorf("error sending GET request: %w", err)
	}
	defer response.Body.Close()

	var responseStruct struct {
		State  int     `json:"state"`
		Result []Asset `json:"result"`
	}

	if err := json.NewDecoder(response.Body).Decode(&responseStruct); err != nil {
		return nil, fmt.Errorf("error decoding response payload: %w", err)
	}

	if response.StatusCode != http.StatusOK || responseStruct.State != 0 {
		return nil, fmt.Errorf("error fetching assets with status %s", response.Status)
	}

	return responseStruct.Result, nil
}
