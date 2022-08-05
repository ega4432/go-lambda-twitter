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
	Message string `json:"message"`
}

func tweetHandler(reqBody string) (events.APIGatewayProxyResponse, error) {
	var b RequestBody
	body := []byte(reqBody)
	err := json.Unmarshal(body, &b)

	tweetText := b.Text

	if err != nil || tweetText == "" {
		res := Response{Message: err.Error()}
		jsonBody, _ := json.Marshal(res)

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
		message := fmt.Sprintf("{ \"consumerKey\": %s, \"consumerSecret\": %s, \"accessToken\": %s, \"accessSecret\": %s }", consumerKey, consumerSecret, accessToken, accessSecret)
		res := Response{Message: message}
		jsonBody, _ := json.Marshal(res)
		return events.APIGatewayProxyResponse{
			Body:       string(jsonBody),
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	client := New(consumerKey, consumerSecret, accessToken, accessSecret)

	err = client.Post(tweetText)
	if err != nil {
		res := Response{Message: err.Error()}
		jsonBody, _ := json.Marshal(res)
		return events.APIGatewayProxyResponse{
			Body:       string(jsonBody),
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Tweeted successfully: %s", tweetText),
		StatusCode: http.StatusCreated,
	}, nil
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch request.HTTPMethod {
	case http.MethodGet:
		res := Response{Message: "OK"}
		jsonBody, _ := json.Marshal(res)
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
