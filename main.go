package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/xarantolus/rfa-launch-bot/bot"
	"github.com/xarantolus/rfa-launch-bot/config"
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

	fmt.Println("Logged in user:", user)
	fmt.Println("Client:", client)
	fmt.Println("Debug:", *flagDebug)
}
