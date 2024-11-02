package cryptomus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// See "Cancel recurring payment" https://doc.cryptomus.com/business/recurring/cancel
//
// # Response example
//
//	{
//		"state": 0,
//		"result": {
//		  "uuid": "bbe5ce96-1126-4843-a0d2-b432e77669c2",
//		  "name": "Access to personal account",
//		  "order_id": "1487555",
//		  "amount": "5",
//		  "currency": "USD",
//		  "payer_currency": "USDT",
//		  "payer_amount_usd": "5.00",
//		  "payer_amount": "5.00",
//		  "url_callback": null,
//		  "discount_days": "30",
//		  "discount_amount": "50.00",
//		  "end_of_discount": null,
//		  "period": "weekly",
//		  "status": "wait_accept",
//		  "url": "https://pay.cryptomus.com/pay/bbe5ce96-1126-4843-a0d2-b432e77669c2",
//		  "last_pay_off": null
//		}
//	}
func (m *Merchant) CancelRecurringPayment(request RecordID) (*RecurringPayment, error) {
	httpResponse, err := m.sendPaymentRequest("POST", urlCancelRecurringPayment, request)
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
		return nil, fmt.Errorf("error with status %s: %v", httpResponse.Status, strings.Join(errs, "; "))
	}

	return &response.Result, nil
}
