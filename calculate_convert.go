package cryptomus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// See "Calculate convert" https://doc.cryptomus.com/personal/converts/calculate
//
// # Request example
//
//	{
//		"from": "BTC",
//		"to": "USDT",
//		"from_amount": "0.0001"
//	}
type Convert struct {
	// (Required) Currency From
	//    max: 10
	From string `json:"from"`
	// (Required) Amount to be calculated in Currency From
	//    max: 10
	To string `json:"to"`
	// (Required) Amount to be calculated in Currency From
	//    required_without: to_amount
	FromAmount string `json:"from_amount"`
	// (Required) Amount to be calculated in Currency To
	//    required_without: from_amount
	ToAmount string `json:"to_amount"`
}

// See "Calculate convert" https://doc.cryptomus.com/personal/converts/calculate
//
// # Response example
//
//	{
//		"state": 0,
//		"result": {
//		  "from": "0.001",
//		  "approximate_rate": "57853.000",
//		  "commission": "3",
//		  "total_amount": "60.000",
//		  "to": "60.000"
//		}
//	}
type CalculateConvertResponse struct {
	// Amount from
	From string `json:"from"`
	// Convert course
	Approximate_rate string `json:"approximate_rate"`
	// Commission
	Commission string `json:"commission"`
	// Total amount
	TotalAmount string `json:"total_amount"`
	// Amount to
	To string `json:"to"`
}

// See "Calculate convert" https://doc.cryptomus.com/personal/converts/calculate
//
// # Response example
//
//	{
//		"state": 0,
//		"result": {
//		  "from": "0.001",
//		  "approximate_rate": "57853.000",
//		  "commission": "3",
//		  "total_amount": "60.000",
//		  "to": "60.000"
//		}
//	}
func (u *User) CalculateConvert(request Convert) (*CalculateConvertResponse, error) {
	httpResponse, err := u.sendPaymentRequest("POST", urlCalculateConvert, struct{}{})
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	var response = struct {
		State   int                      `json:"state"`
		Result  CalculateConvertResponse `json:"result"`
		Message string                   `json:"message"`
		Errors  struct {
			From       []string `json:"from"`
			To         []string `json:"to"`
			FromAmount []string `json:"from_amount"`
			ToAmount   []string `json:"to_amount"`
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
	errs = append(errs, response.Errors.FromAmount...)
	errs = append(errs, response.Errors.ToAmount...)

	if httpResponse.StatusCode != http.StatusOK || response.State != 0 || len(errs) > 0 {
		return nil, fmt.Errorf("error calculating convert with status %s: %v", httpResponse.Status, strings.Join(errs, "; "))
	}

	return &response.Result, nil

}
