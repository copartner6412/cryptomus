package cryptomus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (m *Merchant) CancelRecurringPayment(request RecordID) (*RecurringPayment, error) {
	httpResponse, err := m.sendPaymentRequest("POST", urlCancelRecurringPayment, request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	var response = struct {
		State   int              `json:"state"`
		Result  RecurringPayment `json:"result"`
		Message string           `json:"message"`
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
		return nil, fmt.Errorf("error canceling recurring payment information with status %s: %v", httpResponse.Status, strings.Join(errs, "; "))
	}

	return &response.Result, nil
}
