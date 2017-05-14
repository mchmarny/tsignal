package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

const (
	dataSource          = "Twitter"
	layoutTwitter       = "Mon Jan 02 15:04:05 -0700 2006"
	layoutBigQuery      = "2006-01-02 15:04:05"
	streamLimitPauseMin = 10
)

// start initiates the Tweeter stream subscription and pumps all messages into
// the passed in channel
func subscribeToStream(stock Stock, ch chan<- Content) {

	logInfo.Printf("Subscribing to [%v:%v]...", stock.Symbol, stock.Company)

	consumerKey := os.Getenv("T_CONSUMER_KEY")
	consumerSecret := os.Getenv("T_CONSUMER_SECRET")
	accessToken := os.Getenv("T_ACCESS_TOKEN")
	accessSecret := os.Getenv("T_ACCESS_SECRET")

	if consumerKey == "" || consumerSecret == "" || accessToken == "" || accessSecret == "" {
		logErr.Fatal("Both, consumer key/secret and access token/secret are required")
		return
	}

	// init convif
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)

	// HTTP Client - will automatically authorize Requests
	httpClient := config.Client(appContext, token)
	httpClient.Timeout = time.Duration(30 * time.Second)
	client := twitter.NewClient(httpClient)
	demux := twitter.NewSwitchDemux()

	//Tweet processor
	demux.Tweet = func(tweet *twitter.Tweet) {
		// check if the tweet is a retweet
		if tweet.RetweetedStatus == nil {
			aquiredOn := time.Now()
			username := strings.ToLower(tweet.User.ScreenName)
			msg := Content{
				Post: Post{
					Symbol:   stock.Symbol,
					PostID:   tweet.ID,
					PostedOn: aquiredOn,
					Content:  tweet.Text,
					Username: username,
				},
				Author: Author{
					Username:    username,
					FullName:    tweet.User.Name,
					FriendCount: int64(tweet.User.FollowersCount),
					PostCount:   int64(tweet.User.StatusesCount),
					Source:      dataSource,
					UpdatedOn:   aquiredOn,
				},
			}
			logInfo.Printf("Post [%v:%d]", stock.Symbol, msg.Post.PostID)
			ch <- msg
		}
	}

	// Tweet filter
	filterParams := &twitter.StreamFilterParams{
		Track: []string{
			"#" + stock.Symbol, // hashtag
			"$" + stock.Symbol, // stock sybmob search
			stock.Company,      // just plain name of the company
		},
		StallWarnings: twitter.Bool(true),
		Language:      []string{"en"},
	}

	// Start stream
	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		providerErrors <- ProviderRerun{
			Error:   fmt.Sprintf("Error while creating stream filter: %v", err),
			Channel: ch,
			Stock:   stock,
		}
		return
	}

	demux.StreamLimit = func(limit *twitter.StreamLimit) {
		logErr.Printf("Reached stream limit %v - pausing: %d min",
			limit.Track, streamLimitPauseMin)
		time.Sleep(time.Duration(streamLimitPauseMin * time.Minute))
		providerErrors <- ProviderRerun{
			Error:   fmt.Sprintf("Error while creating stream filter: %v", err),
			Channel: ch,
			Stock:   stock,
		}
		return
	}

	// do the work
	go demux.HandleChan(stream.Messages)
}
