/*
 *  Copyright (C) 2018 Pierre Marchand <pierre.m@atelier-cartographique.be>
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU Affero General Public License as published by
 *  the Free Software Foundation, version 3 of the License.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"log"
	"sync"

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
	cmut    sync.Locker
}

func NewNotifier() *Notifier {
	clients := make([]sucscription, 0)
	source := make(chan Notification)
	mut := sync.Mutex{}
	n := Notifier{source, clients, &mut}
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
	b.cmut.Lock()
	defer b.cmut.Unlock()
	clients := make([]sucscription, 0)
	for _, s := range b.clients {
		if !uuid.Equal(id, s.id) {
			clients = append(clients, s)
		}
	}
	b.clients = clients
}
