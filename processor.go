package main

import (
	lang "cloud.google.com/go/language/apiv1"
	langpb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

var (
	langClient *lang.Client
)

func initProcessor() error {
	var err error
	langClient, err = lang.NewClient(appContext)
	return err
}

func processSentiment(t Content, r chan<- Content) {
	logDebug.Printf("Processing [%v:%d]...", t.Post.Symbol, t.Post.PostID)
	result, err := langClient.AnalyzeSentiment(appContext, &langpb.AnalyzeSentimentRequest{
		Document: &langpb.Document{
			Source: &langpb.Document_Content{
				Content: t.Post.Content,
			},
			Type: langpb.Document_PLAIN_TEXT,
		},
		EncodingType: langpb.EncodingType_UTF8,
	})
	if err != nil {
		logDebug.Printf("Error envoking NLP API: %v", err)
	} else {
		t.Post.SentimentScore = float64(result.DocumentSentiment.Score)
	}
	r <- t

}
