package main

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/gosimple/slug"
)

const minAttachmentLen = 64

type WriteOptions struct {
	sender string
	topic  string
	id     int
	fn     string
	data   []byte
}

type Volume interface {
	Write(WriteOptions) ResultString
	Reader(fp string) ResultReader
}

type volume struct {
	root string
}

func NewVolume(root string) volume {
	return volume{root}
}

func (v volume) Write(o WriteOptions) ResultString {
	sid := strconv.Itoa(o.id)
	dir := path.Join(v.root, o.sender, o.topic, sid)
	ensureDir(dir)
	return writeAttachment(dir, o.fn, o.data)
}

func (v volume) Reader(fp string) ResultReader {
	rfp := path.Join(v.root, fp)
	f, err := os.Open(rfp)
	if err != nil {
		return ErrReader(err)
	}
	return OkReader(f)
}

func ensureDir(dir string) {
	os.MkdirAll(dir, os.ModePerm)
}

func writeAttachment(dir string, fn string, data []byte) ResultString {
	if len(fn) == 0 {
		return ErrString("Filename Can't be Empty")
	}
	ext := path.Ext(fn)
	basename := strings.Split(fn, ".")[0]
	fname := slug.Make(basename) + ext
	err := ioutil.WriteFile(path.Join(dir, fname), data, 0644)
	if err != nil {
		return ErrString(err)
	}
	return OkString(fname)
}
