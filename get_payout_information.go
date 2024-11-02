package cryptomus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// See "Payout information" https://doc.cryptomus.com/business/payouts/payout-information
//
// # Response example
//
//	{
//	    "state": 0,
//	    "result": {
//	        "uuid": "a7c0caec-a594-4aaa-b1c4-77d511857594",
//	        "amount": "3",
//	        "currency": "USDT",
//	        "network": "TRON",
//	        "address": "TJ...",
//	        "txid": null,
//	        "status": "process",
//	        "is_final": false,
//	        "balance": 129,
//	        "payer_currency": "USD",
//	        "payer_amount": 3
//	    }
//	}
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
func (m *Merchant) GetPayoutInformation(request RecordID) (*Payment, error) {
	httpResponse, err := m.sendPayoutRequest("POST", urlGetPayoutInformation, request)
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
		return nil, fmt.Errorf("error with status %s: %v", httpResponse.Status, strings.Join(errs, "; "))
	}

	return &response.Result, nil
}
