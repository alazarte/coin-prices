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

	configFilepath string
	coinConfig api.CoinConfig
)

func init() {
	configureUsage()
	configureFlags()
	validateFlags()

	cfg, err := createApiInstance(configFilepath)
	if err != nil {
		fmt.Printf("Error instancing api: %s\n", err)
		os.Exit(1)
	}
	coinConfig = cfg
}

func main() {
	if exch != 0 && cointo != "" {
		getAndPrintExchange()
		os.Exit(0)
	} else if exch == 0 && cointo == "" {
		getAndPrintPrice()
		os.Exit(0)
	}
}

func configureUsage() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Printf(`Either use -f alone to get the coin price
Add -t and -d to get price exchange
  cpcli -c bitcoin => bitcoin price
  cpcli -c bitcoin -t monero -d 100 => how many bitcoin I get with 100 monero
`)
	}

}

func configureFlags() {
	flag.Float64Var(&exch, "d", 1, "Specify amount to convert")
	flag.StringVar(&coinfrom, "f", "usd", "Specify coin (BTC, ETH, ...)")
	flag.StringVar(&cointo, "t", "usd", "To use with -d (BTC, ETH, USD, ...)")
	flag.StringVar(&configFilepath, "c", "", "Config filepath in json")
	flag.Parse()

}

func validateFlags() {
	if coinfrom == "" {
		flag.Usage()
		os.Exit(1)
	}
	if configFilepath == "" {
		if configFilepath = searchAndGetConfigFilepath(); configFilepath == "" {
			fmt.Println("No config file found")
			os.Exit(1)
		}
	}
}

func createApiInstance(filepath string) (api.CoinConfig, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		os.Exit(1)
	}

	return api.FromJSON(data)
}

func getAndPrintExchange() {
	// TODO make this from cmdline
	f, err := coinConfig.GetExch(coinfrom, cointo, exch)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("%.8f %s are %.8f %s\n", exch, coinfrom, f, cointo)
}

func getAndPrintPrice() {
	price, err := coinConfig.GetPrice(coinfrom)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(price)
}

func searchAndGetConfigFilepath() string {
	homepath := os.Getenv("HOME")
	if homepath == "" {
		return ""
	}

	formatPath := fmt.Sprintf("%s/.config/coin_prices.json", homepath)
	if _, err := os.Stat(formatPath); err != nil {
		return ""
	}

	return formatPath
}
