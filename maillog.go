package main

import (
	"log"
)

func controller(cont chan string) {
	for {
		rec := <-cont
		log.Println(rec)
	}
}

func main() {
	store := NewStore(DbConfig{
		Host:     "10.99.99.51",
		Name:     "mailog",
		User:     "mailog",
		Password: "plokplok",
	})
	cont := make(chan string)
	go StartSMTP(cont, store)
	go StartHTTP(cont, store)
	controller(cont)
}
