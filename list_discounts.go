package cryptomus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// See "Set discount to payment method" https://doc.cryptomus.com/business/discount/set
//
// # Example
//
//	{
//	    "state": 0,
//	    "result": {
//	      "currency": "BUSD",
//	      "network": "bsc",
//	      "discount": "-10"
//	    }
//	}
type Discount struct {
	// Blockchain network code
	Network string `json:"network"`
	// Currency code
	Currency string `json:"currency"`
	// Discount percent
	Discount string `json:"Discount percent"`
}

// See "List of discounts" https://doc.cryptomus.com/business/discount/list
//
//	{
//		"state": 0,
//		"result": [
//		  {
//			"currency": "BTC",
//			"network": "btc",
//			"discount": 0
//		  },
//		  {
//			"currency": "BUSD",
//			"network": "bsc",
//			"discount": -10
//		  },
//		  {
//			"currency": "DASH",
//			"network": "dash",
//			"discount": 0
//		  }
//		]
//	}
func (m *Merchant) ListDiscounts() ([]Discount, error) {
	httpResponse, err := m.sendPaymentRequest("POST", urlListDiscounts, struct{}{})
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	var response = struct {
		State   int        `json:"state"`
		Result  []Discount `json:"result"`
		Message string     `json:"message"`
		Code    int        `json:"code"`
		Error   string     `json:"error"`
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

	if httpResponse.StatusCode != http.StatusOK || response.State != 0 || len(errs) > 0 {
		return nil, fmt.Errorf("error listing discounts with status %s: %v", httpResponse.Status, strings.Join(errs, "; "))
	}

	return response.Result, nil
}
