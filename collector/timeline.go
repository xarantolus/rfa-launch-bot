package collector

import (
	"github.com/dghubble/go-twitter/twitter"
)

// Timeline requests the user home timeline about every minute and puts all new tweets in tweetChan.
// it also includes replies which would normally not be shown in the timeline.
func Timeline(client *twitter.Client, tweetChan chan<- TweetWrapper) {
	defer panic("home timeline follower stopped processing even though it shouldn't")

	var homeTimelineFunc = func(c *twitter.Client, lastTweetID int64) ([]twitter.Tweet, error) {
		t, _, err := client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{
			ExcludeReplies:  twitter.Bool(false), // We want to get everything, including replies to tweets
			TrimUser:        twitter.Bool(false), // We care about the user
			IncludeEntities: twitter.Bool(true),  // We do care about who was mentioned etc.
			SinceID:         lastTweetID,         // everything since our last request
			Count:           200,                 // Maximum number of tweets we can get at once
			TweetMode:       "extended",
		})

		return t, err
	}

	pollingCollector("home timeline", TweetSourceTimeline, client, tweetChan, homeTimelineFunc)
}
