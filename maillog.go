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
//go:generate futil -type func String=string Int=int
//go:generate futil -type option -import time  String=string UInt64=uint64 Node=Node Time=time.Time Error=error  SerializedPart=SerializedPart Int=int UInt=uint
//go:generate futil -type result -import io  -import net/mail -import github.com/jackc/pgx Bool=bool Node=Node ConnPool=*pgx.ConnPool  Error=error Store=Store  Message=*mail.Message SByte=[]byte String=string SerializedMessage=SerializedMessage Int=int  Reader=io.Reader
//go:generate futil -type array   Int=int String=string Node=Node
//go:generate webgen -output queries.go -what sql -prefix Query
//go:generate webgen -output style.go -what css
//go:generate webgen -output js.go -what js
package main

import (
	"flag"
	"log"
	"math"

	"net/http"
	_ "net/http/pprof"

	"github.com/pierremarc/datasize"
)

const maxInt = int(^uint(0) >> 1)

var (
	configFile  string
	smtpdI      string
	httpdI      string
	seed        bool
	siteName    string
	smtpMaxSize string
)

func GetSiteName() string {
	return siteName
}

func GetMaxSize() int {
	var v datasize.ByteSize
	err := v.UnmarshalText([]byte(smtpMaxSize))
	if err != nil {
		log.Fatalf("Could not parse -max-size: %s", err.Error())
	}
	return int(math.Floor(math.Min(float64(v), float64(maxInt))))
}

func init() {
	flag.BoolVar(&seed, "seed", false, "Regenerate attacments")
	flag.StringVar(&smtpMaxSize, "max-size", "12M", "Maximum message size")
	flag.StringVar(&siteName, "name", "log", "A name for the root link")
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
	defer StopThumbnailer()
	flag.Parse()
	dbc := GetDbConfig(configFile)
	tabs := GetTables(configFile)
	vroot := GetVolume(configFile)
	store := NewStore(dbc, tabs)
	volume := NewVolume(vroot)
	notif := NewNotifier()
	RegisterQueries(store)
	if seed {
		SeedAttachments(store, volume)
	} else {
		cont := make(chan string)
		go StartSMTP(cont, smtpdI, store, volume, notif)
		go StartHTTP(cont, httpdI, store, volume, notif)
		// go profiler()
		controller(cont)
	}
}
