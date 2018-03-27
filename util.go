package main

import (
	"crypto/md5"
	"fmt"
)

func encodedSender(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
