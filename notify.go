package main

import (
	"log"

	"github.com/satori/go.uuid"
)

type Receiver func(Notification)

type Notification struct {
	Topic  string `json:"topic"`
	Record int    `json:"record"`
	Parent int    `json:"parent"`
}

type sucscription struct {
	id uuid.UUID
	fn Receiver
}

func MakeNotification(t string, r int, p int) Notification {
	return Notification{t, r, p}
}

type Notifier struct {
	source  chan Notification
	clients []sucscription
}

func NewNotifier() *Notifier {
	clients := make([]sucscription, 0)
	source := make(chan Notification)
	n := Notifier{source, clients}
	go func() {
		for i := range source {
			for _, s := range n.clients {
				s.fn(i)
			}
		}
	}()

	return &n
}

func (b *Notifier) Notify(n Notification) {
	b.source <- n
}

func (b *Notifier) Subscribe(r Receiver) uuid.UUID {
	id := uuid.Must(uuid.NewV4())
	log.Printf("Notifier.Subscribe %s", id.String())
	b.clients = append(b.clients, sucscription{id, r})
	return id
}

func (b *Notifier) Unsubscribe(id uuid.UUID) {
	log.Printf("Notifier.Unsubscribe %s", id.String())
	clients := make([]sucscription, 0)
	for _, s := range b.clients {
		if !uuid.Equal(id, s.id) {
			clients = append(clients, s)
		}
	}
	b.clients = clients
}
