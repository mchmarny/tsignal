package main

import "time"

// Price struct represents price in single point in time.
type Price struct {
	Symbol   string    `json:"symbol"`
	AskPrice float64   `json:"ask_price"`
	SampleOn time.Time `json:"sample_on"`
}

// StockPrice is used as container for stock price data.
type StockPrice struct {
	Stock Stock  `json:"stock"`
	Price *Price `json:"price"`
}

// Content represents simple tweet content
type Content struct {
	Post   Post   `json:"post"`
	Author Author `json:"author"`
}

// Author represents post author
type Author struct {
	Username    string    `json:"username"`
	FullName    string    `json:"full_name"`
	FriendCount int64     `json:"friend_count"`
	PostCount   int64     `json:"post_count"`
	Source      string    `json:"source"`
	UpdatedOn   time.Time `json:"updated_on"`
}

// Post represents simple tweet content
type Post struct {
	Username       string    `json:"username"`
	PostID         int64     `json:"post_id"`
	PostedOn       time.Time `json:"posted_on"`
	Symbol         string    `json:"symbol"`
	Content        string    `json:"content"`
	SentimentScore float64   `json:"sentiment_score"`
}

// Stock represents simple stcik item
type Stock struct {
	Symbol   string `json:"symbol"`
	Company  string `json:"company"`
	Exchange string `json:"exchange"`
}

// ProviderRerun represents state of the provider job
type ProviderRerun struct {
	Stock   Stock
	Channel chan<- Content
	Error   interface{}
}
