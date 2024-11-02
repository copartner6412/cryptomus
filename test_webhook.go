package cryptomus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// You may to pass one of the uuid or order_id parameters, if you pass both, the account will be identified by uuid
//
// See "Testing webhook" https://doc.cryptomus.com/business/payments/testing-webhook
//
// # Request example
//
//	{
//		"uuid": "e1830f1b-50fc-432e-80ec-15b58ccac867",
//		"currency": "ETH",
//		"url_callback": "https://your.site/callback",
//		"network": "eth",
//		"status": "paid"
//	}
type TestWebhookRequest struct {
	RecordID
	// (Required) Invoice currency code
	Currency string `json:"currency"`
	// (Required) Url to which webhooks with payment status will be sent
	//    min: 6
	//    max: 150
	//    url
	URLCallback string `json:"url_callback"`
	// (Required) Invoice network code
	Network string `json:"network"`
	// (Required) Payment status
	//    default: paid
	// Available options for payment/wallet:
	//  - process
	//  - check
	//  - paid
	//  - paid_over
	//  - fail
	//  - wrong_amount
	//  - cancel
	//  - system_fail
	//  - refund_process
	//  - refund_fail
	//  - refund_paid
	// Available options for payout:
	//  - process
	//  - check
	//  - paid
	//  - fail
	//  - cancel
	//  - system_fail
	Status string `json:"status"`
}

// To validate the signature from the webhook data array, use the payment API key.
//
// To ensure that you are correctly receiving webhooks and can validate the signature, you should use this method to test webhooks for payment.
//
// Please note that no data is saved to the database, and any data received in the webhook is only stored in an array for testing purposes to ensure the correctness of the signature and to test the retrieval of this array from us.
//
// To test a webhook with an existing invoice, please provide its uuid or order ID. If these parameters are not provided, the webhook will be sent with a test invoice.
//
// See "Testing webhook" https://doc.cryptomus.com/business/payments/testing-webhook
//
// # Response example
//
//	{
//		"state": 0,
//		"result": []
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
// If you pass uuid or order_id and invoice/static wallet/payout is not found:
//
//	{
//	    "state": 1,
//	    "message": "Not found payment"
//	}
//
// If you don’t pass the required parameters:
//
//	{
//		"state": 1,
//		"errors": {
//			"currency": ["validation.required"]
//	}
//
// }
// If the provided currency and network are not crypto or not supported:
//
//	{
//	    "state": 1,
//	    "message": "Payment service not found"
//	}
func (m *Merchant) TestWebhookPayment(request TestWebhookRequest) error {
	httpResponse, err := m.sendPaymentRequest("POST", urlTestWebhookPayment, request)
	if err != nil {
		return err
	}
	defer httpResponse.Body.Close()

	var response = struct {
		State   int    `json:"state"`
		Message string `json:"message"`
		Errors  struct {
			UUID        []string `json:"uuid"`
			OrderID     []string `json:"order_id"`
			Currency    []string `json:"currency"`
			URLCallback []string `json:"url_callback"`
			Network     []string `json:"network"`
			Status      []string `json:"status"`
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
	errs = append(errs, response.Errors.Currency...)
	errs = append(errs, response.Errors.URLCallback...)
	errs = append(errs, response.Errors.Network...)
	errs = append(errs, response.Errors.Status...)

	if httpResponse.StatusCode != http.StatusOK || response.State != 0 || len(errs) > 0 {
		return fmt.Errorf("error with status %s: %v", httpResponse.Status, strings.Join(errs, "; "))
	}

	return nil
}

// To validate the signature from the webhook data array, use the payment API key.
//
// See "Testing webhook" https://doc.cryptomus.com/business/payments/testing-webhook
//
// # Response example
//
//	{
//			"state": 0,
//			"result": []
//	}
//
// # Possible errors
//
// # Validation errors
//
// code : 422
//
// example errors:
//
// If you pass uuid or order_id and invoice/static wallet/payout is not found:
//
//	{
//	    "state": 1,
//	    "message": "Not found wallet"
//	}
//
// If you don’t pass the required parameters:
//
//	{
//		"state": 1,
//		"errors": {
//			"currency": ["validation.required"]
//		}
//	}
func (m *Merchant) TestWebhookWallet(request TestWebhookRequest) error {
	httpResponse, err := m.sendPaymentRequest("POST", urlTestWebhookWallet, request)
	if err != nil {
		return err
	}
	defer httpResponse.Body.Close()

	var response = struct {
		State   int    `json:"state"`
		Message string `json:"message"`
		Errors  struct {
			UUID        []string `json:"uuid"`
			OrderID     []string `json:"order_id"`
			Currency    []string `json:"currency"`
			URLCallback []string `json:"url_callback"`
			Network     []string `json:"network"`
			Status      []string `json:"status"`
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
	errs = append(errs, response.Errors.Currency...)
	errs = append(errs, response.Errors.URLCallback...)
	errs = append(errs, response.Errors.Network...)
	errs = append(errs, response.Errors.Status...)

	if httpResponse.StatusCode != http.StatusOK || response.State != 0 || len(errs) > 0 {
		return fmt.Errorf("error with status %s: %v", httpResponse.Status, strings.Join(errs, "; "))
	}

	return nil
}

// To validate the signature from the webhook data array, use the payment API key.
//
// See "Testing webhook" https://doc.cryptomus.com/business/payments/testing-webhook
//
// # Response example
//
//	{
//		"state": 0,
//		"result": []
//	  }
//
// # Possible errors
//
// # Validation errors
//
// code : 422
//
// example responses:
//
// If you pass uuid or order_id and invoice/static wallet/payout is not found:
//
//	{
//	    "state": 1,
//	    "message": "Not found payout"
//	}
//
// If you don’t pass the required parameters:
//
//	{
//		"state": 1,
//		"errors": {
//			"currency": ["validation.required"]
//		}
//	}
//
// If the provided currency and network are not crypto or not supported:
//
//	{
//	    "state": 1,
//	    "message": "Payout service not found"
//	}
func (m *Merchant) TestWebhookPayout(request TestWebhookRequest) error {
	httpResponse, err := m.sendPayoutRequest("POST", urlTestWebhookPayout, request)
	if err != nil {
		return err
	}
	defer httpResponse.Body.Close()

	var response = struct {
		State   int    `json:"state"`
		Message string `json:"message"`
		Errors  struct {
			UUID        []string `json:"uuid"`
			OrderID     []string `json:"order_id"`
			Currency    []string `json:"currency"`
			URLCallback []string `json:"url_callback"`
			Network     []string `json:"network"`
			Status      []string `json:"status"`
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
	errs = append(errs, response.Errors.Currency...)
	errs = append(errs, response.Errors.URLCallback...)
	errs = append(errs, response.Errors.Network...)
	errs = append(errs, response.Errors.Status...)

	if httpResponse.StatusCode != http.StatusOK || response.State != 0 || len(errs) > 0 {
		return fmt.Errorf("error with status %s: %v", httpResponse.Status, strings.Join(errs, "; "))
	}

	return nil
}
