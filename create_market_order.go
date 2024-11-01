package cryptomus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// See "Create market order" https://doc.cryptomus.com/personal/converts/market-order
//
// # Request example
//
//	{
//		"from": "BTC",
//		"to": "USDT",
//		"amount": "0.0001"
//	}
type MarketOrderRequest struct {
	// (Required) Currency to convert From
	//    max: 10
	From string `json:"from"`
	// (Required) Currency to convert To
	//    max: 10
	To string `json:"to"`
	// (Required) Amount to be converted in currency from
	//
	// If there are pennies in the amount, then send them with a separator '.'
	//    Example: 10.28
	Amount string `json:"amount"`
}

// See "Create market order" https://doc.cryptomus.com/personal/converts/market-order
//
// # Example response
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
func (u *User) CreateMarketOrder(request MarketOrderRequest) (*MarketOrder, error) {
	httpResponse, err := u.sendPaymentRequest("POST", urlCreateMarketOrder, struct{}{})
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	var response = struct {
		State   int         `json:"state"`
		Result  MarketOrder `json:"result"`
		Message string      `json:"message"`
		Errors  struct {
			From   []string `json:"from"`
			To     []string `json:"to"`
			Amount []string `json:"amount"`
		} `json:"errors"`
		Code  int    `json:"code"`
		Error string `json:"error"`
	}{}

	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	var errs []string
	if response.Message != "" {
		errs = append(errs, response.Message)
	}
	if response.Error != "" {
		errs = append(errs, response.Error)
	}
	errs = append(errs, response.Errors.From...)
	errs = append(errs, response.Errors.To...)
	errs = append(errs, response.Errors.Amount...)

	if httpResponse.StatusCode != http.StatusOK || response.State != 0 || len(errs) > 0 {
		return nil, fmt.Errorf("error creating market order with status %s: %v", httpResponse.Status, strings.Join(errs, "; "))
	}

	return &response.Result, nil
}
