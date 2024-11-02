package cryptomus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// The payouts through API are made only from your business wallets balances.
//
// Payouts can be made in different ways:
//  1. You can choose to receive the payout in a specific cryptocurrency and the payout will then be automatically processed in that specific cryptocurrency. To do so, ensure that you have sufficient balance in that particular currency to cover all associated fees.
//  2. Alternatively, you have the option to specify the payout amount in a fiat currency. In this case, the amount will be automatically converted to a specific cryptocurrency from your available balance. For instance, if you request a payout of 20 USD in LTC, the equivalent value will be deducted from your LTC balance. It is important to have enough funds in the corresponding cryptocurrency to cover all applicable fees.
//  3. Another possibility is to specify the payout amount in a fiat currency, which will be automatically converted to a specific cryptocurrency using your USDT balance. This option is particularly useful when you have autoconvert enabled, as funds from your invoices are automatically converted to USDT. For example, if you want to make a payout of 20 USD in LTC but only have a balance in USDT, make sure you have sufficient USDT funds to cover all fees.
//  4. Additionally, you can choose to specify the payout amount in any cryptocurrency of your preference. The payout will then be automatically processed in that specific cryptocurrency, utilizing your available USDT balance. It is crucial to have enough USDT balance to cover all associated fees.
//
// See "Creating a payout" https://doc.cryptomus.com/business/payouts/creating-payout
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
//	{
//		"state": 1,
//		"errors": {
//			"amount": ["validation.required"]
//		}
//	}
//
// If the payout is blocked on your merchant or user, you will receive this error message. This happens, for example, when you change your password and payouts are blocked for 24 hours.
//
//	{
//	    "state": 1,
//	    "message": "You are forbidden"
//	}
//
// If we don’t support the currency code you provided in currency parameter, you will receive this error message:
//
//	{
//	    "state": 1,
//	    "message": "The currency was not found"
//	}
//
// If we can not find the service for executing payout, you will receive this error message. You have to make sure that the cryptocurrency is supported. Check the correctness of the entered values in. currency, network, from_currency, to_currency:
//
//	{
//	    "state": 1,
//	    "message": "The service was not found"
//	}
//
// If the withdrawal amount is too small, you will receive this error message:
//
//	{
//	    "state": 1,
//	    "message": "The withdrawal amount is too small"
//	}
//
// If you will not have an active merchant wallet with a cryptocurrency of the payout, you will receive this error:
//
//	{
//	    "state": 1,
//	    "message": "Wallet not found"
//	}
//
// If the payout amount, including all commissions, is more than your balance, you will receive this error message:
//
//	{
//	    "state": 1,
//	    "message": "Not enough funds"
//	}
//
// If we don’t find the payout service of the currency that was passed in to_currency parameter, you will receive this error message:
//
//	{
//	    "state": 1,
//	    "message": "Not found service to_currency"
//	}
//
// If accidentally something happens when converting to to_currency, you will receive this error message:
//
//	{
//	    "state": 1,
//	    "message": "Error convert to_currency"
//	}
//
// If the payout amount is less than the minimum supported amount for cryptocurrency, you will receive this error message:
//
//	{
//	    "state": 1,
//	    "message": "Minimum amount 0.5 USDT"
//	}
//
// If the payout amount is greater than the maximum supported amount for cryptocurrency, you will receive this error message:
//
//	{
//	    "state": 1,
//	    "message": "Maximum amount 10000000 USDT"
//	}
//
// If you don’t have enough funds to make a payout, you will receive this error message:
//
//	{
//	    "state": 1,
//	    "message": "Not enough funds"
//	}
//
// If you don't have enough funds to convert the payout amount from from_currency, you will receive this error message:
//
//	{
//	    "state": 1,
//	    "message": "Not enough balance for convert USDT to LTC"
//	}
//
// If technical work occurs and the payout is temporarily unavailable, you will receive this error message
//
//	{
//	    "state": 1,
//	    "message": "The terminal was not found"
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
func (m *Merchant) CreatePayout(request Withdrawal) (*Payout, error) {
	httpResponse, err := m.sendPayoutRequest("POST", urlCreatePayout, request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	var response = struct {
		State   int    `json:"state"`
		Result  Payout `json:"result"`
		Message string `json:"message"`
		// If some parameter is required and not passed
		Errors struct {
			Amount     []string `json:"amount"`
			Currency   []string `json:"currency"`
			OrderID    []string `json:"order_id"`
			Address    []string `json:"address"`
			IsSubtract []string `json:"is_subtract"`
			Network    []string `json:"network"`
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
	errs = append(errs, response.Errors.Address...)
	errs = append(errs, response.Errors.IsSubtract...)
	errs = append(errs, response.Errors.Network...)

	if httpResponse.StatusCode != http.StatusOK || response.State != 0 || len(errs) > 0 {
		return nil, fmt.Errorf("error with status %s: %v", httpResponse.Status, strings.Join(errs, "; "))
	}

	return &response.Result, nil
}
