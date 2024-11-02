package cryptomus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// See "Get directions list" https://doc.cryptomus.com/personal/converts/directions-list
type Direction struct {
	// Currency from
	From string `json:"from"`
	// Currency to
	To string `json:"to"`
	// Min amount from
	MinFrom string `json:"min_from"`
	// Max amount from
	MaxFrom string `json:"max_from"`
	// Min amount to
	MinTo string `json:"min_to"`
	// Max amount to
	MaxTo string `json:"max_to"`
	// Course
	Rate string `json:"rate"`
}

// See "Get directions list" https://doc.cryptomus.com/personal/converts/directions-list
//
// # Response example
//
//	{
//		"state": 0,
//		"result": {
//		  "items": [
//			{
//			  "from": "TRX",
//			  "to": "ETH",
//			  "min_from": "100.00000000",
//			  "min_to": "0.00100000",
//			  "max_from": "100000.00000000",
//			  "max_to": "1000000.00000000",
//			  "rate": "0.00003451"
//			},
//			...
//		  ]
//		}
//	}
func (u *User) ListDirections() ([]Direction, error) {
	httpResponse, err := u.sendPaymentRequest("GET", urlListDirections, nil)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	var response = struct {
		State   int         `json:"state"`
		Result  []Direction `json:"result"`
		Message string      `json:"message"`
		Code    int         `json:"code"`
		Error   string      `json:"error"`
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
