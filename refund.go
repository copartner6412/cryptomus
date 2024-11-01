package cryptomus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// RefundPaymentRequest represents the parameters needed to request a refund.
//
// Invoice is identified by order_id or uuid, if you pass both, the account will be identified by uuid
//
// See "Refund" https://doc.cryptomus.com/business/payments/refund
//
// # Request example
//
//	{
//		"uuid": "8b03432e-385b-4670-8d06-064591096795",
//		"address": "TDD97yguPESTpcrJMqU6h2ozZbibv4Vaqm",
//		"is_subtract": "true",
//	}
type RefundRequest struct {
	RecordID
	// (Required) The address to which the refund should be made
	Address string `json:"address"`
	// (Required) Whether to take a commission from the merchant's balance or from the refund amount
	// true - take the commission from merchant balance
	// false - reduce the refundable amount by the commission amount
	IsSubtract bool `json:"is_subtract"`
}

// RefundPayment initiates a refund for a given invoice using either UUID or OrderID.
//
// See "Refund" https://doc.cryptomus.com/business/payments/refund
//
// # Response example
//
//	{
//	  "state": 0,
//	  "result": []
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
// If the invoice with the uuid or order_id you provided doesn’t exist:
//
//	{
//	    "state": 1,
//	    "message": "Payment was not found"
//	}
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
//
// If payouts are blocked on your merchant, you will receive this error message. This happens, for example, when you change your password and payouts are blocked for 24 hours.
//
//	{
//	    "state": 1,
//	    "message": "You are forbidden"
//	}
//
// If merchant wallet that is connected to an invoice does not exist:
//
//	{
//	    "state": 1,
//	    "message": "Undefined merchant wallet"
//	}
//
// If the invoice is not finalized or not paid.
//
//	{
//	    "state": 1,
//	    "message": "Refunds are made only for completed payments"
//	}
//
// If the amount to be refunded exceeds the amount of the invoice payment:
//
//	{
//	    "state": 1,
//	    "message": "The refund amount should not be more than the amount paid"
//	}
//
// If technical work occurs and the payment is temporarily unavailable, you can receive this error message:
//
//	{
//	    "state": 1,
//	    "message": "Server error"
//	}
func (m *Merchant) Refund(request RefundRequest) error {
	httpResponse, err := m.sendPaymentRequest("POST", urlRefund, request)
	if err != nil {
		return err
	}
	defer httpResponse.Body.Close()

	var response = struct {
		State   int    `json:"state"`
		Message string `json:"message"`
		Errors  struct {
			Address []string `json:"address"`
			UUID    []string `json:"uuid"`
			OrderID []string `json:"order_id"`
		} `json:"errors"`
		Code  int    `json:"code"`
		Error string `json:"error"`
	}{}

	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return fmt.Errorf("error decoding response: %w", err)
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
	errs = append(errs, response.Errors.Address...)

	if httpResponse.StatusCode != http.StatusOK || response.State != 0 || len(errs) > 0 {
		return fmt.Errorf("error processing refund with status %s: %v", httpResponse.Status, strings.Join(errs, "; "))
	}

	return nil
}
