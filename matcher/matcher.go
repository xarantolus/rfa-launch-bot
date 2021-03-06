package matcher

import (
	"log"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/xarantolus/rfa-launch-bot/bot"
	"github.com/xarantolus/rfa-launch-bot/collector"
	"github.com/xarantolus/rfa-launch-bot/util"
)

type Matcher struct {
	IgnoredUsers *bot.UserList

	seenTweets  map[int64]bool
	maxTweetAge time.Duration

	myUserID int64

	client       *twitter.Client
	articleStore *util.ArticleStore

	importantUsers []string

	positiveKeywords         []string
	locationPositiveKeywords []string
	negativeKeywords         []string
}

func NewMatcher(client *twitter.Client, ignoredUsers *bot.UserList, articleStore *util.ArticleStore, myUserID int64) (m *Matcher) {
	m = &Matcher{
		IgnoredUsers: ignoredUsers,
		myUserID:     myUserID,
		client:       client,
		seenTweets:   make(map[int64]bool),
		articleStore: articleStore,

		maxTweetAge: 24 * time.Hour,

		importantUsers: []string{
			"rfa_space",
		},

		positiveKeywords: []string{
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

			"helix engine",
		},

		locationPositiveKeywords: []string{
			"rocket factory", "rocketfactory",

			"statisches feuer", "standfeuer", "dauerfeuer",
			"static fire", "staticfire", "static test fire",
		},

		negativeKeywords: []string{
			"spacex", "blue origin", "blueorigin", "aerojet",
			"electron", "neutron", "rocket lab", "rocketlab", "rklb",
			"falcon", "f9", "starlink", "tesla", "giga press",

			"gigapress", "gigafactory", "openai", "boring", "hyperloop",
			"solarcity", "neuralink", "sls", "nasa_sls", "ula", "vulcan",
			"rogozin", "virgingalactic", "virgin galactic", "virgin orbit",
			"virginorbit", "blueorigin", "boeing", "starliner", "soyuz", "orion",

			"covid", "corona", "pandemie", "pandemic", "impfung", "vaccine", "vax",

			"jesus", "god", "gott", "der herr", "crypto", "eth", "bitcoin", "doge", "coin", "token", "blockchain", "nft", "giveaway",

			"tom sachs", "sachs rocket", "tomsachs",

			"draft register",

			"nintendo", "ring fit adventure", "ringfitadventure",

			"shit", "piss", "anime", "manga", "bronco", "bae",
		},
	}

	return m
}

func (m *Matcher) Match(tweet collector.TweetWrapper) bool {
	if m.tweetMatches(tweet) {
		ignored := m.articleStore.ShouldIgnoreTweet(&tweet.Tweet)
		if ignored {
			log.Printf("Ignoring tweet %s because we've seen an article like this before\n", tweet.URL())
		}
		return !ignored
	}
	return false
}

func (m *Matcher) tweetMatches(tweet collector.TweetWrapper) bool {
	if m.seenTweets[tweet.ID] || tweet.User != nil && tweet.User.ID == m.myUserID {
		return false
	}
	m.seenTweets[tweet.ID] = true

	// t, terr := tweet.CreatedAtTime()
	// if util.LogError(terr, "parsing tweet date") || time.Since(t) > m.maxTweetAge  {
	// 	return false
	// }

	if tweet.Lang != "" && tweet.Lang != "en" && tweet.Lang != "de" && tweet.Lang != "und" {
		return false
	}

	// We don't want to "interrupt" discussions/answers between users by retweeting them
	// However, if someone tweets at themselves (e.g. a thread about space), then it's fine
	if m.isReplyToOtherUser(&tweet.Tweet) {
		return false
	}

	// Anything from important accounts should be retweeted, except for replies to other users
	if tweet.User != nil && containsStringCaseInsensitive(m.importantUsers, tweet.User.ScreenName) {
		return true
	}

	// Some accounts are ignored and should of course not be retweeted
	if m.IgnoredUsers.TweetAssociatedWithAny(tweet.Tweet) {
		return false
	}

	// Now look at the actual tweet text
	text := strings.ToLower(tweet.Text())

	// Some keywords should be ignored
	if anyWordStartsWith(text, m.negativeKeywords...) {
		return false
	}

	// If we have interesting keywords, it's a match
	if anyWordStartsWith(text, m.positiveKeywords...) {
		return true
	}

	// The location stream has additional keywords
	if tweet.Source == collector.TweetSourceLocationStream && anyWordStartsWith(text, m.locationPositiveKeywords...) {
		return true
	}

	return false
}
