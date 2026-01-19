package utils

type Event struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type Broker struct {
	Clients map[chan Event]bool
	Add     chan chan Event
	Remove  chan chan Event
	Events  chan Event
}

func NewBroker() *Broker {
	return &Broker{
		Clients: make(map[chan Event]bool),
		Add:     make(chan chan Event),
		Remove:  make(chan chan Event),
		Events:  make(chan Event),
	}
}

func (b *Broker) Start() {
	go func() {
		for {
			select {
			case c := <-b.Add:
				b.Clients[c] = true

			case c := <-b.Remove:
				delete(b.Clients, c)
				close(c)

			case event := <-b.Events:
				for c := range b.Clients {
					c <- event
				}
			}
		}
	}()
}
