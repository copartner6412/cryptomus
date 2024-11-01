package cryptomus

// See "Creating recurring payment" https://doc.cryptomus.com/business/recurring/creating
//
// See "Payment information" https://doc.cryptomus.com/business/recurring/info
//
// See "List of recurring payments" https://doc.cryptomus.com/business/recurring/list
//
// # Response example
//
// If the request is successful. the response code is 200.
//
//	{
//		"state": 0,
//		"result": {
//			"uuid": "afd050e8-35ea-4129-bbdd-73f510dce556",
//			"name": "Recurring payment",
//			"order_id": null,
//			"amount": "15",
//			"currency": "USDT",
//			"payer_currency": "USDT",
//			"payer_amount_usd": "15.00",
//			"payer_amount": "15.00000000",
//			"url_callback": null,
//			"period": "monthly",
//			"status": "wait_accept",
//			"url": "https://pay.cryptomus.com/recurring/afd050e8-35ea-4129-bbdd-73f510dce556",
//			"last_pay_off": null
//		}
//	}
//
//	{
//		"state": 0,
//		"result": {
//		  "uuid": "bbe5ce96-1126-4843-a0d2-b432e77669c2",
//		  "name": "Access to personal account",
//		  "order_id": "1487555",
//		  "amount": "5",
//		  "currency": "USD",
//		  "payer_currency": "USDT",
//		  "payer_amount_usd": "5.00",
//		  "payer_amount": "5.00",
//		  "url_callback": null,
//		  "discount_days": "30",
//		  "discount_amount": "50.00",
//		  "end_of_discount": null,
//		  "period": "weekly",
//		  "status": "wait_accept",
//		  "url": "https://pay.cryptomus.com/pay/bbe5ce96-1126-4843-a0d2-b432e77669c2",
//		  "last_pay_off": null
//		}
//	  }
type RecurringPayment struct {
	// Recurring uuid
	UUID string `json:"uuid"`
	// Recurring name
	Name string `json:"name"`
	// Order ID in your system
	OrderID *string `json:"order_id"`
	// Recurring amount
	Amount string `json:"amount"`
	// Currency code
	Currency string `json:"currency"`
	// The currency in which the customer must make the payment.
	PayerCurrency string `json:"payer_currency"`
	// The equivalent USD amount in payer_currency that the customer must pay
	PayerAmountUSD string `json:"payer_amount_usd"`
	// The amount in payer_currency that the customer must pay
	PayerAmount string `json:"payer_amount"`
	// Url to which webhooks with payment status will be sent
	URLCallback *string `json:"url_callback"`
	// Recurring payment period
	Period string `json:"period"`
	// Recurring status
	//
	// Available options:
	//  - wait_accept
	//  - cancel_by_merchant
	//  - active
	//  - cancel_by_user
	Status string `json:"status"`
	// The URL of the Cryptomus payment page where the payer will make the payment
	URL string `json:"url"`
	// Date of the last payment. The time zone is UTC+3. If the value is null, no payments were made.
	LastPayOff *string `json:"last_pay_off"`
	// Additional recurring payment details
	AdditionalData *string `json:"additional_data"`
}
