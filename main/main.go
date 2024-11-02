package main

import (
	"fmt"

	"github.com/copartner6412/cryptomus"
)

func main() {
	trades, err := cryptomus.GetExchangeRate("BTC")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(trades)
}