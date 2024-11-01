package cryptomus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// See "Transfer to personal wallet" https://doc.cryptomus.com/business/payouts/transfer-to-personal
//
//	{
//		"amount": "15",
//		"currency": "USDT"
//	}
type TransferRequest struct {
	// (Required) Amount to transfer
	Amount string `json:"amount"`
	// (Required) Currency code. Only cryptocurrency code is allowed.
	Currency string `json:"currency"`
}

// See "Transfer to personal wallet" https://doc.cryptomus.com/business/payouts/transfer-to-personal
//
//	{
//		"state": 0,
//		"result": {
//			"user_wallet_transaction_uuid": "26109ba0-b05b-4ee0-93d1-fd62c822ce95",
//			"user_wallet_balance": "15.00000000",
//			"merchant_transaction_uuid": "95bfcabb-a0ab-48f1-80b3-ce3bc2dbb653",
//			"merchant_balance": "20.00000000"
//		}
//	}
type TransferResponse struct {
	// methods.Uuid of Personal wallet transaction
	UserWalletTransactionUUID string `json:"user_wallet_transaction_uuid"`
	// Personal wallet balance
	UserWalletBalance string `json:"user_wallet_balance"`
	// methods.Uuid of Business wallet transaction
	MerchantTransactionUUID string `json:"merchant_transaction_uuid"`
	// Business wallet balance
	MerchantBalance string `json:"merchant_balance"`
}

// Transfer funds from your business wallet to your personal wallet
//
// See "Transfer to personal wallet" https://doc.cryptomus.com/business/payouts/transfer-to-personal
//
// # Possible errors
//
// # Validation errors
//
// code : 422
//
// example responses:
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
// If the payout is blocked on your user, you will receive this error message. This happens, for example, when you change your password and payouts are blocked for 24 hours.
//
//	{
//	    "state": 1,
//	    "message": "You are forbidden"
//	}
//
// If there is no active personal wallet with provided currency:
//
//	{
//	    "state": 1,
//	    "message": "No found user wallet"
//	}
//
// If there are not enough funds on your personal wallet to transfer the amount:
//
//	{
//	    "state": 1,
//	    "message": "Not enough funds"
//	}
//
// If the transfer amount is less than the minimum supported amount for cryptocurrency, you will receive this error message
//
//	{
//	    "state": 1,
//	    "message": "Minimum amount 0.5 USDT"
//	}
//
// If the transfer amount is greater than the maximum supported amount for cryptocurrency, you will receive this error message:
//
//	{
//	    "state": 1,
//	    "message": "Maximum amount 10000000 USDT"
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
func (m *Merchant) TransferToPersonalWallet(request TransferRequest) (*TransferResponse, error) {
	httpResponse, err := m.sendPayoutRequest("POST", urlTransferToPersonalWallet, request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	var response = struct {
		State   int              `json:"state"`
		Result  TransferResponse `json:"result"`
		Message string           `json:"message"`
		// If some parameter is required and not passed
		Errors *struct {
			Amount   []string `json:"amount"`
			Currency []string `json:"currency"`
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

	if httpResponse.StatusCode != http.StatusOK || response.State != 0 || len(errs) > 0 {
		return nil, fmt.Errorf("error creating invoice with status %s: %v", httpResponse.Status, strings.Join(errs, "; "))
	}

	return &response.Result, nil
}

// See "Transfer to business wallet" https://doc.cryptomus.com/business/payouts/transfer-to-business
//
// # Possible errors
//
// # Validation errors
//
// code : 422
//
// example response:
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
// If the payout is blocked on your user, you will receive this error message. This happens, for example, when you change your password and payouts are blocked for 24 hours.
//
//	{
//	    "state": 1,
//	    "message": "You are forbidden"
//	}
//
// If there is no active personal wallet with provided currency:
//
//	{
//	    "state": 1,
//	    "message": "No found user wallet"
//	}
//
// If there is no active business wallet with provided currency
//
//	{
//	    "state": 1,
//	    "message": "Merchant wallet not found"
//	}
//
// If there are not enough funds on your business wallet to transfer the amount
//
//	{
//	    "state": 1,
//	    "message": "Not enough funds"
//	}
//
// If the transfer amount is less than the minimum supported amount for cryptocurrency, you will receive this error message
//
//	{
//	    "state": 1,
//	    "message": "Minimum amount 0.5 USDT"
//	}
//
// If the transfer amount is greater than the maximum supported amount for cryptocurrency, you will receive this error message:
//
//	{
//	    "state": 1,
//	    "message": "Maximum amount 10000000 USDT"
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
func (m *Merchant) TransferToBusinessWallet(request TransferRequest) (*TransferResponse, error) {
	httpResponse, err := m.sendPayoutRequest("POST", urlCreateInvoice, request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	var response = struct {
		State   int              `json:"state"`
		Result  TransferResponse `json:"result"`
		Message string           `json:"message"`
		// If some parameter is required and not passed
		Errors *struct {
			Amount   []string `json:"amount"`
			Currency []string `json:"currency"`
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

	if httpResponse.StatusCode != http.StatusOK || response.State != 0 || len(errs) > 0 {
		return nil, fmt.Errorf("error creating invoice with status %s: %v", httpResponse.Status, strings.Join(errs, "; "))
	}

	return &response.Result, nil
}
