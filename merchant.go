package cryptomus

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// You need a merchant with different API keys for accepting payment and making payouts.
//
// See "Getting API keys" https://doc.cryptomus.com/business/general/getting-api-keys
type Merchant struct {
	MerchantUUID, PaymentAPIKey, PayoutAPIKey string
	client                                    *http.Client
}

// NewMerchant creates a merchant with different API keys for accepting payment and making payouts.
//
// See "Getting API keys" https://doc.cryptomus.com/business/general/getting-api-keys
func NewMerchant(merchantUUID, paymentAPIKey, PayoutAPIKey string) *Merchant {
	return &Merchant{
		MerchantUUID:  merchantUUID,
		PaymentAPIKey: paymentAPIKey,
		PayoutAPIKey:  PayoutAPIKey,
		client:        &http.Client{Timeout: 10 * time.Second},
	}
}

// signPaymentPayload generates MD5 hash of the body of the POST request encoded in base64 and combined with your payment API key.
//
// See "Request format" https://doc.cryptomus.com/business/general/request-format
func (m *Merchant) signPaymentPayload(jsonData []byte) (string, error) {
	base64Data := base64.StdEncoding.EncodeToString(jsonData)
	hash := md5.Sum([]byte(base64Data + m.PaymentAPIKey))
	return hex.EncodeToString(hash[:]), nil
}

// signPayoutPayload generates MD5 hash of the body of the POST request encoded in base64 and combined with your payout API key.
//
// See "Request format" https://doc.cryptomus.com/business/general/request-format
func (m *Merchant) signPayoutPayload(jsonData []byte) (string, error) {
	base64Data := base64.StdEncoding.EncodeToString(jsonData)
	hash := md5.Sum([]byte(base64Data + m.PayoutAPIKey))
	return hex.EncodeToString(hash[:]), nil
}

func (m *Merchant) sendPaymentRequest(method, url string, request any) (*http.Response, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request data: %w", err)
	}

	httpRequest, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	signature, err := m.signPaymentPayload(jsonData)
	if err != nil {
		return nil, fmt.Errorf("error generating signature: %w", err)
	}

	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("merchant", m.MerchantUUID)
	httpRequest.Header.Set("sign", signature)

	httpResponse, err := m.client.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	return httpResponse, nil
}

func (m *Merchant) sendPayoutRequest(method, url string, request any) (*http.Response, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request data: %w", err)
	}

	httpRequest, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	signature, err := m.signPayoutPayload(jsonData)
	if err != nil {
		return nil, fmt.Errorf("error generating signature: %w", err)
	}

	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("merchant", m.MerchantUUID)
	httpRequest.Header.Set("sign", signature)

	httpResponse, err := m.client.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	return httpResponse, nil
}
