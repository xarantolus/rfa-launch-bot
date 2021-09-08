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
				// Map: https://mapper.acme.com/?ll=48.34986,10.86273&z=10&t=M&marker0=48.12027%2C10.49881%2C12.3%20km%20WxNW%20of%20Turkheim%20DE&marker1=48.59659%2C11.37909%2C46.9%20km%20ExNE%20of%20Stadtbergen%20DE
				"11.37909,48.59659,10.49881,48.12027",
			},

			// Also include search keywords, that way we get more news
			Track: []string{"rocket factory augsburg"},

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
