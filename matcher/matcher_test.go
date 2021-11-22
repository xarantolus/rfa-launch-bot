package matcher

import (
	"strings"
	"testing"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/xarantolus/rfa-launch-bot/collector"
)

func TestKeywordCase(t *testing.T) {
	var m = NewMatcher(nil, nil, 0)

	for _, kw := range m.positiveKeywords {
		if strings.ToLower(kw) != kw {
			t.Errorf("Keyword %q is not lowercase in positiveKeywords, but should be", kw)
		}
	}
	for _, kw := range m.locationPositiveKeywords {
		if strings.ToLower(kw) != kw {
			t.Errorf("Keyword %q is not lowercase in locationPositiveKeywords, but should be", kw)
		}
	}
	for _, kw := range m.negativeKeywords {
		if strings.ToLower(kw) != kw {
			t.Errorf("Keyword %q is not lowercase in negativeKeywords, but should be", kw)
		}
	}
	for _, kw := range m.importantUsers {
		if strings.ToLower(kw) != kw {
			t.Errorf("Name %q is not lowercase in importantUsers, but should be", kw)
		}
	}
}

func TestMatchingTweetPositive(t *testing.T) {
	type tweet struct {
		Username string
		Text     string
	}

	var m = NewMatcher(nil, nil, 0)

	var shouldMatch = []tweet{
		{
			Text: "The rocket factory augsburg is rapidly advancing germanys space access",
		},
	}

	for _, tw := range shouldMatch {
		t.Run(t.Name(), func(t *testing.T) {
			matched := m.Match(collector.TweetWrapper{
				Source: collector.TweetSourceTimeline,
				Tweet: twitter.Tweet{
					FullText: tw.Text,
				},
			})

			if !matched {
				if tw.Username == "" {
					t.Errorf("Didn't match tweet with text %q, but should have", tw.Text)
				} else {
					t.Errorf("Didn't match tweet with text %q by %q, but should have", tw.Text, tw.Username)
				}
			}
		})
	}
}

func TestMatchingTweetNegative(t *testing.T) {
	type tweet struct {
		Username string
		Text     string
	}

	var m = NewMatcher(nil, nil, 0)

	var shouldNotMatch = []tweet{
		{
			Text: "tom sachs rocket factory",
		},
		{
			Text: "One project I’ve been exploring is @Irrelevants_NFT.\nUnique art.\nWonderful team.\nFun arcade roadmap with collabs and utility\n“Build-a-bot” functionality similar to Tom Sachs Rocket Factory\n\nUniqueness and simplicity to users. These will be winning projects.",
		},
		{
			Text: `Rocket Factory IRL Scavenger Hunt-Official Rules: 18 numbered clue/update threads will be added to the TSRF Twitter one at a time. Each clue = 1 Location, 1 Component. When found: post + tag a photo & DM @tsrocketfactory to exchange the Physical Component for its matching NFT Rocket`,
		},
	}

	for _, tw := range shouldNotMatch {
		t.Run(t.Name(), func(t *testing.T) {
			matched := m.Match(collector.TweetWrapper{
				Source: collector.TweetSourceTimeline,
				Tweet: twitter.Tweet{
					FullText: tw.Text,
				},
			})

			if matched {
				if tw.Username == "" {
					t.Errorf("Matched tweet with text %q, but shouldn't have done that", tw.Text)
				} else {
					t.Errorf("Matched tweet with text %q by %q, but shouldn't have done that", tw.Text, tw.Username)
				}
			}
		})
	}
}
