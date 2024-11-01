package cryptomus

import "time"

// See "Create market order" https://doc.cryptomus.com/personal/converts/market-order
//
// See "Create limit order" https://doc.cryptomus.com/personal/converts/limit-order
//
// # Response example
//
//	{
//		"state": 0,
//		"result": {
//		  "order_id": "2d9bf426-98ef-448b-84c2-03cc1ec78feb",
//		  "convert_amount_from": "10.000",
//		  "convert_amount_to": "3.000",
//		  "executed_amount_from": null,
//		  "executed_amount_to": null,
//		  "convert_currency_from": "USDT",
//		  "convert_currency_to": "XMR",
//		  "type": "market",
//		  "status": "completed",
//		  "created_at": "2024-07-11 , 18:06:04",
//		  "current_rate": "100",
//		  "completed_at": "2024-07-11 , 18:06:04"
//		}
//	}
type MarketOrder struct {
	// Id of convert
	OrderID string `json:"order_id"`
	// Convert amount from
	ConvertAmountFrom string `json:"convert_amount_from"`
	// Convert amount to
	ConvertAmountTo string `json:"convert_amount_to"`
	// Executed amount to
	ExecutedAmountTo string `json:"executed_amount_to"`
	// Executed amount from
	ExecutedAmountFrom string `json:"executed_amount_from"`
	// Convert currency from
	ConvertCurrencyFrom string `json:"convert_currency_from"`
	// Convert currency to
	ConvertCurrencyTo string `json:"convert_currency_to"`
	// Type enum
	//
	// Available options:
	//  - market
	//  - limit
	Type string `json:"type"`
	// Status enum
	//
	// Available options:
	//  - active
	//  - completed
	//  - partially_completed
	//  - cancelled
	//  - expired
	//  - failed
	Status string `json:"status"`
	// Date time create
	CreatedAt time.Time `json:"created_at"`
	// Current rate
	CurrentRate string `json:"current_rate"`
	// Limit value (only if type limit)
	Limit string `json:"limit"`
	// Limit expires date time
	ExpiresAt time.Time `json:"expires_at"`
	// Date time when order completed (only if order completed)
	CompletedAt time.Time `json:"completed_at"`
}
