package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Body struct {
	Text string `json:"text"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var b Body
	body := []byte(request.Body)
	err := json.Unmarshal(body, &b)

	if err != nil || b.Text == "" {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: http.StatusBadRequest,
		}, err
	}

	consumerKey := os.Getenv("TWITTER_CONSUMER_KEY")
	consumerSecret := os.Getenv("TWITTER_CONSUMER_SECRET")
	accessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	accessSecret := os.Getenv("TWITTER_ACCESS_SECRET")
	prefix := os.Getenv("PREFIX")

	if consumerKey == "" || consumerSecret == "" || accessToken == "" || accessSecret == "" {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: http.StatusBadRequest,
		}, err
	}

	tweetText := strings.Join([]string{prefix, b.Text}, " ")
	client := New(consumerKey, consumerSecret, accessToken, accessSecret)

	err = client.Post(tweetText)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       "Tweeted successfully: " + tweetText,
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(handler)
}
