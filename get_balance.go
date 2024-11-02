package cryptomus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// See "MerchantWallet" https://doc.cryptomus.com/business/balance
type MerchantWallet struct {
	// Wallet UUID
	UUID string `json:"uuid"`
	// Business/personal wallet balance
	Balance string `json:"balance"`
	// Wallet currency_code
	CurrencyCode string `json:"currency_code"`
}

// See "Get balance" https://doc.cryptomus.com/personal/converts/balance
type UserWallet struct {
	// Wallet UUID
	WalletUUID string `json:"walletUuid"`
	// Wallet currency
	CurrencyCode string `json:"currency_code"`
	// Wallet balance
	Balance string `json:"balance"`
	// Wallet balance in USD
	BalanceUSD string `json:"balanceUsd"`
}

// See "Balance" https://doc.cryptomus.com/business/balance
//
//	{
//	    "state": 0,
//	    "result": [
//	        {
//	            "balance": {
//	                "merchant": [
//	                    {
//	                        "uuid": "abcdabcd-abcd-1234-1234-abcdabcd",
//	                        "balance": "0.00000000",
//	                        "currency_code": "ETH"
//	                    },
//	                    {
//	                        "uuid": "abcdabcd-abcd-1234-1234-abcdabcd",
//	                        "balance": "0.57000000",
//	                        "currency_code": "BTC"
//	                    },
//	                    {
//	                        "uuid": "abcdabcd-abcd-1234-1234-abcdabcd",
//	                        "balance": "23.57327446",
//	                        "currency_code": "TRX"
//	                    },
//	                    {
//	                        "uuid": "abcdabcd-abcd-1234-1234-abcdabcd",
//	                        "balance": "5.00000000",
//	                        "currency_code": "USDT"
//	                    },
//	                    {
//	                        "uuid": "abcdabcd-abcd-1234-1234-abcdabcd",
//	                        "balance": "10.00120000",
//	                        "currency_code": "DASH"
//	                    },
//	                    {
//	                        "uuid": "abcdabcd-abcd-1234-1234-abcdabcd",
//	                        "balance": "0.18500000",
//	                        "currency_code": "LTC"
//	                    }
//	                ],
//	                "user": [
//	                    {
//	                        "uuid": "abcdabcd-abcd-1234-1234-abcdabcd",
//	                        "balance": "0.40000000",
//	                        "currency_code": "BTC"
//	                    },
//	                    {
//	                        "uuid": "abcdabcd-abcd-1234-1234-abcdabcd",
//	                        "balance": "52.00000000",
//	                        "currency_code": "USDT"
//	                    },
//	                    {
//	                        "uuid": "abcdabcd-abcd-1234-1234-abcdabcd",
//	                        "balance": "0.00000000",
//	                        "currency_code": "DASH"
//	                    },
//	                    {
//	                        "uuid": "abcdabcd-abcd-1234-1234-abcdabcd",
//	                        "balance": "0.30000000",
//	                        "currency_code": "LTC"
//	                    },
//	                    {
//	                        "uuid": "abcdabcd-abcd-1234-1234-abcdabcd",
//	                        "balance": "27.00000000",
//	                        "currency_code": "TRX"
//	                    },
//	                    {
//	                        "uuid": "abcdabcd-abcd-1234-1234-abcdabcd",
//	                        "balance": "0.19000000",
//	                        "currency_code": "ETH"
//	                    }
//	                ]
//	            }
//	        }
//	    ]
//	}
func (m *Merchant) GetBalance() (merchantBalances, userBalances []MerchantWallet, err error) {
	httpResponse, err := m.sendPaymentRequest("POST", urlGetBalanceForMerchant, nil)
	if err != nil {
		return nil, nil, err
	}
	defer httpResponse.Body.Close()

	var response = struct {
		State  int `json:"state"`
		Result []struct {
			Balance struct {
				Merchant []MerchantWallet `json:"merchant"`
				User     []MerchantWallet `json:"user"`
			} `json:"balance"`
		} `json:"result"`
		Message string `json:"message"`
		Code    int    `json:"code"`
		Error   string `json:"error"`
	}{}

	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, nil, fmt.Errorf("error decoding response: %w", err)
	}

	var errs []string
	if response.Message != "" {
		errs = append(errs, response.Message)
	}
	if response.Error != "" {
		errs = append(errs, response.Error)
	}

	if httpResponse.StatusCode != http.StatusOK || response.State != 0 || len(errs) > 0 {
		return nil, nil, fmt.Errorf("error with status %s: %v", httpResponse.Status, strings.Join(errs, "; "))
	}

	return response.Result[0].Balance.Merchant, response.Result[0].Balance.User, nil
}

// See "Get balance" https://doc.cryptomus.com/personal/converts/balance
//
// # Response example
//
//	{
//		"state": 0,
//		"result": {
//		  "balances": [
//			{
//			  "walletUuid": "4ba23a47-a182-4d87-8c68-247c974be566",
//			  "currency_code": "BCH",
//			  "balance": "0.00000000",
//			  "balanceUsd": "0.00"
//			},
//			...
//			{
//			  "walletUuid": "539f051f-7ceb-4ac4-831e-21ebdba0a5d0",
//			  "currency_code": "XMR",
//			  "balance": "0.00000000",
//			  "balanceUsd": "0.00"
//			},
//		  ]
//		}
//	  }
func (u *User) GetBalance() ([]UserWallet, error) {
	httpResponse, err := u.sendPaymentRequest("GET", urlGetBalanceForUser, nil)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	var response = struct {
		State   int          `json:"state"`
		Result  []UserWallet `json:"result"`
		Message string       `json:"message"`
		Code    int          `json:"code"`
		Error   string       `json:"error"`
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

	return response.Result, nil
}
