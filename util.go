package main

import (
	"crypto/md5"
	"fmt"
	"mime"
	"strings"

	"github.com/labstack/echo"
)

func encodedSender(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

var decod = new(mime.WordDecoder)

func decodeSubject(s string) string {
	if len(s) > 2 {
		s2 := s[:2]
		if "=?" == s2 {
			d, err := decod.DecodeHeader(s)
			if err != nil {
				return s
			}
			return d
		}
	}
	return s
}

func getHostDomain(c echo.Context) string {
	host := c.Request().Host
	parts := strings.Split(host, ":")
	return parts[0]
}
