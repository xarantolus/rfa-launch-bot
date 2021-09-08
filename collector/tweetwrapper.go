package collector

import "github.com/dghubble/go-twitter/twitter"

type TweetWrapper struct {
	Source TweetSource

	twitter.Tweet
}

func (t TweetWrapper) URL() string {
	if t.User == nil {
		return "https://twitter.com/i/status/" + t.IDStr
	}
	return "https://twitter.com/" + t.User.ScreenName + "/status/" + t.IDStr

}

type TweetSource string

const (
	TweetSourceLocationStream = "location"
	TweetSourceSearchStream   = "search"
	TweetSourceTimeline       = "timeline"
	TweetSourceUser           = "user"
	TweetSourceList           = "list"
)
