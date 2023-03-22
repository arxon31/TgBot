package main

import (
	tgClient "TgBot/clients/telegram"
	event_consumer "TgBot/consumer/event-consumer"
	telegram "TgBot/events/telegram"
	"TgBot/storage/files"
	"flag"
	"log"
)

const (
	host        = "api.telegram.org"
	storagePath = "storage"
	batchSize   = 100
)

func main() {

	eventsProcessor := telegram.New(
		tgClient.New(host, mustToken()),
		files.New(storagePath),
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal()
	}

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

/*
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
*/
