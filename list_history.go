package cryptomus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// See "Payment history" https://doc.cryptomus.com/business/payments/payment-history
//
// See "Payout history" https://doc.cryptomus.com/business/payouts/payout-history
//
// # Request example
//
//	{
//		"date_from": "2023-05-04 00:00:00"
//		"date_to": "2023-05-16 23:59:59"
//	}
type HistoryRequest struct {
	// (Optional) Filtering by creation date, from
	//    format: YYYY-MM-DD H:mm:ss
	//    default: null
	DateFrom *string `json:"date_from,omitempty"`
	// (Optional) Filtering by creation date, to
	//    format: YYYY-MM-DD H:mm:ss
	//    default: null
	DateTo *string `json:"date_to,omitempty"`
}

// See "Payment history" https://doc.cryptomus.com/business/payments/payment-history
//
// # Response example
//
//	{
//		"state": 0,
//		"result": {
//			"items": [{
//				"uuid": "ac1af391-8e98-4335-b9d7-7b6f6b40f268",
//				"order_id": "20fe59c4601dd174985e497e3f6bbcd2",
//				"amount": "20.00",
//				"payment_amount": "0.00000000",
//				"payer_amount": "0.00064860",
//				"discount_percent": 0,
//				"discount": "0.00000000",
//				"payer_currency": "BTC",
//				"currency": "USD",
//				"merchant_amount": null,
//				"comments": null,
//				"network": "btc",
//				"address": "bc1qxm6ehuy6mz2l2h3ag88frcjvl2xxlr9hvnq835",
//				"from": null,
//				"txid": null,
//				"payment_status": "cancel",
//				"url": "https://pay.cryptomus.com/pay/ac1af391-8e98-4335-b9d7-7b6f6b40f268",
//				"expired_at": 1689172103,
//				"status": "cancel",
//				"is_final": true,
//				"additional_data": null,
//				"created_at": "2023-07-12T16:28:24+03:00",
//				"updated_at": "2023-07-12T17:30:16+03:00"
//			}, {
//				"uuid": "1bb48358-2905-4e98-b681-5f1948e818d1",
//				"order_id": "a3329f462eb036dad12b5409147809a3",
//				"amount": "15.00",
//				"payment_amount": "0.00",
//				"payer_amount": "14.25",
//				"discount_percent": 5,
//				"discount": "0.75",
//				"payer_currency": "USDT",
//				"currency": "USDT",
//				"merchant_amount": "15.43500000",
//				"comments": null,
//				"network": "tron",
//				"address": "TSChodGNEJ6D31d9uueFxJAVH9NxiJjTwC",
//				"from": null,
//				"txid": null,
//				"payment_status": "cancel",
//				"url": "https://pay.cryptomus.com/pay/1bb48358-2905-4e98-b681-5f1948e818d1",
//				"expired_at": 1689099958,
//				"status": "cancel",
//				"is_final": true,
//				"additional_data": null,
//				"created_at": "2023-07-11T20:25:58+03:00",
//				"updated_at": "2023-07-11T21:26:18+03:00"
//			}, {
//				"uuid": "70b8db5c-b952-406d-af26-4e1c34c27f15",
//				"order_id": "65bbe87b4098c17a31cff3e71e515243",
//				"amount": "15.00",
//				"payment_amount": "0.00",
//				"payer_amount": "15.75",
//				"discount_percent": -5,
//				"discount": "-0.75",
//				"payer_currency": "USDT",
//				"currency": "USDT",
//				"merchant_amount": "15.43500000",
//				"comments": null,
//				"network": "tron",
//				"address": "TXhfYSWt2oKRrHAJVJeYRuit6ZzKuoEKXj",
//				"from": null,
//				"txid": null,
//				"payment_status": "cancel",
//				"url": "https://pay.cryptomus.com/pay/70b8db5c-b952-406d-af26-4e1c34c27f15",
//				"expired_at": 1689099831,
//				"status": "cancel",
//				"is_final": true,
//				"additional_data": null,
//				"created_at": "2023-07-11T20:23:52+03:00",
//				"updated_at": "2023-07-11T21:24:17+03:00"
//			},
//			...
//
// ],
//
//			"paginate": {
//				"count": 15,
//				"hasPages": true,
//				"nextCursor": "eyJpZCI6MjkxNTU0MywiX3BvaW50c1RvTmV4dEl0ZW1zIjp0cnVlfQ",
//				"previousCursor": null,
//				"perPage": 15
//			}
//		}
//	}
type paymentHistoryResponse struct {
	Items    []Invoice `json:"items"`
	Paginate paginate  `json:"paginate"`
}

