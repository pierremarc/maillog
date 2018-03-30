package main

type Receiver func(Notification)

type Notification struct {
	Topic  string `json:"topic"`
	Record int    `json:"record"`
	Parent int    `json:"parent"`
}

func MakeNotification(t string, r int, p int) Notification {
	return Notification{t, r, p}
}

type Notifier struct {
	source  chan Notification
	clients []Receiver
}

func NewNotifier() *Notifier {
	clients := make([]Receiver, 0)
	source := make(chan Notification)
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

func (b *Notifier) Notify(n Notification) {
	b.source <- n
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
