package collector

import (
	"log"
	"math/rand"
	"sort"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/xarantolus/rfa-launch-bot/util"
)

type tweetFunc func(c *twitter.Client, lastTweetID int64) ([]twitter.Tweet, error)

// pollingCollectors unifies all functionality of collectors that need to poll something on twitter.
// That way the logic isn't duplicated
func pollingCollector(name string, client *twitter.Client, tweetChan chan<- TweetWrapper, tweetFunc tweetFunc) {

	var (
		// lastSeenID is the ID of the last tweet we saw
		lastSeenID int64

		// The first batch of tweets we receive should not acted upon
		isFirstRequest = true
	)

	log.Println("[Twitter] Start watching " + name)

	for {
		tweets, err := tweetFunc(client, lastSeenID)
		if util.LogError(err, name) {
			goto sleep
		}

		// Sort tweets so the first tweet we process is the oldest one
		sort.Slice(tweets, func(i, j int) bool {
			di, _ := tweets[i].CreatedAtTime()
			dj, _ := tweets[j].CreatedAtTime()

			return dj.After(di)
		})

		for _, tweet := range tweets {
			lastSeenID = tweet.ID

			// We only look at tweets that appeared after the bot started
			if isFirstRequest {
				continue
			}

			tweetChan <- TweetWrapper{
				Source: TweetSourceTimeline,
				Tweet:  tweet,
			}
		}

		if isFirstRequest {
			isFirstRequest = false
		}

	sleep:
		// I guess one request about every minute is ok
		time.Sleep(time.Minute + time.Duration(rand.Intn(60))*time.Second)
	}
}
