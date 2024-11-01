package cryptomus

// Payout holds the response structure for a payout transaction request.
//
// See "Creating a payout" https://doc.cryptomus.com/business/payouts/creating-payout
//
// See "Payout history" https://doc.cryptomus.com/business/payouts/payout-history
//
//	{
//	    "state": 0,
//	    "result": {
//	        "uuid": "a7c0caec-a594-4aaa-b1c4-77d511857594",
//	        "amount": "3",
//	        "currency": "USDT",
//	        "network": "TRON",
//	        "address": "TJ...",
//	        "txid": null,
//	        "status": "process",
//	        "is_final": false,
//	        "balance": 129,
//	        "payer_currency": "USD",
//	        "payer_amount": 3
//	    }
//	}
type Payout struct {
	// uuid of the payout
	UUID string `json:"uuid"`
	// Payout amount in currency
	Amount string `json:"amount"`
	// Currency code for the payout
	Currency string `json:"currency"`
	// The code of the blockchain network in which the payment is made
	Network string `json:"network"`
	// The address of the wallet to which the payment is made
	Address string `json:"address"`
	// Transaction ID in the blockchain
	TxID *string `json:"txid"`
	// Payout status
	//
	// The payout status comes in the body of the response of some methods and indicates at what stage the payout is at the moment
	//
	// Available options:
	//  - process: Payout in process
	//  - check: The payout is being verified
	//  - paid: The payout was successful
	//  - fail: Payout failed
	//  - cancel: Payout cancelled
	//  - system_fail: A system error has occurred
	Status string `json:"status"`
	// Whether the payout is finalized
	//
	// The payout process is considered finalized once it has been successfully paid or if it has failed. In the event of a payout failure, the funds will be returned to your balance, requiring you to initiate the payout process again.
	IsFinal bool `json:"is_final"`
	// The remaining funds on the merchant's balance.
	Balance float64 `json:"balance"`
	// Cryptocurrency code in which the payout will be actually made. The payout address will receive the payout currency. (only in CreatePayout)
	PayerCurrency string `json:"payer_currency"`
	// Amount in payer_currency of the payout. (only in CreatePayout)
	PayerAmount float64 `json:"payer_amount"`
	// Creation date of the payout. Timezone is UTC+3 (only in ListPayoutHistory)
	CreatedAt string `json:"created_at"`
	// Last payout updated date. Timezone is UTC+3 (only in ListPayoutHistory)
	UpdatedAt string `json:"updated_at"`
}
