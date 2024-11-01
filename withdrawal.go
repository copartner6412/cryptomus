package cryptomus

// Withdrawal holds the required and optional fields for a payout request.
//
// See "Creating a payout" https://doc.cryptomus.com/business/payouts/creating-payout
//
// # Request example
//
//	{
//	    "amount": "5",
//	    "currency": "USDT",
//	    "network": "TRON",
//	    "order_id": "1",
//	    "address": "TDD97yguPESTpcrJMqU6h2ozZbibv4Vaqm",
//	    "url_callback": "https://your.site/callback",
//	    "is_subtract": "1",
//	}
type Withdrawal struct {
	// (Required) Payout amount
	Amount string `json:"amount"`
	// (Required) Currency code for the payout
	//
	// If Currency is fiat, the to_currency parameter is required.
	Currency string `json:"currency"`
	// (Required) Order ID in your system
	//
	// The parameter should be a string consisting of alphabetic characters, numbers, underscores, and dashes. It should not contain any spaces or special characters.
	//
	// The order_id must be unique within the merchant payouts
	//
	// When we find an existing payout with order_id, we return its details, a new payout will not be created.
	//    min: 1
	//    max: 100
	//    alpha_dash
	OrderID string `json:"order_id"`
	// (Required) The address of the wallet to which the withdrawal will be made
	Address string `json:"address"`
	// (Required) Defines where the withdrawal fee will be deducted
	//
	// true - from your balance
	//
	// false - from payout amount, the payout amount will be decreased
	IsSubtract *bool `json:"is_subtract,omitempty"`
	// (Required) Blockchain network code
	//
	// Not required when the currency/to_currency is a cryptocurrency and has only one network, for example BTC
	Network *string `json:"network,omitempty"`
	// (Optional) URL to which webhooks with payout status will be sent
	URLCallback *string `json:"url_callback,omitempty"`
	// (Optional) Cryptocurrency code in which the payout will be made. It is used when the currency parameter is fiat.
	ToCurrency *string `json:"to_currency,omitempty"`
	// (Optional) The service from which the exchange rates are taken for conversion in the invoice.
	//
	// The parameter is applied only if the currency is fiat, otherwise the default value is taken from the merchant's settings.
	//    default: Value from merchant's settings
	// Available values:
	//  - Binance
	//  - BinanceP2p
	//  - Exmo
	//  - Kucoin
	//  - Garantexio
	CourseSource *string `json:"course_source,omitempty"`
	// (Optional) Allows to automatically convert the withdrawal amount and use the from_currency balance. Only USDT is available.
	//    default: null
	FromCurrency *string `json:"from_currency,omitempty"`
	// (Optional) The parameter for selecting the withdrawal priority. The cost of the withdrawal fee depends on the selected parameter.
	//
	// This parameter is applied only in case of using the BTC, ETH, POLYGON, and BSC networks.
	//    min: 4
	//    max: 11
	//    default: recommended
	// Available values:
	//  - recommended
	//  - economy
	//  - high
	//  - highest
	Priority *string `json:"priority,omitempty"`
	// (Optional) Additional identifier for TON, used to specify a particular recipient or target
	//    min: 1
	//    max: 30
	Memo *string `json:"memo,omitempty"`
}
