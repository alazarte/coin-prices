package api

import (
	"encoding/json"
	"fmt"
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

type CoinConfig struct {
	PriceApi       string   `json:"priceApi"`
	ExchApi        string   `json:"exchApi"`
	SupportedCoins []string `json:"supportedCoins"`
	SupportedFiat  []string `json:"supportedFiat"`
}

func FromJSON(data []byte) (CoinConfig, error) {
	var config CoinConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return CoinConfig{}, err
	}

	// TODO validate config
	return config, nil
}

func (c CoinConfig) GetExch(coinfromStr string, cointoStr string, amount float64) (float64, error) {
	if !c.validCoin(coinfromStr) || !c.validCoin(cointoStr) {
		return 0,
			fmt.Errorf("'from' or 'to' coin not recognized: %s->%s",
				coinfromStr, cointoStr)
	}

	pcoinfrom, err := c.getCoinPrice(coinfromStr)
	if err != nil {
		return 0, err
	}
	pcointo, err := c.getCoinPrice(cointoStr)
	if err != nil {
		return 0, err
	}

	return (pcoinfrom * amount) / pcointo, nil
}

func (c CoinConfig) GetPrice(coinStr string) (float64, error) {
	if !c.validCoin(coinStr) {
		return 0, fmt.Errorf("Coin not recognized: %s", coinStr)
	}
	if c.isFiat(coinStr) {
		return 0, fmt.Errorf("Cannot get price of FIAT: %s", coinStr)
	}

	return c.getCoinPrice(coinStr)
}
