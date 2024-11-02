package cryptomus

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
)

// Webhook is a kind of feedback method for payment information.
// When invoice/withdrawal status changes, a POST request is sent to the url_callback specified when creating the invoice/withdrawal.

// See "Webhook" https://doc.cryptomus.com/business/payments/webhook
//
// See "Webhook" https://doc.cryptomus.com/business/payouts/webhook
//
// # Example
//
//	{
//	  "type": "payment",
//	  "uuid": "62f88b36-a9d5-4fa6-aa26-e040c3dbf26d",
//	  "order_id": "97a75bf8eda5cca41ba9d2e104840fcd",
//	  "amount": "3.00000000",
//	  "payment_amount": "3.00000000",
//	  "payment_amount_usd": "0.23",
//	  "merchant_amount": "2.94000000",
//	  "commission": "0.06000000",
//	  "is_final": true,
//	  "status": "paid",
//	  "from": "THgEWubVc8tPKXLJ4VZ5zbiiAK7AgqSeGH",
//	  "wallet_address_uuid": null,
//	  "network": "tron",
//	  "currency": "TRX",
//	  "payer_currency": "TRX",
//	  "additional_data": null,
//	  "convert": {
//	    "to_currency": "USDT",
//	    "commission": null,
//	    "rate": "0.07700000",
//	    "amount": "0.22638000"
//	  },
//	  "txid": "6f0d9c8374db57cac0d806251473de754f361c83a03cd805f74aa9da3193486b",
//	  "sign": "a76c0d77f3e8e1a419b138af04ab600a"
//	}
//
//	{
//		"type": "payout",
//		"uuid": "2b852d86-3cf1-43fb-b1bb-36f0b7d12151",
//		"order_id": "129359",
//		"amount": "207.00000000",
//		"merchant_amount": "207.30000000",
//		"commission": "0.30000000",
//		"is_final": true,
//		"status": "paid",
//		"txid": "0xcf8",
//		"currency": "USDT",
//		"network": "bsc",
//		"payer_currency": "USDT",
//		"payer_amount": "207.00000000",
//		"sign": "eff3afba8600af59c98b74155934da2d"
//	}
type Update struct {
	// (Common) Available options:
	//  - wallet (invoice)
	//  - payment (invoice)
	//  - payout (withdrawal)
	Type *string `json:"type"`
	// (Common) uuid of the payment/payout
	UUID *string `json:"uuid"`
	// (Common) Order ID in your system to identify the payment/payout order
	OrderID *string `json:"order_id"`
	// (Common) The amount of the payment/payout
	Amount *string `json:"amount"`
	// (Only in Payment) Amount actually paid by client
	PaymentAmount *string `json:"payment_amount"`
	// (Only in Payment) Amount actually paid by client in USD
	PaymentAmountUSD *string `json:"payment_amount_usd"`
	// (Common) The amount added to the merchant's balance, with all commissions subtracted or debited from merchant's balance, including commissions.
	MerchantAmount *string `json:"merchant_amount"`
	// (Common) Cryptomus commission amount
	Commission *string `json:"commission"`
	// (Common) Whether the invoice/withdrawal is finalized.
	IsFinal *bool `json:"is_final"`
	// Payment/payout status
	//
	// Available options for payment (https://doc.cryptomus.com/business/payments/payment-statuses):
	//  - confirm_check: We have seen the transaction in the blockchain and are waiting for the required number of network confirmations.
	//  - paid: The payment was successful and the client paid exactly as much as required.
	//  - paid_over: The payment was successful and client paid more than required.
	//  - fail: Payment error
	//  - wrong_amount: The client paid less than required
	//  - cancel: Payment cancelled, the client did not pay
	//  - system_fail: A system error has occurred
	//  - refund_process: The refund is being processed
	//  - refund_fail: An error occurred during the refund
	//  - refund_paid: The refund was successful
	// Available options for payout (https://doc.cryptomus.com/business/payouts/payout-statuses):
	//  - process: Payout in process
	//  - check: The payout is being verified
	//  - paid: The payout was successful
	//  - fail: Payout failed
	//  - cancel: Payout cancelled
	//  - system_fail: A system error has occurred
	Status *string `json:"status"`
	// (Only in Payment) Payer's wallet address
	From *string `json:"from"`
	// (Only in Payment) uuid of the static wallet
	WalletAddressUUID *string `json:"wallet_address_uuid"`
	// (Common) The blockchain network for payment/payout
	Network *string `json:"network"`
	// (Common) The currency of payment/payout
	Currency *string `json:"currency"`
	// (Common) The cryptocurrency code in which payment/payout will be actually made.
	PayerCurrency *string `json:"payer_currency"`
	// (Only in Payout) Amount in payer_currency of the payout
	PayerAmount *string `json:"payer_amount"`
	// (Only in Payment) Additional information string that you provided when creating an invoice
	AdditionalData *string `json:"additional_data"`
	// (Only in Payment) Information about the currency to which the payment will be automatically converted. Conversion is performed from payer_currency to USDT
	//
	// The convert field will not exist if you have not enabled the automatic conversion function for payer_currency (e.g. auto convert BTC to USDT)
	Convert *AutomaticConvert `json:"convert"`
	// (Common) Transaction hash on the blockchain.
	//
	// The txid for payment field will not exist if:
	//  1. payment was paid by p2p (The payer withdrew funds from his Cryptomus account to the address indicated in the invoice and the payment was made without blockchain, only in our system)
	//  2. Payment was not paid
	//  3. Something was wrong with the payment or the client made a mistake and we marked it as ‘paid’ manually
	// The txid for payout field will not exist if:
	//  1. Payout was paid by p2p to Cryptomus address
	//  2. Payout failed.
	TxID *string `json:"txid"`
	// (Common) Signature
	Sign string `json:"sign"`
}

