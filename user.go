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

type User struct {
	UserID, PaymentAPIKey, PayoutAPIKey string
	client                              *http.Client
}

// You need to release a different API key for accepting payment and making payouts
//
// See "Getting API keys" https://doc.cryptomus.com/personal/general/getting-api-keys
func NewUser(userID, paymentAPIKey, payoutAPIKey string) *User {
	return &User{
		UserID:        userID,
		PaymentAPIKey: paymentAPIKey,
		PayoutAPIKey:  payoutAPIKey,
		client:        &http.Client{Timeout: 10 * time.Second},
	}
}

// signPaymentPayload generates MD5 hash of the body of the POST request encoded in base64 and combined with your payment API key.
//
// See "Request format" https://doc.cryptomus.com/personal/general/request-format
func (u *User) signPaymentPayload(jsonData []byte) (string, error) {
	base64Data := base64.StdEncoding.EncodeToString(jsonData)
	hash := md5.Sum([]byte(base64Data + u.PaymentAPIKey))
	return hex.EncodeToString(hash[:]), nil
}

// signPayoutPayload generates MD5 hash of the body of the POST request encoded in base64 and combined with your payment API key.
//
// See "Request format" https://doc.cryptomus.com/personal/general/request-format
func (u *User) signPayoutPayload(jsonData []byte) (string, error) {
	base64Data := base64.StdEncoding.EncodeToString(jsonData)
	hash := md5.Sum([]byte(base64Data + u.PayoutAPIKey))
	return hex.EncodeToString(hash[:]), nil
}

func (u *User) sendPaymentRequest(method, url string, request any) (*http.Response, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	httpRequest, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	signature, err := u.signPaymentPayload(jsonData)
	if err != nil {
		return nil, fmt.Errorf("error signing request payload: %w", err)
	}

	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("userId", u.UserID)
	httpRequest.Header.Set("sign", signature)

	httpResponse, err := u.client.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	return httpResponse, nil
}

func (u *User) sendPayoutRequest(method, url string, request any) (*http.Response, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request payload: %w", err)
	}

	httpRequest, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	signature, err := u.signPayoutPayload(jsonData)
	if err != nil {
		return nil, fmt.Errorf("error signing request payload: %w", err)
	}

	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("userId", u.UserID)
	httpRequest.Header.Set("sign", signature)

	httpResponse, err := u.client.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	return httpResponse, nil
}
