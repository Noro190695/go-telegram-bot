package events

type IFetcher interface {
	Fetch(limit int) ([]Event, error)
}

type IProcesor interface {
	Process(e Event) error
}

type Type int

const (
	Unknow Type = iota
	Message
)

type Event struct {
	Type Type
	Text string
	Meta interface{}
}
