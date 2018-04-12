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
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type DbConfig struct {
	Host     string
	User     string
	Name     string
	Password string
}

type Tables struct {
	Records     string
	Attachments string
	Domains     string
}

type Config interface {
	Db() DbConfig
	Tables() Tables
}

type config struct {
	Host     string `json:"db/host"`
	User     string `json:"db/user"`
	Name     string `json:"db/name"`
	Password string `json:"db/password"`

	Records     string `json:"tables/records"`
	Attachments string `json:"tables/attachments"`
	Domains     string `json:"tables/domains"`

	Volume string `json:"volume/path"`
}

func (c config) Db() DbConfig {
	return DbConfig{
		Host:     c.Host,
		Name:     c.Name,
		User:     c.User,
		Password: c.Password,
	}
}
func (c config) Tables() Tables {
	return Tables{
		Records:     c.Records,
		Attachments: c.Attachments,
		Domains:     c.Domains,
	}
}

func getConfig(fn string) config {
	file, err := os.Open(fn)
	if err != nil {
		log.Fatalf("Could not Open %v, %v", fn, err)
	}
	c, _ := ioutil.ReadAll(file)
	// log.Printf("%s", c)
	// // file.Seek(0, 0)
	defer file.Close()
	t := config{}
	// var t interface{}
	// decoder := json.NewDecoder(file)
	// err = decoder.Decode(&t)
	err = json.Unmarshal(c, &t)
	if err != nil {
		log.Fatalf("Could not decode %v, %v", fn, err)
	}
	log.Printf("%v", t)
	return t
}

func GetDbConfig(fn string) DbConfig {
	conf := getConfig(fn)
	return conf.Db()

}
func GetTables(fn string) Tables {
	conf := getConfig(fn)
	return conf.Tables()
}
func GetVolume(fn string) string {
	conf := getConfig(fn)
	return conf.Volume
}
