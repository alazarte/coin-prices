package main

import (
	"coin_prices/internal/api"
	"flag"
	"fmt"
	"os"
)

var (
	exch     float64
	coinfrom string
	cointo   string
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Printf(`Either use -f alone to get the coin price
Add -t and -d to get price exchange
  cpcli -c bitcoin => bitcoin price
  cpcli -c bitcoin -t monero -d 100 => how many bitcoin I get with 100 monero
`)
	}

	flag.Float64Var(&exch, "d", 1, "Specify amount to convert")
	flag.StringVar(&coinfrom, "f", "usd", "Specify coin (BTC, ETH, ...)")
	flag.StringVar(&cointo, "t", "usd", "To use with -d (BTC, ETH, USD, ...)")
	flag.Parse()
}

func main() {
	if coinfrom != "" {
		if exch != 0 && cointo != "" {
			// TODO make this from cmdline
			f, err := api.GetExch(coinfrom, cointo, exch)
			if err != nil {
				fmt.Printf("ERROR: %s\n", err)
				os.Exit(1)
			}
			fmt.Printf("%.8f %s are %.8f %s\n", exch, coinfrom, f, cointo)
			os.Exit(0)
		} else if exch == 0 && cointo == "" {
			coin, err := api.GetPrice(coinfrom)
			if err != nil {
				fmt.Printf("ERROR: %s\n", err)
				os.Exit(1)
			}
			fmt.Println(coin.PriceUsd)
			os.Exit(0)
		}
	}
	flag.Usage()
}
