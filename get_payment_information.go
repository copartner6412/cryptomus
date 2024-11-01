package cryptomus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// PaymentInformation retrieves payment information based on either UUID or Order ID.
//
// See "Payment information" https://doc.cryptomus.com/business/payments/payment-information
//
// # Possible errors
//
// code : 422
//
// example responses:
//
//	{
//		"state": 1,
//		"errors": {
//			"uuid": ["validation.required_without"],
//			"order_id": ["validation.required_without"]
//		}
//	}
func (m *Merchant) GetPaymentInformation(request RecordID) (*Payment, error) {
	httpResponse, err := m.sendPaymentRequest("POST", urlGetPaymentInformation, request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	var response = struct {
		State   int     `json:"state"`
		Result  Payment `json:"result"`
		Message string  `json:"message"`
		Errors  struct {
			UUID    []string `json:"uuid"`
			OrderID []string `json:"order_id"`
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
	errs = append(errs, response.Errors.UUID...)
	errs = append(errs, response.Errors.OrderID...)

	if httpResponse.StatusCode != http.StatusOK || response.State != 0 || len(errs) > 0 {
		return nil, fmt.Errorf("error getting payment information with status %s: %v", httpResponse.Status, strings.Join(errs, "; "))
	}

	return &response.Result, nil
}
