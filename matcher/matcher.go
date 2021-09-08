package matcher

import (
	"github.com/xarantolus/rfa-launch-bot/bot"
	"github.com/xarantolus/rfa-launch-bot/collector"
)

type Matcher struct {
	IgnoredUsers *bot.UserList
}

func (m *Matcher) Match(tweet collector.TweetWrapper) bool {

	return false
}
