package generator

import (
	"bytes"
	"errors"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/docker/go-units"
	"github.com/xarantolus/rfa-launch-bot/collector"
	"github.com/xarantolus/rfa-launch-bot/generator/scraper"
	"github.com/xarantolus/rfa-launch-bot/util"
)

const (
	rfaLiveURL = "https://www.youtube.com/channel/UC6PsS67tBgDr5w22ZZSgI9w/live"
)

// LiveStreamTweets tweets when RFA does a live stream
func LiveStreamTweets(client *twitter.Client) {
	var (
		lastTweetedURL      string
		lastTweetedUpcoming bool
	)

	log.Println("[YouTube] Scraping RFA channel for live streams")

	for {
		liveStream := waitForLiveStream(rfaLiveURL)

		currentURL := liveStream.URL()

		var tweetText = describeLiveStream(&liveStream)

		_, du, ok := liveStream.TimeUntil()

		// Check if we already tweeted about this live stream within the last few minutes
		if lastTweetedURL == currentURL && lastTweetedUpcoming == liveStream.IsUpcoming && (!ok || du < waitTime(du)) {
			goto sleep
		}

		// Tweet the stream announcement
		{
			tweet, _, err := client.Statuses.Update(tweetText, nil)
			if util.LogError(err, "tweeting live stream update") {
				goto sleep
			}

			log.Println("[YouTube] Tweeted live stream announcement:", collector.TweetWrapper{Tweet: *tweet}.URL())
		}

	sleep:
		time.Sleep(time.Minute + time.Duration(rand.Intn(60))*time.Second)
	}
}

// waitTime defines how long we need to wait before sending another (same) live stream announcement
func waitTime(durationUntil time.Duration) time.Duration {
	switch {
	case durationUntil < 30*time.Minute:
		return 15 * time.Minute

	case durationUntil < 4*time.Hour:
		return time.Hour

	default:
		return 4 * time.Hour
	}
}

func waitForLiveStream(liveURL string) (lv scraper.LiveVideo) {
	for {
		liveVideo, err := scraper.YouTubeLive(liveURL)
		if err != nil {
			if !errors.Is(err, scraper.ErrNoVideo) {
				log.Println("[YouTube] Unexpected error while scraping YouTube live:", err.Error())
			}

			goto sleep
		}

		// Validate that the video URL will make sense
		if liveVideo.VideoID != "" && liveVideo.Title != "" {
			return liveVideo
		}

	sleep:
		time.Sleep(time.Minute + time.Duration(rand.Intn(60))*time.Second)
	}
}

const streamTweetTemplate = `{{if .IsUpcoming}}Rocket Factory Augsburg live stream starts {{if .HaveStartTime}}in {{.TimeUntil | duration}}{{else}}soon{{end}}:{{else}}Rocket Factory Augsburg is now live on YouTube:{{end}}

{{.Title}}
{{$keywords := (keywords .Title .ShortDescription)}}{{with $keywords}}
{{hashtags $keywords}}{{end}}

{{.URL}}`

var (
	tmplFuncs = map[string]interface{}{
		"hashtags": util.HashTagText,
		"keywords": extractKeywords,
		"duration": func(d time.Duration) string {
			return strings.ToLower(units.HumanDuration(d))
		},
	}
	streamTweetTmpl = template.Must(template.New("streamTweetTemplate").Funcs(tmplFuncs).Parse(streamTweetTemplate))
)

func describeLiveStream(v *scraper.LiveVideo) string {
	var b bytes.Buffer

	t, dur, haveStartTime := v.TimeUntil()

	var data = struct {
		HaveStartTime bool
		TimeUntil     time.Duration
		StartTime     time.Time
		*scraper.LiveVideo
	}{
		HaveStartTime: haveStartTime,
		TimeUntil:     dur,
		StartTime:     t,
		LiveVideo:     v,
	}

	err := streamTweetTmpl.Execute(&b, data)
	if err != nil {
		panic("executing Tweet template: " + err.Error())
	}

	return strings.TrimSpace(b.String())
}

var expectedRegexes = []*regexp.Regexp{
	regexp.MustCompile(`\bRFA\b`),
	regexp.MustCompile(`\bRocket Factory Augsburg\b`),
	regexp.MustCompile(`\bLauncher\b`),
	regexp.MustCompile(`\bRFA Launcher\b`),
	regexp.MustCompile(`\bNewSpace\b`),
	regexp.MustCompile(`\bAugsburg\b`),
}

func extractKeywords(title string, description string) (keywords []string) {
	extr := title + "\n " + description

	for _, matcher := range expectedRegexes {
		res := matcher.FindAllString(extr, 1)
		if len(res) == 0 {
			continue
		}

		if !containsIgnoreCase(keywords, res[0]) {
			keywords = append(keywords, res[0])
		}
	}

	return
}

func containsIgnoreCase(s []string, e string) bool {
	for _, a := range s {
		if strings.EqualFold(a, e) {
			return true
		}
	}
	return false
}
