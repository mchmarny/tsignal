package main

import (
	"encoding/csv"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	// NOTE: At this time only NASDAQ traded stocks are supported
	//       Yahoo has disabled NYSE for some reason while back
	yahooFinURL = "http://finance.yahoo.com/d/quotes.csv?s=%s&f=%s"

	// formatting for onlt the ask price [a]
	baseFormat = "a"
)

func isStockMarketOpened() bool {
	est, _ := time.LoadLocation("EST")
	now := time.Now().In(est)
	day := now.Weekday()
	hr := now.Hour()

	// TODO: Market is open from 9:30
	return day != 6 && day != 7 && hr >= 9 && hr < 16
}

// updatePrices updates process for all stocks
func updatePrices(stocks []Stock, r chan<- bool) {
	var numOfErrors int
	var mu sync.Mutex
	pauseDuration := time.Duration(stockRefreshMin) * time.Minute
	for {
		if isStockMarketOpened() {
			mu.Lock()
			numOfErrors = 0
			mu.Unlock()
			for _, stock := range stocks {
				price, err := getAskPrice(stock.Symbol)
				if err != nil {
					mu.Lock()
					numOfErrors++
					mu.Unlock()
					logErr.Printf("Error getting prices for: %v - %v", stock.Symbol, err)
					continue
				}
				// save the current price
				inserErr := savePrice(price)
				if inserErr != nil {
					mu.Lock()
					numOfErrors++
					mu.Unlock()
					logErr.Printf("Error on insert: %v - %v", stock.Symbol, inserErr)
					continue
				}
			}
			r <- numOfErrors == 0
		} else {
			logDebug.Print("Stock market closed")
		}
		time.Sleep(pauseDuration)
	}
}

// getAskPrice single stock price data
func getAskPrice(symbol string) (*Price, error) {

	data, err := loadSingleStockPrice(symbol)
	if err != nil {
		return nil, err
	}

	val, fpErr := strconv.ParseFloat(data[0], 64)
	if fpErr != nil {
		logErr.Printf("Error while parsing prices: %v - %v", data[0], fpErr)
		return nil, fpErr
	}

	price := &Price{
		AskPrice: val,
		Symbol:   symbol,
		SampleOn: time.Now(),
	}

	return price, nil
}

func getHTTPClient(symbol string) (*http.Client, error) {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	c := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(30 * time.Second),
	}
	return c, nil
}

// Single stock data request
func loadSingleStockPrice(symbol string) ([]string, error) {

	c, cErr := getHTTPClient(symbol)
	if cErr != nil {
		logErr.Printf("Error creating yahoo client %v", cErr)
		return nil, cErr
	}

	url := fmt.Sprintf(yahooFinURL, symbol, baseFormat)
	logDebug.Printf("Stock URL: %v", url)

	resp, err := c.Get(url)
	if err != nil {
		logErr.Printf("Error creating yahoo client %v", err)
		return nil, err
	}

	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)
	data, err := reader.Read()
	if err != nil {
		return nil, err
	}
	logDebug.Printf("Stock Prices Data: %v", data)
	return data, err
}
