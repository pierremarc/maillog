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
	"crypto/md5"
	"fmt"
	"mime"
	"os"
	"strconv"
	"strings"
	"time"

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

func formatTimeDate(t time.Time) string {
	y, m, d := t.Date()
	return fmt.Sprintf("%d-%d-%d", y, m, d)
}

func formatTime(t time.Time) string {
	y, mo, d := t.Date()
	h, mi, _ := t.Clock()
	return fmt.Sprintf("%d-%d-%d %02d:%02d", y, mo, d, h, mi)
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
	domains := ArrayStringFrom(tos...)
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

func getMediaType(contentType string) string {
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "*/*"
	}
	return mediaType
}

func getMainType(mt string) string {
	return strings.Split(mt, "/")[0]
}

func ensureDir(dir string) {
	os.MkdirAll(dir, os.ModePerm)
}

func pathExists(p string) bool {
	_, err := os.Stat(p)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func isSecretTopic(t string) bool {
	return strings.IndexRune(t, '_') == 0
}
