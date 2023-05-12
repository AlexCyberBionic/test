package main

import (
	"context"
	"fmt"
	"sync"

	binance "github.com/aiviaio/go-binance/v2"
)

const (
	apiKey = ""
	apiSecret = ""
	assetLimit = 5
	relativePair = "USDT"
)

func main() {
	// Create binance client
	client := binance.NewClient(apiKey, apiSecret)
	// Create exchange service
	ex := client.NewExchangeInfoService()
	// Create price service
	price := client.NewListPricesService()
	// Get exchange info about all assets
	res, err := ex.Do(context.Background())
	if err != nil {
		panic(err)
	}

	var wg  sync.WaitGroup
	resultChan := make(chan map[string]string, assetLimit)
	for _, s := range res.Symbols[:assetLimit] {
		wg.Add(1)
		// Сreate gorutine for each symbol and get price
		go func(symbol string, res chan map[string]string) {
			listSymbolPrice, err := price.Symbol(symbol + relativePair).Do(context.Background())
			if err != nil {
				fmt.Println(err)
			}
			for _, symbolPrice := range listSymbolPrice {
				res <- map[string]string{symbol:symbolPrice.Price}
			}
			wg.Done()
		}(s.BaseAsset, resultChan)
	}
	// Wait all gorutines
	wg.Wait()
	// Close channel
	close(resultChan)

	// Read and print results 
	for resultMap := range resultChan {
		for key, value := range resultMap{
			fmt.Println(key, value)
		}
	}
}