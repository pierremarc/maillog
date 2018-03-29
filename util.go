package main

import (
	"crypto/md5"
	"fmt"
	"mime"
	"strconv"
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

func getRecipent(to []string) OptionString {
	if len(to) > 0 {
		addr := to[0]
		parts := strings.Split(addr, "@")
		return SomeString(parts[0])
	}
	return NoneString()
}

func getDomains(tos []string) ArrayString {
	domains := NewArrayString(tos...)
	return domains.MapString(getDomain)
}

func getDomain(to string) string {
	parts := strings.Split(to, "@")
	if len(parts) > 1 {
		return parts[1]
	}
	return "local"
}

func getTopic(recipient string) OptionString {
	parts := strings.Split(recipient, "+")
	if len(parts) > 0 {
		return SomeString(parts[0])
	}
	return NoneString()
}

func getAnswer(topic string) OptionUInt64 {
	parts := strings.Split(topic, "+")
	if len(parts) > 1 {
		i, err := strconv.ParseUint(parts[1], 10, 32)
		if err == nil {
			return SomeUInt64(i)
		}
	}
	return NoneUInt64()
}
