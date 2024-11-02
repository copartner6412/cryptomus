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
// # Response example
//
//	{
//		"state": 0,
//		"result": {
//			"uuid": "70b8db5c-b952-406d-af26-4e1c34c27f15",
//			"order_id": "65bbe87b4098c17a31cff3e71e515243",
//			"amount": "15.00",
//			"payment_amount": "0.00",
//			"payer_amount": "15.75",
//			"discount_percent": -5,
//			"discount": "-0.75",
//			"payer_currency": "USDT",
//			"currency": "USDT",
//			"comments": null,
//			"merchant_amount": "15.43500000",
//			"network": "tron",
//			"address": "TXhfYSWt2oKRrHAJVJeYRuit6ZzKuoEKXj",
//			"from": null,
//			"txid": null,
//			"payment_status": "cancel",
//			"url": "https://pay.cryptomus.com/pay/70b8db5c-b952-406d-af26-4e1c34c27f15",
//			"expired_at": 1689099831,
//			"status": "cancel",
//			"is_final": true,
//			"additional_data": null,
//			"created_at": "2023-07-11T20:23:52+03:00",
//			"updated_at": "2023-07-11T21:24:17+03:00"
//		}
//	}
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
		return nil, fmt.Errorf("error with status %s: %v", httpResponse.Status, strings.Join(errs, "; "))
	}

	return &response.Result, nil
}
