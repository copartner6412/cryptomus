package cryptomus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// QRCodeForStaticWalletRequest represents the request to generate a QR code for a static wallet.
//
// See "Generate a QR-code" https://doc.cryptomus.com/business/payments/qr-code-pay-form
//
// # Request example
//
//	{
//	    "wallet_address_uuid": "8b03432e-385b-4670-8d06-064591096795",
//	}
type QRCodeForStaticWalletRequest struct {
	// (Required) uuid of a static wallet
	WalletAddressUUID string `json:"wallet_address_uuid"`
}

// QRCodeForInvoiceRequest represents the request to generate a QR code for an invoice.
//
// See "Generate a QR-code" https://doc.cryptomus.com/business/payments/qr-code-pay-form
//
//	{
//	    "merchant_payment_uuid": "8b03432e-385b-4670-8d06-064591096795",
//	}
type QRCodeForInvoiceRequest struct {
	// (Required) Invoice uuid
	MerchantPaymentUUID string `json:"merchant_payment_uuid"`
}

// QRCodeResponse represents the response with the QR code image in base64 format.
//
// # Response example
//
//	{
//	    "state": 0,
//	    "result": {
//	        "image": "data:image/png;base64,iVBORw0KGgoAAA..."
//		}
//	}
type QRCodeResponse struct {
	// Base64 encode QR-code image
	Image string `json:"image"`
}

// GenerateQRCodeForStaticWallet is a payment method that generates a QR-code for a static wallet address. Scanning it, the user will receive the address for depositing funds.
//
// See "Generate a QR-code" https://doc.cryptomus.com/business/payments/qr-code-pay-form
func (m *Merchant) GenerateQRCodeForStaticWallet(request QRCodeForStaticWalletRequest) (*QRCodeResponse, error) {
	httpResponse, err := m.sendPaymentRequest("POST", urlGenerateQRCodeForStaticWallet, request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	var response = struct {
		State   int            `json:"state"`
		Result  QRCodeResponse `json:"result"`
		Message string         `json:"message"`
		// If some parameter is required and not passed
		Errors struct {
			WalletAddressUUID []string `json:"wallet_address_uuid"`
		} `json:"errors"`
		Code  int    `json:"code"`
		Error string `json:"error"`
	}{}
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	var errs []string
	if httpResponse.StatusCode != http.StatusOK || response.State != 0 {
		if response.Message != "" {
			errs = append(errs, response.Message)
		}
		errs = append(errs, response.Errors.WalletAddressUUID...)
		if response.Error != "" {
			errs = append(errs, response.Error)
		}
	}
	if len(errs) > 0 {
		return nil, fmt.Errorf("error creating QR code for static wallet with status %s: %v", httpResponse.Status, strings.Join(errs, "; "))
	}

	return &response.Result, nil
}

// GenerateQRCodeForInvoice is a payment method that generates a QR-code for an invoice address. Scanning it, the user will receive the address for depositing funds.
//
// See "Generate a QR-code" https://doc.cryptomus.com/business/payments/qr-code-pay-form
func (m *Merchant) GenerateQRCodeForInvoice(request QRCodeForInvoiceRequest) (*QRCodeResponse, error) {
	httpResponse, err := m.sendPaymentRequest("POST", urlGenerateQRCodeForStaticWallet, request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	var response = struct {
		State   int            `json:"state"`
		Result  QRCodeResponse `json:"result"`
		Message string         `json:"message"`
		Errors  struct {
			MerchantPaymentUUID []string `json:"merchant_payment_uuid"`
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
	errs = append(errs, response.Errors.MerchantPaymentUUID...)
	if response.Error != "" {
		errs = append(errs, response.Error)
	}

	if httpResponse.StatusCode != http.StatusOK || response.State != 0 || len(errs) > 0 {
		return nil, fmt.Errorf("error creating QR code for invoice with status %s: %v", httpResponse.Status, strings.Join(errs, "; "))
	}

	return &response.Result, nil
}
