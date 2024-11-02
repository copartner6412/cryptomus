package cryptomus

// Invoice defines the payload for creating an invoice
//
// The invoice will have a specific cryptocurrency and address at the time of creation only if currency or to_currency parameter is a cryptocurrency and the network parameter is passed (or a cryptocurrency has only one network, for example BTC).
//
// Required fields:
//   - Amount
//   - Currency
//   - OrderID
//
// See "Creating an invoice" https://doc.cryptomus.com/business/payments/creating-invoice
//
// See "Payment history" https://doc.cryptomus.com/business/payments/payment-history
//
// # Request example
//
// Below is a sample example with the minimum required number of parameters. In this case, an invoice will be created with the fiat currency USD. On the payment page, the user will be able to select the cryptocurrency and network to pay this bill.
//
//	{
//	    "amount": "15",
//		"currency": "USD",
//		"order_id": "1",
//	}
//
// # More examples
//
// Request data to create an invoice for 20 USDT in tron network. The invoice will have an address at the time of creation.
//
//	{
//		"amount": "20",
//		"currency": "USDT",
//		"order_id": "1",
//		"network": "tron"
//	}
//
// Request data to create an invoice for 25 USD and allow clients to pay only by USDT in any network.
//
//	{
//		"amount": "25",
//		"currency": "USD",
//		"order_id": "1",
//		"to_currency": "USDT"
//	}
//
// Request data to create an invoice for 25 USD and allow clients to pay only by USDT in tron network. The invoice will have an address at the time of creation.
//
//	{
//		"amount": "25",
//		"currency": "USD",
//		"order_id": "1",
//		"to_currency": "USDT",
//		"network": "tron"
//	}
//
// Request data to create an invoice for 20 USDT and allow clients to pay in all available networks
//
//	{
//		"amount": "20",
//		"currency": "USDT",
//		"order_id": "1",
//	}
//
// Request data to create an invoice for 20 USD and allow clients to pay in bitcoin only. The invoice will have an address at the time of creation.
//
//	{
//		"amount": "20",
//		"currency": "USD",
//		"order_id": "1",
//		"to_currency": "BTC"
//	}
//
// Request data to create an invoice for 0.5 BTC and allow clients to pay in bitcoin only. The invoice will have an address at the time of creation.
//
//	{
//		"amount": "0.5",
//		"currency": "BTC",
//		"order_id": "1",
//	}
type Invoice struct {
	// (Required) Amount to be paid.
	//
	// If there are pennies in the amount, then send them with a separator '.'
	//
	// Example:
	//    "10.28"
	Amount string `json:"amount"`
	// (Required) Currency code
	Currency string `json:"currency"`
	// (Required) Order ID in your system
	//
	// The parameter should be a string consisting of alphabetic characters, numbers, underscores, and dashes. It should not contain any spaces or special characters.
	//
	// The order_id must be unique within the merchant invoices/static wallets/recurrence payments
	//
	// When we find an existing invoice with order_id, we return its details, a new invoice will not be created.
	//
	//    min: 1
	//    max: 128
	//    alpha_dash
	OrderID string `json:"order_id"`
	// (Optional) Blockchain network code
	//    default: null
	Network *string `json:"network,omitempty"`
	// (Optional) Before paying, the user can click on the button on the payment form and return to the store page at this URL.
	//    min: 6
	//    max: 255
	//    url
	//    default: null
	URLReturn *string `json:"url_return,omitempty"`
	// (Optional) After successful payment, the user can click on the button on the payment form and return to this URL.
	//    min: 6
	//    max: 255
	//    url
	//    default: null
	URLSuccess *string `json:"url_success,omitempty"`
	// (Optional) Url to which webhooks with payment status will be sent
	//    min: 6
	//    max: 255
	//    url
	//    default: null
	URLCallback *string `json:"url_callback,omitempty"`
	// (Optional) Whether the user is allowed to pay the remaining amount. This is useful when the user has not paid the entire amount of the invoice for one transaction, and you want to allow him to pay up to the full amount. If you disable this feature, the invoice will finalize after receiving the first payment and you will receive funds to your balance.
	IsPaymentMultiple *bool `json:"is_payment_multiple,omitempty"`
	// (Optional) The lifespan of the issued invoice (in seconds)
	//    min: 300
	//    max: 43200
	//    default: true
	Lifetime *int `json:"lifetime,omitempty"`
	// (Optional) The parameter is used to specify the target currency for converting the invoice amount. When creating an invoice, you provide an amount and currency, and the API will convert that amount to the equivalent value in the to_currency.
	//
	// For example, to create an invoice for 20 USD in bitcoin:
	//    amount: 20
	//    currency: USD
	//    to_currency: BTC
	//    default: 3600
	// The API will convert 20 USD amount to its equivalent in BTC based on the current exchange rate and the user will pay in BTC
	//
	// The to_currency should always be the cryptocurrency code, not a fiat currency code.
	ToCurrency *string `json:"to_currency,omitempty"`
	// (Optional) Percentage of the payment commission charged to the client
	//
	// If you have a rate of 1%, then if you create an invoice for 100 USDT with subtract = 100 (the client pays 100% commission), the client will have to pay 101 USDT.
	//    min: 0
	//    max: 100
	//    default: 0
	Subtract *int `json:"subtract,omitempty"`
	// (Optional) Acceptable inaccuracy in payment.
	//
	// For example, if you pass the value 5, the invoice will be marked as Paid even if the client has paid only 95% of the amount.
	//
	// The actual payment amount will be credited to the balance
	//    numeric
	//    min: 0
	//    max: 5
	//    default: 0
	AccuracyPaymentPercent *int `json:"accuracy_payment_percent,omitempty"`
	// (Optional) Additional information for you (not shown to the client)
	//    max: 255
	//    default: null
	AdditionalData *string `json:"additional_data,omitempty"`
	// (Optional) List of allowed currencies for payment. This is useful if you want to limit the list of coins that your customers can use to pay invoices.
	//
	// The array of available currencies from your merchant's settings is used (if the parameter is not set, all supported cryptocurrencies are used by default).
	//
	// The available currencies are set in the merchant settings section.
	Currencies []Currency `json:"currencies,omitempty"`
	// (Optional) List of excluded currencies for payment
	//    default: null
	ExceptCurrencies []Currency `json:"except_currencies,omitempty"`
	// (Optional) The service from which the exchange rates are taken for conversion in the invoice.
	//    min: 4
	//    max: 20
	// Available values:
	//  - Binance
	//  - BinanceP2P
	//  - Exmo
	//  - Kucoin
	//  - Garantexio
	// If not passed, Cryptomus exchange rates are used.
	CourseSource *string `json:"course_source,omitempty"`
	// (Optional) The merchant who makes the request connects to a referrer by code.
	//
	// For example, you are an application that generates invoices via the Cryptomus API and your customers are other stores.
	//
	// They enter their api key and merchant id in your application, and you send requests with their credentials and passing your referral code.
	//
	// Thus, your clients become referrals on your Cryptomus account and you will receive income from their turnover.
	//    default: null
	FromReferralCode *string `json:"from_referral_code,omitempty"`
	// (Optional) Positive numbers:
	//
	// Allows you to set a discount.
	//
	// To set a 5% discount for the payment, you should pass a value: 5
	//
	// Negative numbers:
	//
	// Allows you to set custom additional commission.
	//
	// To set an additional commission of 10% for the payment, you should pass a value: -10
	//
	// The discount percentage when creating an invoice is taken into account only if the invoice has a specific cryptocurrency.
	//    min: -99
	//    max: 100
	//    default: null
	DiscountPercent *int `json:"discount_percent,omitempty"`
	// (Optional) Using this parameter, you can update the lifetime and get a new address for the invoice if the lifetime has expired.
	//
	// To do that, you need to pass all required parameters, and the invoice with passed order_id will be refreshed.
	//
	// Only address, payment_status and expired_at are changed. No other fields are changed, regardless of the parameters passed.
	//    default: false
	IsRefresh *bool `json:"is_refresh,omitempty"`
}

type Currency struct {
	// (Required) Currency code
	Currency string `json:"currency"`
	// (Optional) Blockchain network code
	Network *string `json:"network"`
}
