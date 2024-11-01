package cryptomus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// See "Set discount to payment method" https://doc.cryptomus.com/business/discount/set
//
// # Request example
//
//	{
//	    "network": "bsc"
//	    "currency": "BUSD",
//	    "discount_percent": -20,
//	}
type DiscountRequest struct {
	// (Required) Currency code
	Network string `json:"network"`
	// (Required) Blockchain network code
	Currency string `json:"currency"`
	// (Required) Discount percent
	// About discount percent:
	//
	// Positive Numbers (>0). Gives buyers a discount for paying with a coin. Good promotional tool if you want to give extra support to a particular coin.
	//
	// Negative Numbers (<0). Adds a certain percentage (padding) for paying with a coin. This could be used to cover your crypto/fiat conversion costs, make adjustments to match your local exchange, etc.
	DiscountPercent int `json:"discount_percent"`
}

// See "Set discount to payment method" https://doc.cryptomus.com/business/discount/set
//
// # Possible errors
//
//	{
//		"message": "The service BUSD bsc was not found",
//		"state": 1
//	}
func (m *Merchant) SetDiscount(request DiscountRequest) (*Discount, error) {
	httpResponse, err := m.sendPaymentRequest("POST", urlSetDiscount, request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	var response = struct {
		State   int      `json:"state"`
		Result  Discount `json:"result"`
		Message string   `json:"message"`
		Errors  struct {
			Currency []string `json:"currency"`
			Network  []string `json:"network"`
			Discount []string `json:"discount"`
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
	errs = append(errs, response.Errors.Currency...)
	errs = append(errs, response.Errors.Network...)
	errs = append(errs, response.Errors.Discount...)

	if httpResponse.StatusCode != http.StatusOK || response.State != 0 || len(errs) > 0 {
		return nil, fmt.Errorf("error setting discount with status %s: %v", httpResponse.Status, strings.Join(errs, "; "))
	}

	return &response.Result, nil

}
