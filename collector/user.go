package collector

import (
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
)

// User requests tweets from the twitter profile with the given screen name every 1-2 minutes
func User(username string, client *twitter.Client, tweetChan chan<- TweetWrapper) {
	userName := fmt.Sprintf("user (@%s)", username)
	defer panic(userName + " follower stopped processing even though it shouldn't")

	var userTimelineFunc = func(c *twitter.Client, lastTweetID int64) ([]twitter.Tweet, error) {
		t, _, err := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
			ScreenName:     username,
			TweetMode:      "extended",
			ExcludeReplies: twitter.Bool(false),
			SinceID:        lastTweetID,
		})

		return t, err
	}

	pollingCollector(userName, client, tweetChan, userTimelineFunc)
}
