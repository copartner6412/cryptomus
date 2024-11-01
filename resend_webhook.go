package cryptomus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// ResendWebhook resends the webhook for a finalized invoice identified by either UUID or OrderID.
//
// Resend the webhook by invoice. You can resend the webhook only for finalized invoices, that is, invoices in statuses: wrong_amount, paid, paid_over.
//
// To resend the webhook on the invoice, the url_callback must be specified at the time of invoice creation.
//
// See "Resend webhook" https://doc.cryptomus.com/business/payments/resend-webhook
//
// # Response example
//
//	{
//		"state": 0,
//	  "result": []
//	}
//
// # Possible errors
//
// # Validation Errors
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
//
// If the invoice is not found
//
//	{
//		"state": 1,
//		"message": "Payment not found"
//	}
//
// If the url_callback was not specified when creating an invoice, and no notification was send, there is nothing to resend and you will receive this error message:
//
//	{
//		"state": 1,
//		"message": "Notification not found"
//	}
//
// The maximum number of resendings is 10, if you exceed this limit, you will receive an error message. If you want to resend the notification over the limit, you can contact the support.
//
//	{
//		"state": 1,
//		"message": "Too much resend"
//	}
func (m *Merchant) ResendWebhook(request RecordID) error {
	httpResponse, err := m.sendPaymentRequest("POST", urlResendWebhook, request)
	if err != nil {
		return err
	}
	defer httpResponse.Body.Close()

	var response = struct {
		State   int    `json:"state"`
		Message string `json:"message"`
		Errors  struct {
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

	if httpResponse.StatusCode != http.StatusOK || response.State != 0 || len(errs) > 0 {
		return fmt.Errorf("error resending webhook with status %s: %v", httpResponse.Status, strings.Join(errs, "; "))
	}

	return nil
}
