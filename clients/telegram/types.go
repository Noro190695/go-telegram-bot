package telegram

type Updates struct {
	ID      int              `json:"update_id"`
	Message *IncomingMessage `json:"message"`
}

type UpdatesResponse struct {
	Ok     bool      `json:"ok"`
	Result []Updates `json:"result"`
}

type IncomingMessage struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

type From struct {
	UserName string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}
