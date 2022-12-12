package main

import (
	"coin_prices/internal/api"
	"coin_prices/internal/graph"
	"coin_prices/internal/store"
	"flag"
	"fmt"
	"os"
)

var (
	exch       float64
	coinfrom   string
	cointo     string
	plotPoints bool

	dbFilepath     string
	configFilepath string

	// TODO decide for a name, either plot or graph
	plotFilepath string
	coinClient   api.CoinConfig
)

func init() {
	configureUsage()
	configureFlags()
	validateFlags()

	store.SetDBFilepath(dbFilepath)

	cfg, err := clientFromJSON(configFilepath)
	if err != nil {
		fmt.Printf("Error configuring client:", err)
		os.Exit(1)
	}
	coinClient = cfg
}

func main() {
	if plotPoints {
		if err := plotPointsAndExit(); err != nil {
			fmt.Printf("ERROR: %s\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if err := getAndPrintExchange(); err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}

	if err := printPriceHistory(); err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}

	if err := store.DeleteOlder(coinfrom); err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}

	os.Exit(0)

	/* else if exch == 0 && cointo == "" {
		getAndPrintPrice()
		getPriceHistory()
		os.Exit(0)
	} */
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
	flag.StringVar(&dbFilepath, "db", "/home/al/.config/prices.sql", "DB filepath")
	flag.StringVar(&plotFilepath, "plotfile", "./points.png", "Output filepath for graph")
	flag.BoolVar(&plotPoints, "plot", false, "Plot points in graph and exit")
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

func clientFromJSON(filepath string) (api.CoinConfig, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		os.Exit(1)
	}

	return api.FromJSON(data)
}

func getAndPrintExchange() error {
	f, err := coinClient.GetExch(coinfrom, cointo, exch)
	if err != nil {
		return err
	}

	// Only record the full price of the coin, if exch is not 1 then
	// the history would look weird
	if exch == 1 {
		if err := store.RecordPrice(coinfrom, fmt.Sprintf("%f", f)); err != nil {
			fmt.Printf("WARN: %s\n", err)
		}
	}

	fmt.Printf("%.8f %s are %.8f %s\n", exch, coinfrom, f, cointo)
	return nil
}

/*
func getAndPrintPrice() {
	price, err := coinClient.GetPrice(coinfrom)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
	if err := store.RecordPrice(coinfrom, fmt.Sprintf("%f", price)); err != nil {
		fmt.Printf("WARN: %s\n", err)
	}
	fmt.Println(price)
}
*/

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

func printPriceHistory() error {
	values, err := store.GetPriceHistory(coinfrom)
	if err != nil {
		return err
	}
	for _, row := range values {
		fmt.Printf("%s\t%s\t%s\n", row[0], row[1], row[2])
	}

	return nil
}

func plotPointsAndExit() error {
	values, err := store.GetPriceHistory(coinfrom)
	if err != nil {
		return err
	}

	graph.XLabel = "Points"
	graph.YLabel = coinfrom
	graph.OutputFilepath = plotFilepath

	points, err := graph.PointsFromValues(values)
	if err != nil {
		return err
	}

	if err := graph.GraphPoints(points); err != nil {
		return err
	}

	return nil
}
