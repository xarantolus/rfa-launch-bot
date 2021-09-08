package collector

import (
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
)

// List requests tweets from the twitter list with the given ID every 1-2 minutes
func List(list twitter.List, client *twitter.Client, tweetChan chan<- TweetWrapper) {
	listName := fmt.Sprintf("list (%s, id=%d)", list.FullName, list.ID)
	defer panic(listName + " follower stopped processing even though it shouldn't")

	var listTimelineFunc = func(c *twitter.Client, lastTweetID int64) ([]twitter.Tweet, error) {
		t, _, err := client.Lists.Statuses(&twitter.ListsStatusesParams{
			ListID:          list.ID,
			IncludeRetweets: twitter.Bool(true),
			IncludeEntities: twitter.Bool(true),
			SinceID:         lastTweetID, // everything since our last request
			Count:           200,         // Maximum number of tweets we can get at once
		})

		return t, err
	}

	pollingCollector(listName, client, tweetChan, listTimelineFunc)
}
