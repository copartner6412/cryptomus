package cryptomus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Required fields:
//   - Currency
//   - Network
//   - OrderID
//
// See "Creating a Static wallet" https://doc.cryptomus.com/business/payments/creating-static
//
// # Request example
//
//	{
//	    "currency": "USDT",
//	    "network": "tron",
//	    "order_id": "1",
//	    "url_callback": "https://your.site/callback"
//	}
type StaticWalletRequest struct {
	// (Required) Currency code
	Currency string `json:"currency"`
	// (Required) Blockchain network code
	Network string `json:"network"`
	// (Required) Order ID in your system
	//
	// The parameter should be a string consisting of alphabetic characters, numbers, underscores, and dashes. It should not contain any spaces or special characters.
	//
	// The order_id must be unique within the merchant invoices/static wallets/recurrence payments
	//
	// When we find an existing invoice with order_id, we return its details, a new invoice will not be created.
	//    min: 1
	//    max: 100
	//    alpha_dash
	OrderID string `json:"order_id"`
	// (Optional) URL, to which the webhook will be sent after each top-up of the wallet
	//    min: 6
	//    max: 255
	//    url
	//    default: null
	URLCallback *string `json:"url_callback,omitempty"`
	// (Optional) The merchant who makes the request connects to a referrer by code.
	//
	// For example, you are an application that generates invoices via the Cryptomus API and your customers are other stores.
	//
	// They enter their api key and merchant id in your application, and you send requests with their credentials and passing your referral code.
	//
	// Thus, your clients become referrals on your Cryptomus account and you will receive income from their turnover.
	//
	//    default: null
	FromReferralCode *string `json:"from_referral_code,omitempty"`
}

// See "Creating a Static wallet" https://doc.cryptomus.com/business/payments/creating-static
//
// # Response example
//
//	{
//		"state": 0,
//		"result": {
//			"wallet_uuid": "de15b0f6-883f-4585-b27b-73a648044a92",
//			"uuid": "87961ae5-80c5-413a-a4fe-d38199894940",
//			"address": "TTEtddVZyNtLD9wbq4PzomjBhtxenSMXbb",
//			"network": "tron",
//			"currency": "USDT",
//			"url": "https://pay.cryptomus.com/wallet/3901446a-4b74-4796-b50a-14e14dafe3ed"
//		}
//	}
type StaticWalletResponse struct {
	// uuid of the merchant's wallet
	WalletUUID string `json:"wallet_uuid"`
	// uuid of the wallet in a particular network
	UUID string `json:"uuid"`
	// Wallet address in a particular network
	Address string `json:"address"`
	// Wallet network code
	Network string `json:"network"`
	// Wallet network currency
	Currency string `json:"currency"`
	// Link to payment form
	URL string `json:"url"`
}

// CreateStaticWallet is a payment method that creates a new static wallet for merchant on Cryptomus (Suitable for balance top-up)
//
// You can create a static address in a specific currency and network. The address will be attached to the order_id
//
// All transactions sent to this address will be credited regardless of the amount.
//
// See "Creating a Static wallet" https://doc.cryptomus.com/business/payments/creating-static
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
//			"currency": ["validation.required"]
//		}
//	}
//
// If your payments are blocked. In this case you need to contact support.
//
//	{
//	    "state": 1,
//	    "message": "You are forbidden"
//	}
//
// If we don’t find the payment service with the currency and network:
//
//	{
//	    "state": 1,
//	    "message": "The service was not found"
//	}
//
// If technical work occurs and the payment is temporarily unavailable, you can receive this error message:
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
//	{
//	    "state": 1,
//	    "message": "Wallet not found"
//	}
func (m *Merchant) CreateStaticWallet(request StaticWalletRequest) (*StaticWalletResponse, error) {
	httpResponse, err := m.sendPaymentRequest("POST", urlCreateStaticWallet, request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	var response = struct {
		State   int                  `json:"state"`
		Result  StaticWalletResponse `json:"result"`
		Message string               `json:"message"`
		// If some parameter is required and not passed
		Errors struct {
			Currency []string `json:"currency"`
			Network  []string `json:"network"`
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
	errs = append(errs, response.Errors.Currency...)
	errs = append(errs, response.Errors.Network...)
	errs = append(errs, response.Errors.OrderID...)

	if httpResponse.StatusCode != http.StatusOK || response.State != 0 || len(errs) > 0 {
		return nil, fmt.Errorf("error creating static wallet with status %s: %v", httpResponse.Status, strings.Join(errs, "; "))
	}

	return &response.Result, nil
}
