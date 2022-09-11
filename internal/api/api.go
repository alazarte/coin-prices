package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const (
	PRICE_API = "https://api.coincap.io/v2/assets"
	EXCH_API  = "https://blockchain.info/ticker"
)

var (
	supportedCoins = []string{"bitcoin", "monero"}
	supportedFiat  = []string{"usd"}
)

// TODO do I need this struct?
type DataResponse struct {
	Coin CoinResponse `json:"data"`
}

type CoinResponse struct {
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	PriceUsd string `json:"priceUsd"`
	Explorer string `json:"explorer"`
}

func handleGet(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err := res.Body.Close(); err != nil {
		return nil, err
	}

	// ReadAll err
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 ||
		res.StatusCode > 299 {
		err = fmt.Errorf("Status: %s", res.Status)
	}

	return body, err
}

func getCoinPrice(coinStr string) (float64, error) {
	if isFiat(coinStr) {
		return 1, nil
	}
	coin, err := GetPrice(coinStr)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(coin.PriceUsd, 64)
}

func GetExch(coinfromStr string, cointoStr string, amount float64) (float64, error) {
	if !validCoin(coinfromStr) || !validCoin(cointoStr) {
		return 0,
			fmt.Errorf("'from' or 'to' coin not recognized: %s->%s",
				coinfromStr, cointoStr)
	}

	pcoinfrom, err := getCoinPrice(coinfromStr)
	if err != nil {
		return 0, err
	}
	pcointo, err := getCoinPrice(cointoStr)
	if err != nil {
		return 0, err
	}

	return (pcoinfrom * amount) / pcointo, nil

	/*
		return (amount * pcoinfrom) / pcointo, nil

		// perform exchange

		if cointoStr == "usd" {
			return amount / pcoin, nil
		}

		cointo, err := GetPrice(cointoStr)
		if err != nil {
			return 0, err
		}

		pcointo, err := strconv.ParseFloat(cointo.PriceUsd, 64)
		if err != nil {
			return 0, err
		}

		return (amount * pcoinfrom) / pcointo, nil
	*/
}

func isStringInArray(s string, arr []string) bool {
	for _, a := range arr {
		if s == a {
			return true
		}
	}
	return false
}

func validCoin(coin string) bool {
	return isStringInArray(coin, append(supportedCoins, supportedFiat...))
}

func isFiat(coin string) bool {
	return isStringInArray(coin, supportedFiat)
}

func GetPrice(coinStr string) (CoinResponse, error) {
	var data DataResponse

	if !validCoin(coinStr) {
		return data.Coin, fmt.Errorf("Coin not recognized: %s", coinStr)
	}
	if isFiat(coinStr) {
		return data.Coin, fmt.Errorf("Cannot get price of FIAT: %s", coinStr)
	}

	res, err := handleGet(fmt.Sprintf("%s/%s", PRICE_API, coinStr))
	if err != nil {
		return data.Coin, fmt.Errorf("Failed to get coin price: %s (%s)", res, err)
	}

	if err := json.Unmarshal(res, &data); err != nil {
		return data.Coin, err
	}

	return data.Coin, nil
}
