package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"golang.org/x/net/context"
)

const (
	defaultStockRefreshMin = 25
)

var (

	// logs
	logInfo  *log.Logger
	logDebug *log.Logger
	logErr   *log.Logger

	appContext context.Context
	projectID  string
	instanceID string
	dbID       string

	debug           bool
	stockRefreshMin int

	providerErrors        = make(chan ProviderRerun)
	maxProviderErrorCount = 10
	subErrors             = make(map[string]int)
	subMu                 sync.Mutex
)

func main() {

	logInfo = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	logErr = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	// START CONFIG
	flag.StringVar(&projectID, "project", os.Getenv("GCLOUD_PROJECT"),
		"GCP Project ID")
	flag.StringVar(&instanceID, "instance", os.Getenv("GCLOUD_INSTANCE"),
		"GCP Spanner Instance ID")
	flag.StringVar(&dbID, "db", os.Getenv("GCLOUD_DB"), "GCP Spanner DB ID")
	flag.BoolVar(&debug, "debug", debug, "Use verbose logging")
	flag.IntVar(&stockRefreshMin, "refresh", defaultStockRefreshMin,
		"Frequency of stock refresh in minutes")
	flag.Parse()

	if projectID == "" || instanceID == "" || dbID == "" {
		logErr.Fatalf("Missing required argument: project:%v, instance:%v, db:%v",
			projectID, instanceID, dbID)
	}
	logInfo.Printf("CONF - project:%v, instance:%v, db:%v, debug:%v, refresh:%v",
		projectID, instanceID, dbID, debug, stockRefreshMin)
	// END CONFIG

	// configure loggers
	if debug {
		logDebug = log.New(os.Stderr, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		logDebug = log.New(ioutil.Discard, "", log.Ldate)
	}

	// context
	ctx, cancel := context.WithCancel(context.Background())
	appContext = ctx
	go func() {
		// Wait for SIGINT and SIGTERM (HIT CTRL-C)
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		logInfo.Println(<-ch)
		cancel()
		os.Exit(0)
	}()

	// init data store
	storeErr := initStore()
	if storeErr != nil {
		logErr.Fatal(storeErr)
		return
	}

	// get stocks to process
	stocks, err := getStocks()
	if err != nil {
		logErr.Fatalf("Error getting stocks: %v", err)
		return
	}

	// channels for pipeline
	chanells := len(stocks) * 5 // number of channels per stock
	tweets := make(chan Content, chanells)
	results := make(chan Content, chanells)
	pricecheckResult := make(chan bool)

	// update pricess
	go updatePrices(stocks, pricecheckResult)

	// initialize processor
	processErr := initProcessor()
	if processErr != nil {
		logErr.Fatalf("Error from processor: %v", processErr)
		return
	}

	for _, s := range stocks {
		go subscribeToStream(s, tweets)
	}

	// report
	for {
		select {
		case <-appContext.Done():
			break
		case t := <-tweets:
			go processSentiment(t, results)
		case p := <-pricecheckResult:
			logDebug.Printf("Price Check Succeeded: %v", p)
		case t := <-providerErrors:
			checkSubErrors(&t)
			go subscribeToStream(t.Stock, tweets)
		case r := <-results:
			savePost(&r)
			logDebug.Printf("Result [%d:%v]", r.Post.PostID, r.Post.SentimentScore)
		}
	}
}

func checkSubErrors(t *ProviderRerun) {
	logDebug.Printf("Provider rerun for %v", t.Stock.Symbol)
	if subErrors[t.Stock.Symbol] > maxProviderErrorCount {
		logErr.Fatalf("Max numbers of retries for %v reached: %d - %v",
			t.Stock.Symbol,
			subErrors[t.Stock.Symbol],
			t.Error,
		)
	}
	subMu.Lock()
	subErrors[t.Stock.Symbol] = subErrors[t.Stock.Symbol] + 1
	subMu.Unlock()
}
