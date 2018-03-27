//go:generate futil -type option -import time  String=string UInt64=uint64 Node=Node Time=time.Time Error=error  SerializedPart=SerializedPart
//go:generate futil -type result -import io  -import errors -import net/mail -import github.com/jackc/pgx Bool=bool Node=Node ConnPool=*pgx.ConnPool  Error=error Store=Store  Message=*mail.Message SByte=[]byte String=string SerializedMessage=SerializedMessage Int=int  Reader=io.Reader
//go:generate futil -type array   Int=int String=string
package main

import (
	"flag"
	"log"

	"net/http"
	_ "net/http/pprof"
)

var configFile string
var smtpdI string
var httpdI string
var migrate bool

func init() {
	flag.BoolVar(&migrate, "migrate", false, "rather migrate")
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

func profiler() {
	log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
}

func main() {
	flag.Parse()
	dbc := GetDbConfig(configFile)
	tabs := GetTables(configFile)
	vroot := GetVolume(configFile)
	store := NewStore(dbc, tabs)
	volume := NewVolume(vroot)
	if migrate {
		MakeMigration(store, volume)
	} else {
		cont := make(chan string)
		go StartSMTP(cont, smtpdI, store, volume)
		go StartHTTP(cont, httpdI, store, volume)
		go profiler()
		controller(cont)
	}
}
