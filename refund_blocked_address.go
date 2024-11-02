package cryptomus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// RefundBlockedAddressRequest represents the parameters needed to refund payments on a blocked wallet address.
//
// To refund payments you need to pass either uuid or order_id, if you pass both, the static wallet will be identified by uuid
//
// See "Refund payments on blocked address" https://doc.cryptomus.com/business/payments/refundblocked
//
// # Request example
//
//	{
//	    "order_id": "1",
//	    "address": "TK8..."
//	}
type RefundBlockedAddressRequest struct {
	RecordID
	// (Required) Refund all blocked funds to this address
	//    min: 10
	//    max: 128
	//    alpha_dash
	Address string `json:"address"`
}

// RefundBlockedAddressResponse represents the result of a successful refund request.
//
// See "Refund payments on blocked address" https://doc.cryptomus.com/business/payments/refundblocked
//
// # Response example
//
//	{
//	    "state": 0,
//	    "result": {
//	        "commission": "0.50",
//	        "amount": "9.50"
//	    }
//	}
type RefundBlockedAddressResponse struct {
	// Commission of refund
	Commission string `json:"commission"`
	// Amount of refund
	Amount string `json:"amount"`
}

// RefundBlockedAddress refunds blocked funds to a specified address, identified by either a UUID or Order ID.
//
// See "Refund payments on blocked address" https://doc.cryptomus.com/business/payments/refundblocked
//
// # Response example
//
//	{
//	    "state": 0,
//	    "result": {
//	        "commission": "0.50",
//	        "amount": "9.50"
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
//
// If static wallet was not found:
//
//	{
//	    "state": 1,
//	    "message": "Not found"
//	}
//
// Al long as withdrawal is allowed only once from blocked static wallet, If the withdrawal was already made you will receive this error message:
//
//	{
//	    "state": 1,
//	    "message": "You can withdraw only once"
//	}
//
// Al long as withdrawal is allowed only once from blocked static wallet, If you blocked a static wallet with a parameter is_force_refund=true, the withdrawal was already made and you will receive this error message:
//
//	{
//	    "state": 1,
//	    "message": "You can withdraw only once"
//	}
//
// If refund is already in process:
//
//	{
//	    "state": 1,
//	    "message": "Refund is in process"
//	}
//
// If the payout is blocked on your merchant or user, you will receive this error message. This happens, for example, when you change your password and payouts are blocked for 24 hours.
//
//	{
//	    "state": 1,
//	    "message": "You are forbidden"
//	}
//
// If no payments were made to the static wallet:
//
//	{
//	    "state": 1,
//	    "message": "Nothing to withdraw"
//	}
//
// If the payout amount is less than the minimum supported amount for cryptocurrency, you will receive this error message:
//
//	{
//	    "state": 1,
//	    "message": "The amount is too small"
//	}
//
// If technical work occurs or the payout is temporarily unavailable, you will receive this error messages:
//
//	{
//	    "state": 1,
//	    "message": "Server error"
//	}
//
//	{
//	    "state": 1,
//	    "message": "The service was not found"
//	}
//
// # Internal server error
//
// code : 500
//
// example responses:
//
//	{
//		"message": "Server error, #1",
//		"code": 500,
//		"error": null
//	}
func (m *Merchant) RefundBlockedAddress(request RefundBlockedAddressRequest) (*RefundBlockedAddressResponse, error) {
	httpResponse, err := m.sendPaymentRequest("POST", urlRefundBlockedAddress, request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	var response = struct {
		State   int                          `json:"state"`
		Result  RefundBlockedAddressResponse `json:"result"`
		Message string                       `json:"message"`
		Errors  struct {
			UUID    []string `json:"uuid"`
			OrderID []string `json:"order_id"`
			Address []string `json:"address"`
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
	errs = append(errs, response.Errors.Address...)

	if httpResponse.StatusCode != http.StatusOK || response.State != 0 || len(errs) > 0 {
		return nil, fmt.Errorf("error with status %s: %v", httpResponse.Status, strings.Join(errs, "; "))
	}

	return &response.Result, nil
}
