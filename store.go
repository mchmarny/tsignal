package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

const (
	defaultStockData   = "scripts/stocks.csv"
	stockTable         = "Stocks"
	authorTable        = "Authors"
	postTable          = "Posts"
	postTableTimeIndex = "PostSymbolIndex"
	priceTable         = "Prices"
)

var (
	client                    *spanner.Client
	postPostprocessingWindow  = time.Duration(-1 * time.Hour)
	selectStocksStatement     = spanner.NewStatement("SELECT * FROM Stocks ORDER BY Symbol")
	selectPostsSinceStatement = spanner.NewStatement(
		"SELECT * FROM Posts@{FORCE_INDEX=PostSymbolIndex} WHERE Symbol = @symbol AND PostedOn >= @since")
)

func initStore() error {

	// projects/my-project/instances/my-instance/databases/example-db
	db := fmt.Sprintf("projects/%v/instances/%v/databases/%v",
		projectID, instanceID, dbID)

	logInfo.Printf("Initializing store [%v]...", db)
	c, err := spanner.NewClient(appContext, db)
	if err != nil {
		return err
	}
	client = c

	return loadDefaultData()
}

func savePrice(p *Price) error {

	mP, err := spanner.InsertOrUpdateStruct(priceTable, p)
	if err != nil {
		logErr.Fatal(err)
		return err
	}

	_, insertErr := client.Apply(appContext, []*spanner.Mutation{mP})
	if insertErr != nil {
		logErr.Fatal(insertErr)
	}

	return insertErr

}

func savePost(p *Content) error {

	// insert author
	mAuth, maErr := spanner.InsertOrUpdateStruct(authorTable, p.Author)
	if maErr != nil {
		logErr.Fatal(maErr)
		return maErr
	}

	mPost, mpErr := spanner.InsertOrUpdateStruct(postTable, p.Post)
	if mpErr != nil {
		logErr.Fatal(mpErr)
		return mpErr
	}

	_, insertErr := client.Apply(appContext, []*spanner.Mutation{mAuth, mPost})
	if insertErr != nil {
		logErr.Fatal(insertErr)
	}

	return insertErr

}

func loadDefaultData() error {

	f, err := os.Open(defaultStockData)
	if err != nil {
		logErr.Fatal(err)
		return err
	}
	defer f.Close()

	r := csv.NewReader(f)
	stocks, err := r.ReadAll()
	if err != nil {
		logErr.Fatal(err)
		return err
	}

	records := make([]*spanner.Mutation, len(stocks))
	for i, stock := range stocks {
		record, mErr := spanner.InsertOrUpdateStruct(stockTable, Stock{
			Symbol:   stock[0],
			Company:  stock[1],
			Exchange: stock[2],
		})
		if mErr != nil {
			logErr.Print(mErr)
		}
		records[i] = record
	}

	_, insertErr := client.Apply(appContext, records)
	if insertErr != nil {
		logErr.Fatal(insertErr)
		return insertErr
	}
	return nil

}

func getPostsSince(symbol string, since time.Time) ([]*Post, error) {

	selectPostsSinceStatement.Params["symbol"] = symbol
	selectPostsSinceStatement.Params["since"] = since

	posts := []*Post{}
	iter := client.Single().Query(appContext, selectPostsSinceStatement)
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			logDebug.Print("Done post processing")
			return posts, nil
		}
		if err != nil {
			logInfo.Printf("Error while processing posts: %v", err)
			return posts, err
		}

		var p Post
		if err := row.ToStruct(&p); err != nil {
			logInfo.Printf("Error while processing post row: %v", err)
			return posts, err
		}
		logDebug.Printf("POST PROCESS: %v", p)
		posts = append(posts, &p)
	}
}

func getStocks() ([]Stock, error) {

	stocks := []Stock{}
	iter := client.Single().Query(appContext, selectStocksStatement)
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			return stocks, nil
		}
		if err != nil {
			return stocks, err
		}

		var colSymbol, colComp, colExchange string
		if err := row.Columns(&colSymbol, &colComp, &colExchange); err != nil {
			return stocks, err
		}
		stock := Stock{
			Symbol:   colSymbol,
			Company:  colComp,
			Exchange: colExchange,
		}
		stocks = append(stocks, stock)
	}
}

func query(sql string) error {
	stmt := spanner.Statement{SQL: sql}
	iter := client.Single().Query(appContext, stmt)
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			return err
		}
		var singerID, albumID int64
		var albumTitle string
		if err := row.Columns(&singerID, &albumID, &albumTitle); err != nil {
			return err
		}
		logInfo.Printf("%d %d %s\n", singerID, albumID, albumTitle)
	}
}

func read(table string, args []string) error {
	iter := client.Single().Read(appContext, table, spanner.AllKeys(), args)
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			return err
		}
		var singerID, albumID int64
		var albumTitle string
		if err := row.Columns(&singerID, &albumID, &albumTitle); err != nil {
			return err
		}
		logDebug.Printf("%d %d %s\n", singerID, albumID, albumTitle)
	}
}

/*
func update(ctx context.Context, w io.Writer, client *spanner.Client) error {
	cols := []string{"SingerId", "AlbumId", "MarketingBudget"}
	_, err := client.Apply(ctx, []*spanner.Mutation{
		spanner.Update("Albums", cols, []interface{}{1, 1, 100000}),
		spanner.Update("Albums", cols, []interface{}{2, 2, 500000}),
	})
	return err
}
*/
