package cryptomus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// To get the recurring payment status you need to pass one of the required parameters, if you pass both, the account will be identified by order_id
//
// See "Payment information" https://doc.cryptomus.com/business/recurring/info
//
// # Possible errors
//
// # Validation errors
//
// code : 422
//
// example responses:
//
// If you will not pass a uuid and order_id, you will receive this message. It means that at least one of the parameters is required.
//
//	{
//		"state": 1,
//		"errors": {
//			"uuid": ["validation.required_without"],
//			"order_id": ["validation.required_without"]
//		}
//	}
func (m *Merchant) GetRecurringPaymentInformation(request RecordID) (*RecurringPayment, error) {
	httpResponse, err := m.sendPaymentRequest("POST", urlGetRecurringPaymentInformation, request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	var response = struct {
		State   int              `json:"state"`
		Result  RecurringPayment `json:"result"`
		Message string           `json:"message"`
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
		return nil, fmt.Errorf("error getting recurring payment information with status %s: %v", httpResponse.Status, strings.Join(errs, "; "))
	}

	return &response.Result, nil
}
