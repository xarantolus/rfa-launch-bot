package collector

import (
	"log"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/xarantolus/rfa-launch-bot/util"
)

// LocationStream checks out tweets from a large area around Augsburg
func LocationStream(client *twitter.Client, tweetChan chan<- TweetWrapper) {
	defer panic("location stream ended even though it never should")

	var backoff int = 1
	for {
		s, err := client.Streams.Filter(&twitter.StreamFilterParams{
			Locations: []string{
				// This is a large area around Augsburg
				// Map: http://bboxfinder.com/#48.241138,10.656738,48.553887,11.188202
				"10.263290,48.181654,11.581650,48.612938",

				// Andøya Space Center, see https://twitter.com/rfa_space/status/1425335693484625923
				// Map: http://bboxfinder.com/#68.856583,14.864502,69.377411,16.677246
				"14.864502,68.856583,16.677246,69.377411",
			},
			FilterLevel: "none",
			Language:    []string{"de", "en"},
		})
		if util.LogError(err, "location stream") {
			goto sleep
		}

		log.Println("[Twitter] Connected to location stream")

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
				Source: TweetSourceLocationStream,
				Tweet:  *t,
			}
		}

		backoff *= 2

		log.Printf("[Twitter] Location stream ended for some reason, trying again in %d seconds", backoff*5)
	sleep:
		time.Sleep(time.Duration(backoff) * 5 * time.Second)
	}
}
