package cryptomus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Discount:
//
// The first period in days (the discount_days parameter), which will cost a different amount (the discount_amount parameter). For example, you can set a discount_amount for the first 30 days and a monthly period, then the user will pay the first payment with a lower price, the second and subsequent payments will be at full price.
//
// The discount amount must be in the currency specified by the parameter ‘currency’.
//
// See "Creating recurring payment" https://doc.cryptomus.com/business/recurring/creating
//
// # Request example
//
// Below is a sample example with the minimum required number of parameters In this case, an invoice will be created with the fiat currency USD. On the payment page, the user will be able to select the cryptocurrency and network to pay this bill.
//
//	{
//		"amount": "15",
//		"currency": "USDT",
//		"name": "Recurring payment",
//		"period": "monthly"
//	}
type RecurringInvoice struct {
	// (Required) Recurring payment amount
	Amount string `json:"amount"`
	// (Required) Currency code
	Currency string `json:"currency"`
	// (Required) Recurring payment name
	//    min: 3
	//    max: 60
	Name string `json:"name"`
	// (Required) Recurring payment period
	// Available options:
	//  - weekly
	//  - monthly
	//  - three_month
	Period string `json:"period"`
	// (Optional) Currency code for accepting payments
	//
	// The parameter is used to specify the target currency for converting the recurrent payment amount.
	//
	// For example, to create an recurrent payment for 20 USD in bitcoin:
	//    amount: 20
	//    currency: USD
	//    to_currency: BTC
	// The API will convert 20 USD amount to its equivalent in BTC based on the current exchange rate and the user will pay in BTC
	//
	// The to_currency should always be the cryptocurrency code, not a fiat currency code.
	ToCurrency *string `json:"to_currency,omitempty"`
	// (Optional) Order ID in your system
	//    default: null
	OrderID *string `json:"order_id,omitempty"`
	// (Optional) Url to which webhooks with payment status will be sent
	//    default: null
	URLCallback *string `json:"url_callback,omitempty"`
	// (Optional) Url to which webhooks with payment status will be sent
	//    default: null
	DiscountDays *int `json:"discount_days,omitempty"`
	// (Optional) Discount period days (required with 'discount_amount')
	//    min: 1
	//    max: 365
	//    default: 0
	DiscountAmount *string `json:"discount_amount,omitempty"`
	// (Optional) Discount amount (required with 'discount_days').Here the amount in the currency of the parameter ‘currency’
	//    default: null
	AdditionalData *string `json:"additional_data,omitempty"`
	// Additional recurring payment details
	//    default: null
}

