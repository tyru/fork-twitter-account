package util

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
)

type Twitter struct {
	api *anaconda.TwitterApi
}

// https://qiita.com/konojunya/items/d51672f900f4912a5563
func NewTwitter() (*Twitter, error) {
	consumerKey := os.Getenv("TWITTER_CONSUMER_KEY")
	consumerSecret := os.Getenv("TWITTER_CONSUMER_SECRET")
	accessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret := os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")
	if consumerKey == "" {
		return nil, errors.New("environment variable 'TWITTER_CONSUMER_KEY' is empty")
	}
	if consumerSecret == "" {
		return nil, errors.New("environment variable 'TWITTER_CONSUMER_SECRET' is empty")
	}
	if accessToken == "" {
		return nil, errors.New("environment variable 'TWITTER_ACCESS_TOKEN' is empty")
	}
	if accessTokenSecret == "" {
		return nil, errors.New("environment variable 'TWITTER_ACCESS_TOKEN_SECRET' is empty")
	}
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)
	return &Twitter{api: api}, nil
}

type UserResult struct {
	ScreenName string
	Error      error
}

// GetFollowingUsers returns channel which sends screen name until the end.
// The stream will be closed at the end.
func (tw *Twitter) GetFollowingUsers(screenName string) chan UserResult {
	ch := make(chan UserResult)
	go func() {
		defer close(ch)
		v := url.Values{}
		v.Set("screen_name", screenName)
		v.Set("count", "200") // maximum value
		for page := range tw.api.GetFriendsListAll(v) {
			if page.Error != nil {
				ch <- UserResult{Error: page.Error}
				return
			}
			for i := range page.Friends {
				ch <- UserResult{ScreenName: page.Friends[i].ScreenName}
			}
		}
	}()
	return ch
}

func (tw *Twitter) FollowUsers(screenName string) error {
	_, err := tw.api.FollowUser(screenName)
	return err
}

func LogInfo(msg string) {
	fmt.Fprintln(os.Stderr, "[INFO]", msg)
}

func LogError(msg string) {
	fmt.Fprintln(os.Stderr, "[ERROR]", msg)
}
