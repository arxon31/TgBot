package main

import (
	tgClient "TgBot/clients/telegram"
	event_consumer "TgBot/consumer/event-consumer"
	telegram "TgBot/events/telegram"
	"TgBot/storage/sqlite"
	"context"
	"flag"
	"log"
)

const (
	host              = "api.telegram.org"
	sqliteStoragePath = "data/sqlite/storage.db"
	batchSize         = 100
)

func main() {
	//s := files.New(storagePath)

	s, err := sqlite.New(sqliteStoragePath)

	if err != nil {
		log.Fatalf("can't connect to storage: %w", err)
	}

	if err := s.Init(context.TODO()); err != nil {
		log.Fatalf("can't init storage: %v", err)
	}

	eventsProcessor := telegram.New(
		tgClient.New(host, mustToken()),
		s,
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatalf("can't start consumer: %v", err)
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
		log.Fatalf("Please enter the BOT TOKEN")
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
