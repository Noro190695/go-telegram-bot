package telegram

import (
	"bot/clients/telegram"
	"bot/events"
	"bot/lib/myError"
	"bot/storage"
	"errors"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	UserName string
}

var (
	ErrorUnknownEventType = errors.New("unknown event type")
	ErrorUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      nil,
		offset:  0,
		storage: nil,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)

	if err != nil {
		return nil, myError.Wrap("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return myError.Wrap("can't process message", ErrorUnknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return myError.Wrap("can't process message", err)
	}
	if err := p.doCmd(event.Text, meta.ChatID, meta.UserName); err != nil {
		return myError.Wrap("can't process message", err)
	}
	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, myError.Wrap("can't get meta", ErrorUnknownMetaType)
	}

	return res, nil
}

func event(u telegram.Updates) events.Event {
	typeEvent := fetchType(u)
	res := events.Event{
		Type: typeEvent,
		Text: fetchText(u),
	}

	if typeEvent == events.Message {
		res.Meta = Meta{
			ChatID:   u.Message.Chat.ID,
			UserName: u.Message.From.UserName,
		}
	}
	return res
}

func fetchType(u telegram.Updates) events.Type {
	if u.Message == nil {
		return events.Unknow
	}
	return events.Message
}

func fetchText(u telegram.Updates) string {
	if u.Message == nil {
		return ""
	}
	return u.Message.Text
}
