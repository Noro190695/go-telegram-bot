package main

import (
	"bot/clients/telegram"
	"flag"
	"fmt"

	"log"
)

const (
	tgHost = "api.telegram.org"
)

func main() {
	token := mustToken()

	tgClient := telegram.New(tgHost, token)
	fmt.Println(tgClient)

	// fetcher fetcher.New()

	// processor processor.New()

	//consumer.Start(fetcher, processor)

}

func mustToken() string {
	token := flag.String("token-bot", "", "token for access telegram bot")
	flag.Parse()
	if *token == "" {
		log.Fatal("token is not corect")
	}
	return *token
}
