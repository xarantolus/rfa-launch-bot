package matcher

import (
	"time"

	"github.com/xarantolus/rfa-launch-bot/bot"
	"github.com/xarantolus/rfa-launch-bot/collector"
	"github.com/xarantolus/rfa-launch-bot/util"
)

type Matcher struct {
	IgnoredUsers *bot.UserList

	seenTweets  map[int64]bool
	maxTweetAge time.Duration

	importantUsers []string

	positiveKeywords []string
	negativeKeywords []string
}

func NewMatcher(ignoredUsers *bot.UserList) (m *Matcher) {
	m = &Matcher{
		IgnoredUsers: ignoredUsers,
		seenTweets:   make(map[int64]bool),

		maxTweetAge: 24 * time.Hour,

		importantUsers: []string{
			"rfa_space",
		},

		positiveKeywords: []string{
			"rocket factory augsburg",
		},

		negativeKeywords: []string{},
	}

	return m
}

func (m *Matcher) Match(tweet collector.TweetWrapper) bool {
	if m.seenTweets[tweet.ID] {
		return false
	}
	m.seenTweets[tweet.ID] = true

	t, terr := tweet.CreatedAtTime()
	if util.LogError(terr, "parsing tweet date") || time.Since(t) > m.maxTweetAge {
		return false
	}

	// Anything from important accounts should be retweeted
	if tweet.User != nil && containsStringCaseInsensitive(m.importantUsers, tweet.User.ScreenName) {
		return true
	}

	// Some accounts are ignored and should of course not be retweeted
	if m.IgnoredUsers.TweetAssociatedWithAny(tweet.Tweet) {
		return false
	}

	return false
}
