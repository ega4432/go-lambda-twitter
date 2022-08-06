package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type RequestBody struct {
	Text string `json:"text"`
}

type Response struct {
	Message   string `json:"message"`
	TweetText string `json:"tweet_text,omitempty"`
	TweetUrl  string `json:"tweet_url,omitempty"`
}

func tweetHandler(reqBody string) (events.APIGatewayProxyResponse, error) {
	var b RequestBody
	body := []byte(reqBody)
	err := json.Unmarshal(body, &b)
	tweetText := b.Text

	if err != nil || tweetText == "" {
		jsonBody, _ := json.Marshal(Response{Message: err.Error()})
		return events.APIGatewayProxyResponse{
			Body:       string(jsonBody),
			StatusCode: http.StatusBadRequest,
		}, err
	}

	consumerKey := os.Getenv("TWITTER_CONSUMER_KEY")
	consumerSecret := os.Getenv("TWITTER_CONSUMER_SECRET")
	accessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	accessSecret := os.Getenv("TWITTER_ACCESS_SECRET")

	if consumerKey == "" || consumerSecret == "" || accessToken == "" || accessSecret == "" {
		message := fmt.Sprintf("consumerKey: %s, consumerSecret: %s, accessToken: %s, accessSecret: %s\n", consumerKey, consumerSecret, accessToken, accessSecret)
		jsonBody, _ := json.Marshal(Response{Message: message})
		return events.APIGatewayProxyResponse{
			Body:       string(jsonBody),
			StatusCode: http.StatusInternalServerError,
		}, errors.New("failed to get environment variable")
	}

	client := New(consumerKey, consumerSecret, accessToken, accessSecret)

	twRes, err := client.Post(tweetText)
	if err != nil {
		jsonBody, _ := json.Marshal(Response{Message: "Failed to post tweet"})
		return events.APIGatewayProxyResponse{
			Body:       string(jsonBody),
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	res := Response{
		Message:   "Tweeted successfully",
		TweetText: tweetText,
		TweetUrl:  fmt.Sprintf("https://twitter.com/ega4432/status/%s", twRes.Data.ID),
	}
	jsonBody, _ := json.Marshal(res)
	return events.APIGatewayProxyResponse{
		Body:       string(jsonBody),
		StatusCode: http.StatusCreated,
	}, nil
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch request.HTTPMethod {
	case http.MethodGet:
		jsonBody, _ := json.Marshal(Response{Message: "OK"})
		return events.APIGatewayProxyResponse{
			Body:       string(jsonBody),
			StatusCode: http.StatusOK,
		}, nil
	case http.MethodPost:
		return tweetHandler(request.Body)
	default:
		return events.APIGatewayProxyResponse{
			Body:       "",
			StatusCode: http.StatusMethodNotAllowed,
		}, errors.New("method not allowed")
	}
}

func main() {
	lambda.Start(handler)
}
