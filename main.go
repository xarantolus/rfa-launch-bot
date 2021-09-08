package main

import (
	"flag"
	"log"
	"strconv"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/xarantolus/rfa-launch-bot/bot"
	"github.com/xarantolus/rfa-launch-bot/collector"
	"github.com/xarantolus/rfa-launch-bot/config"
	"github.com/xarantolus/rfa-launch-bot/matcher"
	"github.com/xarantolus/rfa-launch-bot/util"
)

func main() {
	var (
		flagDebug  = flag.Bool("debug", false, "Whether to activate debug mode")
		configFile = flag.String("cfg", "config.yaml", "Path to configuration file")
	)
	flag.Parse()

	cfg, err := config.Parse(*configFile)
	if err != nil {
		log.Fatalf("parsing configuration file: %s\n", err.Error())
	}

	client, user, err := bot.Login(cfg)
	if err != nil {
		log.Fatalf("failed to log in: %s\n", err.Error())
	}
	log.Printf("[Twitter] Logged in @%s\n", user.ScreenName)

	// Load all ignored users
	var ignoredUsers = bot.ListMembers(client, cfg.Lists.NegativeIDs...)

	var matcher = matcher.Matcher{
		IgnoredUsers: ignoredUsers,
	}

	// This channel receives all tweets that should be checked if they are on topic
	var tweetChan = make(chan collector.TweetWrapper, 250)

	// Start all background jobs
	{
		// Timeline
		go collector.Timeline(client, tweetChan)

		// Important users (Company/Investor accounts)
		go collector.User("rfa_space", client, tweetChan)
		go collector.User("OHB_SE", client, tweetChan)

		// All positive lists
		for _, listID := range cfg.Lists.PositiveIDs {
			list, _, err := client.Lists.Show(&twitter.ListsShowParams{
				ListID: listID,
			})
			if util.LogError(err, "loading list details for list with id "+strconv.FormatInt(listID, 10)) {
				continue
			}
			go collector.List(*list, client, tweetChan)
		}
	}

	var retweet = func(t collector.TweetWrapper) {
		if *flagDebug {
			log.Println("Not retweeting", t.URL(), "because we're in debug mode")
			return
		}
		_, _, err := client.Statuses.Retweet(t.ID, &twitter.StatusRetweetParams{})
		util.LogError(err, "retweet")
	}

	for tweet := range tweetChan {
		if matcher.Match(tweet) {
			retweet(tweet)
		}
	}
}
