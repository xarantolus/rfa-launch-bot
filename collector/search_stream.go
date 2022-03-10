package collector

import (
	"log"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/xarantolus/rfa-launch-bot/util"
)

// SearchStream checks out tweets from a large area around Augsburg
func SearchStream(client *twitter.Client, tweetChan chan<- TweetWrapper) {
	defer panic("search stream ended even though it never should")

	var backoff int = 1
	for {
		s, err := client.Streams.Filter(&twitter.StreamFilterParams{
			Track: []string{
				"rocket factory augsburg",
				"rocketfactory augsburg",
				"rocketfactoryaugsburg",

				"rfa_space",

				"rfa one",
				"rfa1 ",
				"rfa 1 ",
				"rfa launcher",
				"rocket factory one",
				"rocket factory launcher",

				"esrange space center",
				"helix engine",
			},
			FilterLevel: "none",
		})
		if util.LogError(err, "search stream") {
			goto sleep
		}

		log.Println("[Twitter] Connected to search stream")

		// Stream all tweets and serve them to the channel
		for m := range s.Messages {
			backoff = 1
			t, ok := m.(*twitter.Tweet)
			if !ok || t == nil {
				continue
			}

			// If we have truncated text, we try to get the whole tweet
			if t.Truncated {
				t, _, err = client.Statuses.Show(t.ID, &twitter.StatusShowParams{
					TweetMode: "extended",
				})
				if err != nil {
					continue
				}
			}

			tweetChan <- TweetWrapper{
				Source: TweetSourceSearchStream,
				Tweet:  *t,
			}
		}

		backoff *= 2

		log.Printf("[Twitter] Search stream ended for some reason, trying again in %d seconds", backoff*5)
	sleep:
		time.Sleep(time.Duration(backoff) * 5 * time.Second)
	}
}
