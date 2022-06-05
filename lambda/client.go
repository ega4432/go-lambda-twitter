package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dghubble/oauth1"
)

// Twitter API v2 tweet endpoint
// ref. https://developer.twitter.com/en/docs/twitter-api/tweets/manage-tweets/api-reference/post-tweets
var endpoint string = "https://api.twitter.com/2/tweets"

type Reply struct {
	ExcludeReplyUserIds []string `json:"exclude_reply_user_ids,omitempty"`
	InReplyToTweetId    string   `json:"in_reply_to_tweet_id,omitempty"`
}

type Poll struct {
	DurationMinutes int      `json:"duration_minutes"`
	Options         []string `json:"options"`
}

type Media struct {
	MediaIds      []string `json:"media_ids,omitempty"`
	TaggedUserIds []string `json:"tagged_user_ids,omitempty"`
}

type Geo struct {
	PlaceId string `json:"place_id"`
}

type TweetRequest struct {
	Text                  string `json:"text"`
	ReplySettings         string `json:"reply_settings,omitempty"`
	Reply                 *Reply `json:"reply,omitempty"`
	QuoteTweetId          string `json:"quote_tweet_id,omitempty"`
	Poll                  *Poll  `json:"poll,omitempty"`
	ForSuperFollowersOnly bool   `json:"for_super_followers_only,omitempty"`
	DirectMessageDeepLink string `json:"direct_message_deep_link,omitempty"`
}

type Data struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type TweetResponse struct {
	Data *Data `json:"data"`
}

type Client struct {
	Config *oauth1.Config
	Token  *oauth1.Token
}

func New(consumerKey, consumerSecret, accessToken, accessSecret string) *Client {
	return &Client{
		Config: oauth1.NewConfig(consumerKey, consumerSecret),
		Token:  oauth1.NewToken(accessToken, accessSecret),
	}
}

func (c *Client) Post(tweetText string) error {
	httpClient := c.Config.Client(oauth1.NoContext, c.Token)
	body := &TweetRequest{Text: tweetText}
	buf, err := json.Marshal(body)
	if err != nil {
		return err
	}

	res, err := httpClient.Post(endpoint, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return errors.New(fmt.Sprintln(res))
	}

	var twRes TweetResponse
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(resBody, &twRes)
	if err != nil {
		return err
	}

	log.Printf("[INFO]: tweet result\n\ttweet id: %s\n\ttweet text: %s\n", twRes.Data.ID, twRes.Data.Text)
	return nil
}
