package cryptomus_test

import (
	"testing"

	"github.com/copartner6412/cryptomus"
)

var codes = map[string]struct{ Currency, Network string }{
	"Arbitrum USDT":         {"USDT", "arbitrum"},
	"Arbitrum USDC":         {"USDC", "arbitrum"},
	"Arbitrum Ethereum":     {"ETH", "arbitrum"},
	"Avalanche C-Chain USDT": {"USDT", "avalanche"},
	"Avalanche C-Chain AVAX": {"AVAX", "avalanche"},
	"Avalanche C-Chain USDC": {"USDC", "avalanche"},
	"Bitcoin Cash":          {"BCH", "bch"},
	"ChainGPT":              {"CGPT", "bsc"},
	"Binance USDT":          {"USDT", "bsc"},
	"Binance USDC":          {"USDC", "bsc"},
	"Binance Ethereum":      {"ETH", "bsc"},
	"Binance DAI":           {"DAI", "bsc"},
	"Binance BNB":           {"BNB", "bsc"},
	"Bitcoin":               {"BTC", "btc"},
	"Dash":                  {"DASH", "dash"},
	"Dogecoin":              {"DOGE", "doge"},
	"VERSE Token":           {"VERSE", "eth"},
	"Ethereum MATIC":        {"MATIC", "eth"},
	"Ethereum USDC":         {"USDC", "eth"},
	"Ethereum":              {"ETH", "eth"},
	"Ethereum USDT":         {"USDT", "eth"},
	"Ethereum DAI":          {"DAI", "eth"},
	"Litecoin":              {"LTC", "ltc"},
	"CRMS":                  {"CRMS", "polygon"},
	"Polygon USDC":          {"USDC", "polygon"},
	"Polygon MATIC":         {"MATIC", "polygon"},
	"Polygon DAI":           {"DAI", "polygon"},
	"Polygon USDT":          {"USDT", "polygon"},
	"Solana USDT":           {"USDT", "sol"},
	"Solana":                {"SOL", "sol"},
	"TON USDT":              {"USDT", "ton"},
	"TON":                   {"TON", "ton"},
	"Tron TRX":              {"TRX", "tron"},
	"Tron USDC":             {"USDC", "tron"},
	"Tron USDT":             {"USDT", "tron"},
	"Monero":                {"XMR", "xmr"},
}

var currencies = []string{
	"USDT", "USDC", "ETH", "AVAX", "BCH", "CGPT", "DAI", "BNB", "BTC", "DASH",
	"DOGE", "VERSE", "MATIC", "LTC", "CRMS", "SOL", "TON", "TRX", "XMR",
}

var networks = []string{
	"arbitrum", "avalanche", "bch", "bsc", "btc", "dash", "doge", "eth", 
	"ltc", "polygon", "sol", "ton", "tron", "xmr",
}

func TestGetAssets(t *testing.T) {
	assets, err := cryptomus.GetAssets()
	t.Log(assets)
	if err != nil {
		t.Errorf("error getting assets: %v", err)
	}
}

func TestGetExchangeRate(t *testing.T) {
	for _, currency := range currencies {
		rates, err := cryptomus.GetExchangeRate(currency)
		t.Log(rates)
		if err != nil {
			t.Errorf("error getting exchange rate for %s: %v", currency, err)
		}
	}
}

func TestGetOrderBook(t *testing.T) {
	t.Parallel()
	timestamp, bids, asks, err := cryptomus.GetOrderBook("BTC_USDT", 1)
	t.Log(timestamp)
	t.Log(bids)
	t.Log(asks)
	if err != nil {
		t.Errorf("error getting order book for currency pair %s and level %d: %v", "BTC_USDT", 1, err)
	}
}

func TestGetTrades(t *testing.T) {
	t.Parallel()
	currencyPair := "BTC_USDT"
	trades, err := cryptomus.GetTrades(currencyPair)
	t.Log(trades)
	if err != nil {
		t.Errorf("error getting trades for currency pair %s: %v", currencyPair, err)
	}
}
