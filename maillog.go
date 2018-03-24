//go:generate futil -type option -import time  String=string UInt64=uint64 Node=Node Time=time.Time Error=error
//go:generate futil -type result -import database/sql Bool=bool Node=Node SqlDB=*sql.DB  Error=error Store=Store
//go:generate futil -type array   Int=int String=string
package main

import (
	"flag"
	"log"
)

var configFile string
var smtpdI string
var httpdI string

func init() {
	flag.StringVar(&configFile, "config", "config.json", "configuration file")
	flag.StringVar(&smtpdI, "smtp", "0.0.0.0:2525", "interface for smtpd")
	flag.StringVar(&httpdI, "http", "0.0.0.0:8080", "interface for httpd")
}

func controller(cont chan string) {
	for {
		rec := <-cont
		log.Println(rec)
	}
}

func main() {
	flag.Parse()
	dbc := GetDbConfig(configFile)
	tabs := GetTables(configFile)
	store := NewStore(dbc, tabs)
	cont := make(chan string)
	go StartSMTP(cont, smtpdI, store)
	go StartHTTP(cont, httpdI, store)
	controller(cont)
}