// To get next/previous page entries, specify the next/previous cursor hash in the query parameters (?cursor=nextCursorHash)
//
// See "Payment history" https://doc.cryptomus.com/business/payments/payment-history
type paginate struct {
	// Number of items on the current page
	Count int `json:"count"`
	// Whether there are enough elements to split into multiple pages (true / false)
	HasPages bool `json:"hasPages"`
	// Cursor to next page
	NextCursor string `json:"nextCursor"`
	// Cursor to previous page
	PreviousCursor string `json:"previousCursor"`
	// Maximum number of elements per page
	PerPage int `json:"perPage"`
}

// See "Payment history" https://doc.cryptomus.com/business/payments/payment-history
func (m *Merchant) nextPaymentHistoryPage(currentPage *paymentHistoryResponse) (*paymentHistoryResponse, error) {
	if currentPage.Paginate.NextCursor == "" {
		return nil, nil
	}

	url := urlListPaymentHistory + "?cursor=" + currentPage.Paginate.NextCursor

	httpResponse, err := m.sendPaymentRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	var response = struct {
		State   int                    `json:"state"`
		Result  paymentHistoryResponse `json:"result"`
		Message string                 `json:"message"`
		Code    int                    `json:"code"`
		Error   string                 `json:"error"`
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

	return &response.Result, nil
}

// See "Payment history" https://doc.cryptomus.com/business/payments/payment-history
//
// # Response example
//
//	{
//		"state": 0,
//		"result": {
//			"items": [{
//				"uuid": "ac1af391-8e98-4335-b9d7-7b6f6b40f268",
//				"order_id": "20fe59c4601dd174985e497e3f6bbcd2",
//				"amount": "20.00",
//				"payment_amount": "0.00000000",
//				"payer_amount": "0.00064860",
//				"discount_percent": 0,
//				"discount": "0.00000000",
//				"payer_currency": "BTC",
//				"currency": "USD",
//				"merchant_amount": null,
//				"comments": null,
//				"network": "btc",
//				"address": "bc1qxm6ehuy6mz2l2h3ag88frcjvl2xxlr9hvnq835",
//				"from": null,
//				"txid": null,
//				"payment_status": "cancel",
//				"url": "https://pay.cryptomus.com/pay/ac1af391-8e98-4335-b9d7-7b6f6b40f268",
//				"expired_at": 1689172103,
//				"status": "cancel",
//				"is_final": true,
//				"additional_data": null,
//				"created_at": "2023-07-12T16:28:24+03:00",
//				"updated_at": "2023-07-12T17:30:16+03:00"
//			}, {
//				"uuid": "1bb48358-2905-4e98-b681-5f1948e818d1",
//				"order_id": "a3329f462eb036dad12b5409147809a3",
//				"amount": "15.00",
//				"payment_amount": "0.00",
//				"payer_amount": "14.25",
//				"discount_percent": 5,
//				"discount": "0.75",
//				"payer_currency": "USDT",
//				"currency": "USDT",
//				"merchant_amount": "15.43500000",
//				"comments": null,
//				"network": "tron",
//				"address": "TSChodGNEJ6D31d9uueFxJAVH9NxiJjTwC",
//				"from": null,
//				"txid": null,
//				"payment_status": "cancel",
//				"url": "https://pay.cryptomus.com/pay/1bb48358-2905-4e98-b681-5f1948e818d1",
//				"expired_at": 1689099958,
//				"status": "cancel",
//				"is_final": true,
//				"additional_data": null,
//				"created_at": "2023-07-11T20:25:58+03:00",
//				"updated_at": "2023-07-11T21:26:18+03:00"
//			}, {
//				"uuid": "70b8db5c-b952-406d-af26-4e1c34c27f15",
//				"order_id": "65bbe87b4098c17a31cff3e71e515243",
//				"amount": "15.00",
//				"payment_amount": "0.00",
//				"payer_amount": "15.75",
//				"discount_percent": -5,
//				"discount": "-0.75",
//				"payer_currency": "USDT",
//				"currency": "USDT",
//				"merchant_amount": "15.43500000",
//				"comments": null,
//				"network": "tron",
//				"address": "TXhfYSWt2oKRrHAJVJeYRuit6ZzKuoEKXj",
//				"from": null,
//				"txid": null,
//				"payment_status": "cancel",
//				"url": "https://pay.cryptomus.com/pay/70b8db5c-b952-406d-af26-4e1c34c27f15",
//				"expired_at": 1689099831,
//				"status": "cancel",
//				"is_final": true,
//				"additional_data": null,
//				"created_at": "2023-07-11T20:23:52+03:00",
//				"updated_at": "2023-07-11T21:24:17+03:00"
//			},
//			...
//
// ],
//
//			"paginate": {
//				"count": 15,
//				"hasPages": true,
//				"nextCursor": "eyJpZCI6MjkxNTU0MywiX3BvaW50c1RvTmV4dEl0ZW1zIjp0cnVlfQ",
//				"previousCursor": null,
//				"perPage": 15
//			}
//		}
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
// If the date_from or date_to parameter is passed in the wrong format:
//
//	{
//		"state": 1,
//		"errors": {
//			"date_from": ["validation.regex"]
//		}
//	}
func (m *Merchant) ListPaymentHistory(request HistoryRequest) ([]Invoice, error) {
	httpResponse, err := m.sendPaymentRequest("POST", urlListPaymentHistory, request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	var response = struct {
		State   int                    `json:"state"`
		Result  paymentHistoryResponse `json:"result"`
		Message string                 `json:"message"`
		Errors  struct {
			DateFrom []string `json:"date_from"`
			DateTo   []string `json:"date_to"`
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
		errs = append(errs, response.Errors.DateFrom...)
		errs = append(errs, response.Errors.DateTo...)
		if response.Error != "" {
			errs = append(errs, response.Error)
		}
	}
	if httpResponse.StatusCode != http.StatusOK || response.State != 0 || len(errs) > 0 {
		return nil, fmt.Errorf("error with status %s: %v", httpResponse.Status, strings.Join(errs, "; "))
	}

	var invoices []Invoice
	invoices = append(invoices, response.Result.Items...)
	page := response.Result

	for page.Paginate.NextCursor != "" {
		page, err := m.nextPaymentHistoryPage(&page)
		if err != nil {
			return nil, fmt.Errorf("error paging payment history: %w", err)
		}
		if page != nil {
			invoices = append(invoices, page.Items...)
		}
	}

	return invoices, nil
}

// payoutHistoryResponse represents the response structure for a payout history request.
//
// See "Payout history" https://doc.cryptomus.com/business/payouts/payout-history
//
// # Response example
//
//	{
//		"state": 0,
//		"result": {
//			"items": [{
//				"uuid": "a7c0caec-a594-4aaa-b1c4-77d511857594",
//				"amount": "3",
//				"currency": "USDT",
//				"network": "TRON",
//				"address": "TJ...",
//				"txid": null,
//				"status": "process",
//				"is_final": false,
//				"balance": "129.00000000",
//				"created_at": "2023-06-21T17:25:55+03:00",
//				"updated_at": "2023-06-21T17:34:38+03:00"
//			}, {
//				"uuid": "92c39264-d180-4503-9c16-ee16f083bbb8",
//				"amount": "5.40000000",
//				"currency": "DOGE",
//				"network": "doge",
//				"address": "DEw8CJLfxg9fhumeXP1zvVNjZicsqtDv7V",
//				"txid": "5e5810946152ea569d2a2aa9aa32a45c0e4223a4f9aad8e31d2fc660d2cdedb8",
//				"order_id": null,
//				"payment_status": null,
//				"status": "paid",
//				"is_final": true,
//				"balance": "26.77966652",
//				"created_at": "2023-07-21T17:25:55+03:00",
//				"updated_at": "2023-07-21T17:34:38+03:00"
//			},
//			...
//			],
//			"paginate": {
//				"count": 15,
//				"hasPages": true,
//				"nextCursor": "eyJpZCI6MjkxNTU0MywiX3BvaW50c1RvTmV4dEl0ZW1zIjp0cnVlfQ",
//				"previousCursor": null,
//				"perPage": 15
//			}
//		}
//	}
type payoutHistoryResponse struct {
	// Merchant's UUID
	MerchantUUID string `json:"merchant_uuid"`
	// Array of Payouts
	Items []Payout `json:"items"`
	// Data for pagination
	Paginate paginate `json:"paginate"`
}

// See "Payout history" https://doc.cryptomus.com/business/payouts/payout-history
func (m *Merchant) nextPayoutHistoryPage(currentPage *payoutHistoryResponse) (*payoutHistoryResponse, error) {
	if currentPage.Paginate.NextCursor == "" {
		return nil, nil
	}

	url := urlListPayoutHistory + "?cursor=" + currentPage.Paginate.NextCursor
	httpResponse, err := m.sendPayoutRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	var response = struct {
		State   int                   `json:"state"`
		Result  payoutHistoryResponse `json:"result"`
		Message string                `json:"message"`
		Code    int                   `json:"code"`
		Error   string                `json:"error"`
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

	return &response.Result, nil
}

// See "Payout history" https://doc.cryptomus.com/business/payouts/payout-history
//
// # Response example
//
//	{
//		"state": 0,
//		"result": {
//			"items": [{
//				"uuid": "a7c0caec-a594-4aaa-b1c4-77d511857594",
//				"amount": "3",
//				"currency": "USDT",
//				"network": "TRON",
//				"address": "TJ...",
//				"txid": null,
//				"status": "process",
//				"is_final": false,
//				"balance": "129.00000000",
//				"created_at": "2023-06-21T17:25:55+03:00",
//				"updated_at": "2023-06-21T17:34:38+03:00"
//			}, {
//				"uuid": "92c39264-d180-4503-9c16-ee16f083bbb8",
//				"amount": "5.40000000",
//				"currency": "DOGE",
//				"network": "doge",
//				"address": "DEw8CJLfxg9fhumeXP1zvVNjZicsqtDv7V",
//				"txid": "5e5810946152ea569d2a2aa9aa32a45c0e4223a4f9aad8e31d2fc660d2cdedb8",
//				"order_id": null,
//				"payment_status": null,
//				"status": "paid",
//				"is_final": true,
//				"balance": "26.77966652",
//				"created_at": "2023-07-21T17:25:55+03:00",
//				"updated_at": "2023-07-21T17:34:38+03:00"
//			},
//			...
//			],
//			"paginate": {
//				"count": 15,
//				"hasPages": true,
//				"nextCursor": "eyJpZCI6MjkxNTU0MywiX3BvaW50c1RvTmV4dEl0ZW1zIjp0cnVlfQ",
//				"previousCursor": null,
//				"perPage": 15
//			}
//		}
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
// If the date_from or date_to parameter is passed in the wrong format.
//
//	{
//		"state": 1,
//		"errors": {
//		"date_from": ["validation.regex"]
//		}
//	}
func (m *Merchant) ListPayoutHistory(request HistoryRequest) ([]Payout, error) {
	httpResponse, err := m.sendPayoutRequest("POST", urlListPayoutHistory, request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	var response = struct {
		State   int                   `json:"state"`
		Result  payoutHistoryResponse `json:"result"`
		Message string                `json:"message"`
		Errors  struct {
			DateFrom []string `json:"date_from"`
			DateTo   []string `json:"date_to"`
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
		errs = append(errs, response.Errors.DateFrom...)
		errs = append(errs, response.Errors.DateTo...)
		if response.Error != "" {
			errs = append(errs, response.Error)
		}
	}
	if httpResponse.StatusCode != http.StatusOK || response.State != 0 || len(errs) > 0 {
		return nil, fmt.Errorf("error with status %s: %v", httpResponse.Status, strings.Join(errs, "; "))
	}

	var payouts []Payout
	payouts = append(payouts, response.Result.Items...)
	page := response.Result

	for page.Paginate.NextCursor != "" {
		page, err := m.nextPayoutHistoryPage(&page)
		if err != nil {
			return nil, fmt.Errorf("error paging payout history: %w", err)
		}
		if page != nil {
			payouts = append(payouts, page.Items...)
		}
	}

	return payouts, nil
}

// See "List of recurring payments" https://doc.cryptomus.com/business/recurring/list
//
//	{
//		"state": 0,
//		"result": {
//		  "items": [
//			{
//			  "uuid": "bbe5ce96-1126-4843-a0d2-b432e77669c2",
//			  "name": "Access to personal account",
//			  "order_id": "1487555",
//			  "amount": "5.00",
//			  "currency": "USD",
//			  "payer_currency": "USDT",
//			  "payer_amount_usd": "5.00",
//			  "payer_amount": "5.00",
//			  "url_callback": null,"discount_days": "30",
//			  "discount_amount": "50.00",
//			  "end_of_discount": null,
//			  "period": "weekly",
//			  "status": "wait_accept",
//			  "url": "https://api.cryptomus.com/pay/bbe5ce96-1126-4843-a0d2-b432e77669c2",
//			  "last_pay_off": null
//			}
//		  ],
//		  "paginate": {
//			"count": 1,
//			"hasPages": false,
//			"nextCursor": null,
//			"previousCursor": null,
//			"perPage": 8
//		  }
//		}
//	  }
type recurringPaymentHistoryResponse struct {
	Items    []RecurringPayment `json:"items"`
	Paginate paginate           `json:"paginate"`
}

// See "List of recurring payments" https://doc.cryptomus.com/business/recurring/list
func (m *Merchant) nextRecurringPaymentHistoryPage(currentPage *recurringPaymentHistoryResponse) (*recurringPaymentHistoryResponse, error) {
	if currentPage.Paginate.NextCursor == "" {
		return nil, nil
	}

	url := urlListRecurringPayments + "?cursor=" + currentPage.Paginate.NextCursor

	httpResponse, err := m.sendPaymentRequest("POST", url, struct{}{})
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	var response = struct {
		State   int                             `json:"state"`
		Result  recurringPaymentHistoryResponse `json:"result"`
		Message string                          `json:"message"`
		Code    int                             `json:"code"`
		Error   string                          `json:"error"`
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
		return nil, fmt.Errorf("error retrieving next page of payment history with status %s: %v", httpResponse.Status, strings.Join(errs, "; "))
	}

	return &response.Result, nil
}

// See "List of recurring payments" https://doc.cryptomus.com/business/recurring/list
func (m *Merchant) ListRecurringPayments() ([]RecurringPayment, error) {
	httpResponse, err := m.sendPaymentRequest("POST", urlListRecurringPayments, struct{}{})
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	var response = struct {
		State   int                             `json:"state"`
		Result  recurringPaymentHistoryResponse `json:"result"`
		Message string                          `json:"message"`
		Code    int                             `json:"code"`
		Error   string                          `json:"error"`
	}{}

	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	var errs []string
	if httpResponse.StatusCode != http.StatusOK || response.State != 0 {
		if response.Message != "" {
			errs = append(errs, response.Message)
		}
		if response.Error != "" {
			errs = append(errs, response.Error)
		}
	}
	if httpResponse.StatusCode != http.StatusOK || response.State != 0 || len(errs) > 0 {
		return nil, fmt.Errorf("error with status %s: %v", httpResponse.Status, strings.Join(errs, "; "))
	}

	var recurringPayments []RecurringPayment
	recurringPayments = append(recurringPayments, response.Result.Items...)
	page := response.Result

	for page.Paginate.NextCursor != "" {
		page, err := m.nextRecurringPaymentHistoryPage(&page)
		if err != nil {
			return nil, fmt.Errorf("error paging recurring payments: %w", err)
		}
		if page != nil {
			recurringPayments = append(recurringPayments, page.Items...)
		}
	}

	return recurringPayments, nil
}

// See "Get orders list" https://doc.cryptomus.com/personal/converts/orders-list
//
// # Response example
//
//	{
//	  "state": 0,
//	  "result": {
//	    "items": [
//	      {
//	        "order_id": 49347,
//	        "convert_amount_from": "0.03700249",
//	        "convert_amount_to": "2476.39230892",
//	        "executed_amount_from": "0.03700249",
//	        "executed_amount_to": "2476.39230892",
//	        "convert_currency_from": "BTC",
//	        "convert_currency_to": "USDT",
//	        "type": "market",
//	        "status": "completed",
//	        "created_at": "2024-03-25 , 11:24:55",
//	        "current_rate": "66925.01798999",
//	        "completed_at": "2024-03-25 , 11:25:03"
//	      },
//	      ...
//	    ],
//	    "paginate": {
//	      "count": 15,
//	      "hasPages": true,
//	      "nextCursor": "eyJpdiI6ImRtQnpjUUorSkNPNXRtcGNzdVB6Mnc9PSIsInZhbHVlIjoiajArWmVOTTFEUzdJK2xUZmtLWm9Td2NEcTBmZmRQRUNkN203L0JZTUhicDVxYjBjQXpsZ0NDUkVBQjc3cFBXemZWbENQbmRVNEZPVm5FeCthM2tsV0VHODZ2aVp3MEg3SDc1VnRqeTNXbU8zTERORVp2L1JlQ0hHRWtjdkE2bTIiLCJtYWMiOiI1MjNjNjU0MzJkYmNmZWY1ZjQ4Zjc5ODMxOTExZmEzYWQzMDM3YTQxOTYyMmMyMWZlN2ZjYzJlZWE1MDAwMTY2IiwidGFnIjoiIn0=",
//	      "previousCursor": "eyJpdiI6IlRDMzlwWnI3MkZsTHRWdWdPUTZmenc9PSIsInZhbHVlIjoicWxiQTIvNzZycCtZKzFQRzlXQUFGcmdWLzRLMXhmUk1hZmZkSTBlNW9qSUkzVWplakdqUjlQNjNRdTM1TEZ5cTRVOW1acGs5VnRuMHdJVVZJN2QrcVE1aC8zZnFNWmRtZmpkaDloNFNkTTQ0YTF0N2tyNXh2c3VKNXkxSkRsRVUiLCJtYWMiOiJmNDU1ZTJiNjJlZjRmYWY3MDhmMTE3ZDVkNWE4OTUyNjg1MDM4OWUyZmRiM2I0YWZiMmUwNDRkMDRjMTUzNzJkIiwidGFnIjoiIn0=",
//	      "perPage": 15
//	    }
//	  }
//	}
type listOrdersResponse struct {
	Items    []MarketOrder `json:"items"`
	Paginate paginate      `json:"Paginate"`
}


// Available options for type:
//   - market
//   - limit
//
// Available options for status:
//   - active
//   - completed
//   - partially_completed
//   - cancelled
//   - expired
//   - failed
func (u *User) nextOrderHistoryPage(cursor, orderType, orderStatus string) (*listOrdersResponse, error) {
	url := urlListOrderHistory
	if cursor != "" {
		url = url + "?cursor=" + cursor
	}
	if orderType != "" {
		url = url + "?type=" + orderType
	}
	if orderStatus != "" {
		url = url + "?status=" + orderStatus
	}

	httpResponse, err := u.sendPaymentRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	var response = struct {
		State   int                `json:"state"`
		Result  listOrdersResponse `json:"result"`
		Message string             `json:"message"`
		Code    int                `json:"code"`
		Error   string             `json:"error"`
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

	return &response.Result, nil

}


// See "Get orders list" https://doc.cryptomus.com/personal/converts/orders-list
//
// # Response example
//
//	{
//	  "state": 0,
//	  "result": {
//	    "items": [
//	      {
//	        "order_id": 49347,
//	        "convert_amount_from": "0.03700249",
//	        "convert_amount_to": "2476.39230892",
//	        "executed_amount_from": "0.03700249",
//	        "executed_amount_to": "2476.39230892",
//	        "convert_currency_from": "BTC",
//	        "convert_currency_to": "USDT",
//	        "type": "market",
//	        "status": "completed",
//	        "created_at": "2024-03-25 , 11:24:55",
//	        "current_rate": "66925.01798999",
//	        "completed_at": "2024-03-25 , 11:25:03"
//	      },
//	      ...
//	    ],
//	    "paginate": {
//	      "count": 15,
//	      "hasPages": true,
//	      "nextCursor": "eyJpdiI6ImRtQnpjUUorSkNPNXRtcGNzdVB6Mnc9PSIsInZhbHVlIjoiajArWmVOTTFEUzdJK2xUZmtLWm9Td2NEcTBmZmRQRUNkN203L0JZTUhicDVxYjBjQXpsZ0NDUkVBQjc3cFBXemZWbENQbmRVNEZPVm5FeCthM2tsV0VHODZ2aVp3MEg3SDc1VnRqeTNXbU8zTERORVp2L1JlQ0hHRWtjdkE2bTIiLCJtYWMiOiI1MjNjNjU0MzJkYmNmZWY1ZjQ4Zjc5ODMxOTExZmEzYWQzMDM3YTQxOTYyMmMyMWZlN2ZjYzJlZWE1MDAwMTY2IiwidGFnIjoiIn0=",
//	      "previousCursor": "eyJpdiI6IlRDMzlwWnI3MkZsTHRWdWdPUTZmenc9PSIsInZhbHVlIjoicWxiQTIvNzZycCtZKzFQRzlXQUFGcmdWLzRLMXhmUk1hZmZkSTBlNW9qSUkzVWplakdqUjlQNjNRdTM1TEZ5cTRVOW1acGs5VnRuMHdJVVZJN2QrcVE1aC8zZnFNWmRtZmpkaDloNFNkTTQ0YTF0N2tyNXh2c3VKNXkxSkRsRVUiLCJtYWMiOiJmNDU1ZTJiNjJlZjRmYWY3MDhmMTE3ZDVkNWE4OTUyNjg1MDM4OWUyZmRiM2I0YWZiMmUwNDRkMDRjMTUzNzJkIiwidGFnIjoiIn0=",
//	      "perPage": 15
//	    }
//	  }
//	}
func (u *User) ListOrderHistory(orderType, orderStatus string) ([]MarketOrder, error) {
	url := urlListOrderHistory
	if orderType != "" {
		url = url + "?type=" + orderType
	}
	if orderStatus != "" {
		url = url + "?status=" + orderStatus
	}

	httpResponse, err := u.sendPaymentRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	var response = struct {
		State   int                `json:"state"`
		Result  listOrdersResponse `json:"result"`
		Message string             `json:"message"`
		Code    int                `json:"code"`
		Error   string             `json:"error"`
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

	var orders []MarketOrder
	orders = append(orders, response.Result.Items...)
	page := &response.Result
	for page.Paginate.NextCursor != "" {
		page, err = u.nextOrderHistoryPage(page.Paginate.NextCursor, orderType, orderStatus)
		if err != nil {
			return nil, fmt.Errorf("error paging orders history: %w", err)
		}
		if page != nil {
			orders = append(orders, page.Items...)
		}
	}
	return orders, nil
}