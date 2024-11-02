package cryptomus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// See "List of services" https://doc.cryptomus.com/business/payments/list-of-services
//
// See "List of services" https://doc.cryptomus.com/business/payouts/list-of-services
type Service struct {
	// Blockchain network code
	Network string `json:"network"`
	// Currency code
	Currency string `json:"currency"`
	// Whether a payment service is available (true / false)
	IsAvailable bool `json:"is_available"`
	Limit       struct {
		// Minimum amount available for payment/payout
		MinAmount string `json:"min_amount"`
		// Maximum amount available for payment/payout
		MaxAmount string `json:"max_amount"`
	} `json:"limit"`
	Commission struct {
		// Fixed fee amount
		FeeAmount string `json:"fee_amount"`
		// Percentage of Cryptomus payment commission
		Percent string `json:"percent"`
	} `json:"commission"`
}

// Returns a list of available payment services. Payment services store settings that are taken into account when creating an invoice. For example. currencies, networks, minimum and maximum limits, commissions.
//
// No parameters. Pass {}
//
// The reply comes with a collection of all available payment services.
//
// See "List of services" https://doc.cryptomus.com/business/payments/list-of-services
func (m *Merchant) ListPaymentServices() ([]Service, error) {
	httpResponse, err := m.sendPaymentRequest("POST", urlListPaymentServices, nil)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	var response = struct {
		State   int       `json:"state"`
		Result  []Service `json:"result"`
		Message string    `json:"message"`
		Code    int       `json:"code"`
		Error   string    `json:"error"`
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

	if httpResponse.StatusCode != http.StatusOK || response.State != 0 || len(errs) > 0 {
		return nil, fmt.Errorf("error with status %s: %v", httpResponse.Status, strings.Join(errs, "; "))
	}

	return response.Result, nil
}

// Returns a list of available payout services. Payout services store settings that are taken into account when creating a payout. For example. currencies, networks, minimum and maximum limits, commissions.
//
// No parameters. Pass {}
//
// The reply comes with a collection of all available payment services.
//
// See "List of services" https://doc.cryptomus.com/business/payouts/list-of-services
func (m *Merchant) ListPayoutServices() ([]Service, error) {
	httpResponse, err := m.sendPayoutRequest("POST", urlListPayoutServices, nil)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	var response = struct {
		State   int       `json:"state"`
		Result  []Service `json:"result"`
		Message string    `json:"message"`
		Code    int       `json:"code"`
		Error   string    `json:"error"`
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

	if httpResponse.StatusCode != http.StatusOK || response.State != 0 || len(errs) > 0 {
		return nil, fmt.Errorf("error with status %s: %v", httpResponse.Status, strings.Join(errs, "; "))
	}

	return response.Result, nil
}