// Recurring payments in cryptocurrency are a way to automate regular transactions using digital assets. They can be useful for subscription-based services, donations, memberships, and other recurring payments.
//
// To use recurring payments, you need to create a payment that specifies the amount, currency, and frequency of the payments, and then share it with your payers. The payer will be redirected to the cryptomus website, where he will need to log in to confirm the payment plan and make the first payment. After that, payments will be made automatically according to the plan.
//
// See "Creating recurring payment" https://doc.cryptomus.com/business/recurring/creating
//
// # Examples
//
// Request data example to create a recurrence payment for 15 USD
//
//	{
//		"amount": "15",
//		"currency": "USD",
//		"name": "Recurring payment",
//		"period": "monthly",
//		"url_callback": "https://your.site/callback"
//	}
//
// Response:
//
//	{
//		"state": 0,
//		"result": {
//			"uuid": "1d4c1ef7-1dcd-4f98-9433-2e19d2e62820",
//			"name": "Recurring payment",
//			"order_id": null,
//			"amount": "15",
//			"currency": "USD",
//			"payer_currency": null,
//			"payer_amount_usd": "0.00",
//			"payer_amount": null,
//			"url_callback": "https://your.site/callback",
//			"period": "monthly",
//			"status": "wait_accept",
//			"url": "https://pay.cryptomus.com/recurring/1d4c1ef7-1dcd-4f98-9433-2e19d2e62820",
//			"last_pay_off": null
//		}
//	}
//
// Request data example to create a recurrence payment for 15 USDT (only USDT will be allowed for payment)
//
//	{
//		"amount": "15",
//		"currency": "USDT",
//		"name": "Recurring payment",
//		"period": "monthly",
//		"url_callback": "https://your.site/callback"
//	}
//
// Response:
//
//	{
//		"state": 0,
//		"result": {
//			"uuid": "ec04dc64-fea4-445a-9698-59baa474317a",
//			"name": "Recurring payment",
//			"order_id": null,
//			"amount": "15",
//			"currency": "USDT",
//			"payer_currency": "USDT",
//			"payer_amount_usd": "15.00",
//			"payer_amount": "15.00000000",
//			"url_callback": "https://your.site/callback",
//			"period": "monthly",
//			"status": "wait_accept",
//			"url": "https://pay.cryptomus.com/recurring/ec04dc64-fea4-445a-9698-59baa474317a",
//			"last_pay_off": null
//		}
//	}
//
// Request data example to create a recurrence payment for 15 USD and allow only BTC currency for payment
//
//	{
//		"amount": "15",
//		"currency": "USD",
//		"to_currency": "BTC",
//		"name": "Recurring payment",
//		"period": "monthly",
//		"url_callback": "https://your.site/callback"
//	}
//
// Response:
//
//	{
//		"state": 0,
//		"result": {
//			"uuid": "e7d9c265-3849-4ff8-84c1-15435368f05f",
//			"name": "Recurring payment",
//			"order_id": null,
//			"amount": "15",
//			"currency": "USD",
//			"payer_currency": "BTC",
//			"payer_amount_usd": "15.00",
//			"payer_amount": "0.00051495",
//			"url_callback": "https://your.site/callback",
//			"period": "monthly",
//			"status": "wait_accept",
//			"url": "https://pay.cryptomus.com/recurring/e7d9c265-3849-4ff8-84c1-15435368f05f",
//			"last_pay_off": null
//		}
//	}
//
// Request data example to create a recurring payment in the amount of $15. The first 30 days the price will be $1, after that - $15:
//
//	{
//		"amount": "15",
//		"currency": "USD",
//		"name": "Recurring payment",
//		"period": "monthly",
//		"discount_days": 30,
//	  "discount_amount": "1"
//		"url_callback": "https://your.site/callback"
//	}
//
// Response:
//
//	{
//		"state": 0,
//		"result": {
//			"uuid": "1ee00fbd-42e7-4653-b316-1ab12f8a447b",
//			"name": "Recurring payment",
//			"order_id": null,
//			"amount": "15",
//			"currency": "USD",
//			"payer_currency": null,
//			"payer_amount_usd": "0.00",
//			"payer_amount": null,
//			"url_callback": "https://your.site/callback",
//			"discount_days": "30",
//	                "discount_amount": "1",
//			"end_of_discount": "2023-07-11T20:23:52+03:00",
//			"period": "monthly",
//			"status": "wait_accept",
//			"url": "https://pay.cryptomus.com/recurring/1ee00fbd-42e7-4653-b316-1ab12f8a447b",
//			"last_pay_off": null
//		}
//	}
func (m *Merchant) CreateRecurringInvoice(request RecurringInvoice) (RecurringPayment, error) {
	httpResponse, err := m.sendPaymentRequest("POST", urlCreateRecurringPayment, request)
	if err != nil {
		return RecurringPayment{}, err
	}
	defer httpResponse.Body.Close()

	var response = struct {
		State   int              `json:"state"`
		Result  RecurringPayment `json:"result"`
		Message string           `json:"message"`
		// If some parameter is required and not passed
		Errors struct {
			Amount   []string `json:"amount"`
			Currency []string `json:"currency"`
			Name     []string `json:"name"`
			Period   []string `json:"period"`
		} `json:"errors"`
		Code  int    `json:"code"`
		Error string `json:"error"`
	}{}

	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return RecurringPayment{}, fmt.Errorf("error decoding response: %w", err)
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
	errs = append(errs, response.Errors.Name...)
	errs = append(errs, response.Errors.Period...)

	if httpResponse.StatusCode != http.StatusOK || response.State != 0 || len(errs) > 0 {
		return RecurringPayment{}, fmt.Errorf("error with status %s: %v", httpResponse.Status, strings.Join(errs, "; "))
	}

	return response.Result, nil
}
