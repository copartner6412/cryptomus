package cryptomus

import "time"

// Payment defines the payment information from Cryptomus
//
// See "Creating an invoice" https://doc.cryptomus.com/business/payments/creating-invoice
//
// See "Payment information" https://doc.cryptomus.com/business/payments/payment-information
//
// # Payment example
//
//	{
//		"state": 0,
//		"result": {
//			"uuid": "26109ba0-b05b-4ee0-93d1-fd62c822ce95",
//			"order_id": "1",
//			"amount": "15.00",
//			"payment_amount": null,
//			"payer_amount": null,
//			"discount_percent": null,
//			"discount": "0.00000000",
//			"payer_currency": null,
//			"currency": "USD",
//			"merchant_amount": null,
//			"network": null,
//			"address": null,
//			"from": null,
//			"txid": null,
//			"payment_status": "check",
//			"url": "https://pay.cryptomus.com/pay/26109ba0-b05b-4ee0-93d1-fd62c822ce95",
//			"expired_at": 1689098133,
//			"status": "check",
//			"is_final": false,
//			"additional_data": null,
//			"created_at": "2023-07-11T20:23:52+03:00",
//			"updated_at": "2023-07-11T21:24:17+03:00"
//		}
//	}
//
//	{
//		"state": 0,
//		"result": {
//			"uuid": "70b8db5c-b952-406d-af26-4e1c34c27f15",
//			"order_id": "65bbe87b4098c17a31cff3e71e515243",
//			"amount": "15.00",
//			"payment_amount": "0.00",
//			"payer_amount": "15.75",
//			"discount_percent": -5,
//			"discount": "-0.75",
//			"payer_currency": "USDT",
//			"currency": "USDT",
//			"comments": null,
//			"merchant_amount": "15.43500000",
//			"network": "tron",
//			"address": "TXhfYSWt2oKRrHAJVJeYRuit6ZzKuoEKXj",
//			"from": null,
//			"txid": null,
//			"payment_status": "cancel",
//			"url": "https://pay.cryptomus.com/pay/70b8db5c-b952-406d-af26-4e1c34c27f15",
//			"expired_at": 1689099831,
//			"status": "cancel",
//			"is_final": true,
//			"additional_data": null,
//			"created_at": "2023-07-11T20:23:52+03:00",
//			"updated_at": "2023-07-11T21:24:17+03:00"
//		}
//	}
type Payment struct {
	// Invoice uuid
	UUID string `json:"uuid"`
	// Order ID in your system
	OrderID string `json:"order_id"`
	// The amount of the invoice
	Amount string `json:"amount"`
	// Amount paid by client
	PaymentAmount string `json:"payment_amount"`
	// The amount in payer_currency that the customer must pay, including a discount or additional commission.
	PayerAmount string `json:"payer_amount"`
	// Percentage of discount or additional commission, that was passed in request parameters
	DiscountPercent int `json:"discount_percent"`
	// Actual amount of discount or additional commission in cryptocurrency.
	//
	// For example, if invoice amount is 15 USDT and discount_percent is -5, the discount value will be -0.75
	//
	// i.e. amount + discount = payer_amount
	Discount string `json:"discount"`
	// The currency in which the customer must make the payment.
	PayerCurrency string `json:"payer_currency"`
	// Invoice currency code
	Currency string `json:"currency"`
	// Amount in crypto that will be credited to your balance. If invoice payer_currency is not specified, the value will be null.
	MerchantAmount string `json:"merchant_amount"`
	// Blockchain network code
	Network string `json:"network"`
	// Wallet address for payment
	Address string `json:"address"`
	// The address of the wallet from which the payment was made
	From string `json:"from"`
	// Transaction hash
	TxID string `json:"txid"`
	// Payment status
	//  - paid: The payment was successful and the client paid exactly as much as required.
	//  - paid_over: The payment was successful and client paid more than required.
	//  - wrong_amount: The client paid less than required
	//  - process: Payment in processing
	//  - confirm_check: We have seen the transaction in the blockchain and are waiting for the required number of network confirmations.
	//  - wrong_amount_waiting: The client paid less than required, with the possibility of an additional payment
	//  - check: Waiting for the transaction to appear on the blockchain
	//  - fail: Payment error
	//  - cancel: Payment cancelled, the client did not pay
	//  - system_fail: A system error has occurred
	//  - refund_process: The refund is being processed
	//  - refund_fail: An error occurred during the refund
	//  - refund_paid: The refund was successful
	//  - locked: Funds are locked due to the AML program
	PaymentStatus string `json:"payment_status"`
	// URL payment page
	URL string `json:"url"`
	// Timestamp of expiration of the invoice
	ExpiredAt int64 `json:"expired_at"`
	// Whether the invoice is finalized.
	//
	// When invoice is finalized it is impossible to pay an invoice (it's either paid or expired)
	IsFinal bool `json:"is_final"`
	// Additional information
	AdditionalData string `json:"additional_data"`
	// Creation date of the invoice. Timezone is UTC+3
	CreatedAt time.Time `json:"created_at"`
	// Last invoice updated date. Timezone is UTC+3
	UpdatedAt time.Time `json:"updated_at"`
}
