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
	RawMails string
	Answers  string
}

type Config interface {
	Db() DbConfig
	Tables() Tables
}

type config struct {
	Host     string `json:host`
	User     string `json:user`
	Name     string `json:name`
	Password string `json:password`

	RawMails string `rawmails`
	Answers  string `answers`
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
		RawMails: c.RawMails,
		Answers:  c.Answers,
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
