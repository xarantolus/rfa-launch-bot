package main

import (
	"flag"
	"log"
	"strconv"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/xarantolus/rfa-launch-bot/bot"
	"github.com/xarantolus/rfa-launch-bot/collector"
	"github.com/xarantolus/rfa-launch-bot/config"
	"github.com/xarantolus/rfa-launch-bot/generator"
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

	articleStore, err := util.NewArticleStore("articles.json", 12*time.Hour)
	if err != nil {
		log.Printf("[ArticleStore] Failed loading store, but still starting: %s\n", err.Error())
	}

	// Load all ignored & known users
	var (
		knownUsers   = bot.ListMembers(client, "known", cfg.Lists.PositiveIDs...)
		ignoredUsers = bot.ListMembers(client, "ignored", cfg.Lists.NegativeIDs...)
	)

	var matcher = matcher.NewMatcher(client, ignoredUsers, articleStore, user.ID)

	// This channel receives all tweets that should be checked if they are on topic
	var tweetChan = make(chan collector.TweetWrapper, 250)

	// Start all background jobs
	{
		// Timeline
		go collector.Timeline(client, tweetChan)

		// Important users (Company/Investor accounts)
		go collector.User("rfa_space", client, tweetChan)
		go collector.User("OHB_SE", client, tweetChan)

		// Get tweets from around the area
		go collector.LocationStream(client, tweetChan)

		// Get tweets mentioning rfa
		go collector.SearchStream(client, tweetChan)

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

		// And now the jobs that generate tweets by themselves:
		go generator.LiveStreamTweets(client)
	}

	var retweet = func(t collector.TweetWrapper) {
		if *flagDebug {
			log.Println("Not retweeting", t.URL(), "because we're in debug mode")
			return
		}
		_, _, err := client.Statuses.Retweet(t.ID, &twitter.StatusRetweetParams{})
		if util.LogError(err, "retweet") {
			return
		}
		log.Println("[Twitter] Retweeted", t.URL())

		// Add users we don't know yet to a "staging" list. That way, I can add them to the positive list
		if t.User != nil && !knownUsers.ContainsByID(t.User.ID) {
			_, err = client.Lists.MembersCreate(&twitter.ListsMembersCreateParams{
				ListID: cfg.Lists.Staging,
				UserID: t.User.ID,
			})
			if err != nil {
				log.Printf("adding user to list: %s\n", err.Error())
			}
		}
	}

	for tweet := range tweetChan {
		if matcher.Match(tweet) {
			retweet(tweet)
		}
	}
}
