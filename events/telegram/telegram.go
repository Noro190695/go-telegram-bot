package telegram

import "bot/clients/telegram"

type Processor struct {
	th     *telegram.Client
	offset int
}

// func New(client *telegram.Client, storage) {
// 	client
// }
