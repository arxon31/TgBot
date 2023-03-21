package main

import (
	"TgBot/clients/telegram"
	"flag"
	"log"
)

func main() {

	tgClient := telegram.New(mustHost(), mustToken())

	// fetcher = fetcher.New(tgClient)

	// processor = processor.New(tgClient)

	// consumer.start(fetcher, processor)
}

func mustToken() string {
	token := flag.String(
		"token",
		"",
		"token for your bot",
	)
	flag.Parse()

	if *token == "" {
		log.Fatal()
	}
	return *token
}

func mustHost() string {
	host := flag.String(
		"host",
		"",
		"host for your bot",
	)
	flag.Parse()

	if *host == "" {
		log.Fatal()
	}
	return *host
}
