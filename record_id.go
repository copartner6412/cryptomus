package cryptomus

// Record (invoice/payout) represents parameters required to identify an entity by either UUID or OrderID.
// Only one of UUID or OrderID is required; if both are provided, OrderID will be prioritized for identification.
//
// See "Payment information" https://doc.cryptomus.com/business/payments/payment-information
//
// See "Payout information" https://doc.cryptomus.com/business/payouts/payout-information
//
// See "Resend webhook" https://doc.cryptomus.com/business/payments/resend-webhook
//
// Example
//
//	{
//	    "uuid": "8b03432e-385b-4670-8d06-064591096795"
//	}
type RecordID struct {
	// (Required) Record UUID
	//    uuid
	//    required_without: order_id
	UUID *string `json:"uuid,omitempty"`
	// (Required) Record Order ID
	//    min: 1
	//    max: 128 (max: 32 for GetPayoutInformation, BlockStaticWallet RefundBlockedAddress)
	//    alpha_dash
	//    required_without: uuid
	OrderID *string `json:"order_id,omitempty"`
}
