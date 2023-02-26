# Coin Price

Simple utility made for me to keep up-to-date with coin prices. It uses the API
from coincap, I don't think it works with others APIs.

## Description

My idea was to request the price of a coin, and keep a history of prices.

## Configuration

Uses `$HOME/.config/` path to keep a SQLite DB for the price history, and a
JSON config file.

- `$HOME/.config/coin_prices.json`: some configs variables
- `$HOME/.config/prices.sql`: prices history

### Config template

Check `./sample_coin_prices.json` file, or just copy this:

```
{
  "priceApi": "https://api.coincap.io/v2/assets",
  "exchApi": "https://blockchain.info/ticker",
  "supportedCoins": [
    "bitcoin",
    "monero",
    "litecoin"
  ],
  "supportedFiat": [
    "usd"
  ]
}
```

## Usage

Check `./Makefile` for examples, and once build run `./cpcli -help`.
