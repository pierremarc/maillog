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
	"fmt"
	"time"
)

var (
	RssRss         = htmlFactory("rss")
	RssChannel     = htmlFactory("channel")
	RssTitle       = htmlFactory("title")
	RssLink        = htmlFactory("link")
	RssDescription = htmlFactory("description")
	RssItem        = htmlFactory("item")
	RssMedia       = htmlFactory("media:content")
	RssPubDate     = htmlFactory("pubDate")
	RssBuildDate   = htmlFactory("lastBuildDate")
	RssGUID        = htmlFactory("guid")
	RssCategory    = htmlFactory("category")
	RssAuthor      = htmlFactory("dc:creator")
	RssAtomLink    = htmlFactory("atom:link")
)

func MakeRSS() Node {
	return RssRss(NewAttr().
		Set("version", "2.0").
		Set("xmlns:media", "http://search.yahoo.com/mrss/").
		Set("xmlns:atom", "http://www.w3.org/2005/Atom").
		Set("xmlns:dc", "http://purl.org/dc/elements/1.1/"))
}

func MakeRssChannel(title string, link string, desc string, self string) Node {
	return RssChannel(NewAttr(),
		RssAtomLink(NewAttr().
			Set("href", self).
			Set("rel", "self").
			Set("type", "application/rss+xml")),
		RssTitle(NewAttr(), Text(title)),
		RssLink(NewAttr(), Text(link)),
		RssDescription(NewAttr(), Text(desc)))
}

func cdata(s string) string {
	return "<![CDATA[" + s + "]]>"
}

func MakeRssItem(topic string, sender string, title string, link string, desc string, t time.Time) Node {
	return RssItem(NewAttr(),
		RssTitle(NewAttr(), Text(cdata(title))),
		RssDescription(NewAttr(), Text(cdata(desc))),
		RssLink(NewAttr(), Text(link)),
		RssAuthor(NewAttr(), Text(cdata(sender))),
		RssCategory(NewAttr(), Text(cdata(topic))),
		RssPubDate(NewAttr(), Text(t.Format(time.RFC822Z))),
		RssGUID(NewAttr().Set("isPermaLink", "true"), Text(link)))
}

const xmlDocElem = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"

func RenderRss(n Node) string {
	return fmt.Sprintf("%s\n%s", xmlDocElem, n.Render())
}
