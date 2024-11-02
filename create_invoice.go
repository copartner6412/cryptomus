package cryptomus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// CreateInvoice is a payment method that creates an invoice for merchant by sending a POST request to Cryptomus
//
// This payment method is suitable for paying for goods with a specific amount.
//
// You can specify the currency, then the payer can select the cryptocurrency and network on the payment page.You can also specify the crypto-currency and network immediately, then a crypto-currency address will be generated and you can give it to the client.
//
// Once the client has paid, you will receive a notification - a webhook, which indicates the status, the payment amount and hash, as well as additional data (more details - Link )
//
//   - If the client didn’t pay, you will receive the cancel status.
//   - If the client sent the exact amount, you will receive a paid status.
//   - If the client sent a lower amount, you will receive the wrong_amount status.
//   - If the client sent a higher amount, you will receive the paid_over status.
//
// To allow the client to pay the rest, pass the parameter is_payment_multiple = true
//
// On the payment page, the client will be notified that there is a surcharge to be paid. If the client doesn't pay the rest of an amount, you will receive a webhook with the wrong_amount status, when the invoice expires.
//
// See "Creating an invoice" https://doc.cryptomus.com/business/payments/creating-invoice
//
// # Payment example
//
//	{
//		"state": 0,
//		"result": {
//			"uuid": "26109ba0-b05b-4ee0-93d1-fd62c822ce95",
//			"order_id": "1",
//			"amount": "15.00",
//			"payment_amount": null,
//			"payer_amount": null,
//			"discount_percent": null,
//			"discount": "0.00000000",
//			"payer_currency": null,
//			"currency": "USD",
//			"merchant_amount": null,
//			"network": null,
//			"address": null,
//			"from": null,
//			"txid": null,
//			"payment_status": "check",
//			"url": "https://pay.cryptomus.com/pay/26109ba0-b05b-4ee0-93d1-fd62c822ce95",
//			"expired_at": 1689098133,
//			"status": "check",
//			"is_final": false,
//			"additional_data": null,
//			"created_at": "2023-07-11T20:23:52+03:00",
//			"updated_at": "2023-07-11T21:24:17+03:00"
//		}
//	}
//
// # Possible errors
//
// # Validation errors
//
// code : 422
//
// If some parameter is required and not passed:
//
//	{
//		"state": 1,
//		"errors": {
//			"amount": ["validation.required"]
//		}
//	}
//
// If you provided the network code that is not supported
//
//	{
//	    "state": 1,
//	    "message": "The network was not found"
//	}
//
// If we don’t support the currency code you provided in currency parameter:
//
//	{
//	    "state": 1,
//	    "message": "The currency was not found"
//	}
//
// If your payments are blocked. In this case you need to contact support.
//
//	{
//	    "state": 1,
//	    "message": "You are forbidden"
//	}
//
// If we don’t find the payment service of the currency that was passed in to_currency parameter:
//
//	{
//	    "state": 1,
//	    "message": "Not found service to_currency"
//	}
//
// If accidentally something happens when converting to to_currency:
//
//	{
//	    "state": 1,
//	    "message": "Error convert to_currency"
//	}
//
// If the payment amount is less than the minimum supported amount for cryptocurrency, you will receive this error message:
//
//	{
//	    "state": 1,
//	    "message": "Minimum amount 0.5 USDT"
//	}
//
// If the payment amount is greater than the maximum supported amount for cryptocurrency, you will receive this error message:
//
//	{
//	    "state": 1,
//	    "message": "Maximum amount 10000000 USDT"
//	}
//
// If you will not have an active merchant wallet with a cryptocurrency of the payment, you will receive this error:
//
//	{
//	    "state": 1,
//	    "message": "Wallet not found"
//	}
//
// If technical work occurs and the payment is temporarily unavailable, you can receive this error messages:
//
//	{
//	    "state": 1,
//	    "message": "Gateway error"
//	}
//
//	{
//	    "state": 1,
//	    "message": "The terminal was not found"
//	}
//
//	{
//	    "state": 1,
//	    "message": "Server error"
//	}
//
// # Internal server error
//
// code : 500
//
//	{
//		"message": "Server error, #1",
//		"code": 500,
//		"error": null
//	}
func (m *Merchant) CreateInvoice(request Invoice) (*Payment, error) {
	httpResponse, err := m.sendPaymentRequest("POST", urlCreateInvoice, request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	var response = struct {
		State   int     `json:"state"`
		Result  Payment `json:"result"`
		Message string  `json:"message"`
		// If some parameter is required and not passed
		Errors struct {
			Amount   []string `json:"amount"`
			Currency []string `json:"currency"`
			OrderID  []string `json:"order_id"`
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
	errs = append(errs, response.Errors.Amount...)
	errs = append(errs, response.Errors.Currency...)
	errs = append(errs, response.Errors.OrderID...)

	if httpResponse.StatusCode != http.StatusOK || response.State != 0 || len(errs) > 0 {
		return nil, fmt.Errorf("error with status %s: %v", httpResponse.Status, strings.Join(errs, "; "))
	}

	return &response.Result, nil
}