// See "Webhook" https://doc.cryptomus.com/business/payments/webhook
type AutomaticConvert struct {
	// The currency code to which the payment will be converted
	ToCurrency *string `json:"to_currency"`
	// Conversion fee
	Commission *string `json:"commission"`
	// Conversion rate
	Rate *string `json:"rate"`
	// Conversion amount in to_currency that was added to the merchant's balance, with all commissions subtracted.
	//
	// amount = merchant_amount * rate
	Amount *string `json:"amount"`
}

// Your api keys are secret and no one except you and cryptomus should know them. So, when verifying the signature, you will be sure that the webhook was sent by cryptomus.
//
// We create a sign using this algorithm. MD5 hash of the body of the POST request encoded in base64 and combined with your API key.
//
// As the signature comes in the body of the request, to verify it, you need to extract the sign from the response body, generate a hash from the body and your API KEY and match it with the sign parameter.
//
// See "Webhook" https://doc.cryptomus.com/business/payments/webhook
func (m *Merchant) VerifySign(update Update) error {
	var sign string
	switch *update.Type {
	case "payment", "wallet":
		paymentUpdateWithoutSignature := struct {
			Type              *string           `json:"type"`
			UUID              *string           `json:"uuid"`
			OrderID           *string           `json:"order_id"`
			Amount            *string           `json:"amount"`
			PaymentAmount     *string           `json:"payment_amount"`
			PaymentAmountUSD  *string           `json:"payment_amount_usd"`
			MerchantAmount    *string           `json:"merchant_amount"`
			Commission        *string           `json:"commission"`
			IsFinal           *bool             `json:"is_final"`
			Status            *string           `json:"status"`
			From              *string           `json:"from"`
			WalletAddressUUID *string           `json:"wallet_address_uuid"`
			Network           *string           `json:"network"`
			Currency          *string           `json:"currency"`
			PayerCurrency     *string           `json:"payer_currency"`
			AdditionalData    *string           `json:"additional_data"`
			Convert           *AutomaticConvert `json:"convert"`
			TxID              *string           `json:"txid"`
		}{
			Type:              update.Type,
			UUID:              update.UUID,
			OrderID:           update.OrderID,
			Amount:            update.Amount,
			PaymentAmount:     update.PaymentAmount,
			PaymentAmountUSD:  update.PaymentAmountUSD,
			MerchantAmount:    update.MerchantAmount,
			Commission:        update.Commission,
			IsFinal:           update.IsFinal,
			Status:            update.Status,
			From:              update.From,
			WalletAddressUUID: update.WalletAddressUUID,
			Network:           update.Network,
			Currency:          update.Currency,
			PayerCurrency:     update.PayerCurrency,
			AdditionalData:    update.AdditionalData,
			Convert:           update.Convert,
			TxID:              update.TxID,
		}

		jsonData, err := json.Marshal(paymentUpdateWithoutSignature)
		if err != nil {
			return fmt.Errorf("error marshalling update payload: %w", err)
		}

		sign, err = m.signPaymentPayload(jsonData)
		if err != nil {
			return fmt.Errorf("error generating payment signature: %w", err)
		}

		if subtle.ConstantTimeCompare([]byte(sign), []byte(update.Sign)) == 0 {
			return fmt.Errorf("signature mismatch")
		}
	case "payout":
		payoutUpdateWithoutSignature := struct {
			Type           *string `json:"type"`
			UUID           *string `json:"uuid"`
			OrderID        *string `json:"order_id"`
			Amount         *string `json:"amount"`
			MerchantAmount *string `json:"merchant_amount"`
			Commission     *string `json:"commission"`
			IsFinal        *bool   `json:"is_final"`
			Status         *string `json:"status"`
			TxID           *string `json:"txid"`
			Currency       *string `json:"currency"`
			Network        *string `json:"network"`
			PayerCurrency  *string `json:"payer_currency"`
			PayerAmount    *string `json:"payer_amount"`
		}{
			Type:           update.Type,
			UUID:           update.UUID,
			OrderID:        update.OrderID,
			Amount:         update.Amount,
			MerchantAmount: update.MerchantAmount,
			Commission:     update.Commission,
			IsFinal:        update.IsFinal,
			Status:         update.Status,
			TxID:           update.TxID,
			Currency:       update.Currency,
			Network:        update.Network,
			PayerCurrency:  update.PayerCurrency,
			PayerAmount:    update.PayerAmount,
		}

		jsonData, err := json.Marshal(payoutUpdateWithoutSignature)
		if err != nil {
			return fmt.Errorf("error marshalling update payload: %w", err)
		}

		sign, err = m.signPayoutPayload(jsonData)
		if err != nil {
			return fmt.Errorf("error generating payment signature: %w", err)
		}
	default:
		return fmt.Errorf("unsupported type: %T", update)
	}

	if subtle.ConstantTimeCompare([]byte(sign), []byte(update.Sign)) == 0 {
		return fmt.Errorf("signature mismatch")
	}

	return nil
}
