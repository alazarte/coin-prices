package api

import (
        "fmt"
        "encoding/json"
	"strconv"
)

func (c CoinConfig) validCoin(coin string) bool {
	return isStringInArray(coin, append(c.SupportedCoins, c.SupportedFiat...))
}

func (c CoinConfig) isFiat(coin string) bool {
	return isStringInArray(coin, c.SupportedFiat)
}

func (c CoinConfig) getCoinPrice(coinStr string) (float64, error) {
	var data DataResponse

	if c.isFiat(coinStr) {
		return 1, nil
	}

	res, err := handleGet(fmt.Sprintf("%s/%s", c.PriceApi, coinStr))
	if err != nil {
		return 0, fmt.Errorf("Failed to get coin price: %s (%s)", res, err)
	}

	if err := json.Unmarshal(res, &data); err != nil {
		return 0, err
	}


	return strconv.ParseFloat(data.Coin.PriceUsd, 64)
}

