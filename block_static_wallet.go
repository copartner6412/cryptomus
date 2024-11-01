package cryptomus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// You need to pass one of the required parameters, if you pass both, the account will be identified by order_id
//
// See "Block static wallet" https://doc.cryptomus.com/business/payments/block-wallet
//
// # Request example
//
//	{
//	    "order_id": "1"
//	}
type BlockStaticWalletRequest struct {
	RecordID
	// (Optional) Refund all incoming payments to sender’s address
	//    default: false
	IsForceRefund *bool `json:"is_force_refund,omitempty"`
}

// See "Block static wallet" https://doc.cryptomus.com/business/payments/block-wallet
//
// # Response example
//
//	{
//		"state": 0,
//		"result": {
//			"uuid": "fcc40793-39f9-4fa9-85b2-93148039a72b",
//			"status": "blocked"
//		}
//	}
type BlockStaticWalletResponse struct {
	// uuid of the static wallet
	UUID string `json:"uuid"`
	// Status of the static wallet.
	//
	// Available options:
	//  - blocked
	//  - active
	//  - in_active
	Status string `json:"status"`
}

// When you need to block your clients static wallet, all the further payments will not be credited to his balance. You can make a refund of this funds only once. The funds will be returned to the addresses from which they came.
//
// See "Block static wallet" https://doc.cryptomus.com/business/payments/block-wallet
//
// # Possible errors
//
// # Validation errors
//
// code : 422
//
//	{
//		"state": 1
//	}
//
//	{
//		"state": 1,
//		"errors": {
//			"uuid": ["validation.required_without"],
//			"order_id": ["validation.required_without"]
//		}
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
func (m *Merchant) BlockStaticWallet(request BlockStaticWalletRequest) (*BlockStaticWalletResponse, error) {
	httpResponse, err := m.sendPaymentRequest("POST", urlBlockStaticWallet, request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	var response = struct {
		State   int                       `json:"state"`
		Result  BlockStaticWalletResponse `json:"result"`
		Message string                    `json:"message"`
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
		return nil, fmt.Errorf("error blocking static wallet with status %s: %v", httpResponse.Status, strings.Join(errs, "; "))
	}

	return &response.Result, nil
}
