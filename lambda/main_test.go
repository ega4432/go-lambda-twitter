package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	// t.Run("Unable to get IP", func(t *testing.T) {
	// 	_, err := handler(events.APIGatewayProxyRequest{})
	// 	if err == nil {
	// 		t.Fatal("Error failed to trigger with an invalid request")
	// 	}
	// })

	// t.Run("Non 200 Response", func(t *testing.T) {
	// 	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		w.WriteHeader(500)
	// 	}))
	// 	defer ts.Close()

	// 	handler(events.APIGatewayProxyRequest{})
	// 	if err != nil && err.Error() != ErrNon200Response.Error() {
	// 		t.Fatalf("Error failed to trigger with an invalid HTTP response: %v", err)
	// 	}
	// })

	// t.Run("Unable decode IP", func(t *testing.T) {
	// 	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		w.WriteHeader(500)
	// 	}))
	// 	defer ts.Close()

	// 	_, err := handler(events.APIGatewayProxyRequest{})
	// 	if err == nil {
	// 		t.Fatal("Error failed to trigger with an invalid HTTP response")
	// 	}
	// })

	t.Run("Successful Request", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "hello world")
		}))

		defer ts.Close()

		res, err := handler(events.APIGatewayProxyRequest{
			HTTPMethod: "GET",
		})

		if err != nil {
			t.Fatal("Everything should be ok")
		}

		if res.StatusCode != http.StatusOK {
			t.Fatal("Everything should be 200")
		}

		var r Response
		err = json.Unmarshal([]byte(res.Body), &r)

		if err != nil {
			t.Fatal("Failed to parse response body")
		}

		if r.Message != "OK" {
			t.Fatal("Response message should be \"OK\"")
		}
	})

	// TODO: wip
	// t.Run("Successful Request: tweet", func(t *testing.T) {
	// 	var rb RequestBody
	// 	rb.Text = "test message"

	// 	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		w.WriteHeader(http.StatusCreated)
	// 		fmt.Fprintf(w, "Tweeted successfully: %s", rb.Text)
	// 	}))

	// 	defer ts.Close()

	// 	res, err := handler(events.APIGatewayProxyRequest{
	// 		Body:       `{"text": "test message"}`,
	// 		HTTPMethod: "POST",
	// 	})

	// 	if err != nil {
	// 		t.Fatalf("error: %s", err.Error())
	// 	}

	// 	var r Response
	// 	err = json.Unmarshal([]byte(res.Body), &r)

	// 	if err != nil {
	// 		t.Fatalf("error2: %s", err.Error())
	// 	}

	// 	if r.Message != fmt.Sprintf("Tweeted successfully: %s", rb.Text) {
	// 		t.Fatalf("expected: %v. but actual: %v\n", r, res)
	// 	}
	// })
}
