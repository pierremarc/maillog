package main

type Receiver func(i interface{})

type Notifier struct {
	source  chan interface{}
	clients []Receiver
}

func NewNotifier() *Notifier {
	clients := make([]Receiver, 0)
	source := make(chan interface{})
	n := Notifier{source, clients}
	go func() {
		for i := range source {
			for _, r := range n.clients {
				r(i)
			}
		}
	}()

	return &n
}

func (b *Notifier) Notify(i interface{}) {
	b.source <- i
}

func (b *Notifier) Subscribe(r Receiver) int {
	b.clients = append(b.clients, r)
	return len(b.clients) - 1
}

func (b *Notifier) Unsubscribe(i int) {
	if i > len(b.clients) {
		return
	}
	head := b.clients[:i]
	if i2 := i + 1; i2 < len(b.clients) {
		tail := b.clients[i2:]
		head = append(head, tail...)
	}
	b.clients = head
}
