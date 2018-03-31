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
	"strings"

	"github.com/microcosm-cc/bluemonday"
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

const textNodeTag = "__TEXT_NODE_TAG__"

type pair struct {
	k string
	v string
}

type attributes []pair

func NewAttr() attributes {
	return attributes{}
}

func ClassAttr(class string) attributes {
	return NewAttr().Set("class", class)
}

func (a attributes) Set(k string, v string) attributes {
	return append(a, pair{k, v})
}

type TextNode interface {
	Content() string
}

type textNode struct {
	content string
}

func (t textNode) Content() string {
	return t.content
}

type Node interface {
	Text() TextNode
	Tag() string
	Attrs() attributes
	SetAttr(string, string) Node
	Children() ArrayNode
	Append(n ...Node) Node
	Render() string
}

type node struct {
	txt      TextNode
	tag      string
	attrs    attributes
	children ArrayNode
}

type raw string

func NewRawNode(s string) raw {
	return raw(s)
}

func (r raw) Text() TextNode {
	return nil
}
func (r raw) Tag() string {
	return "section"
}
func (r raw) Attrs() attributes {
	return NewAttr()
}
func (r raw) SetAttr(string, string) Node {
	return r
}
func (r raw) Children() ArrayNode {
	return NewArrayNode()
}
func (r raw) Append(ns ...Node) Node {
	return r
}
func (r raw) Render() string {
	d := strings.Join(strings.Split(string(r), "\n"), "\n\n")
	unsafe := blackfriday.Run([]byte(d))
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	return string(html)
}

func (n node) Text() TextNode {
	return n.txt
}
func (n node) Tag() string {
	return n.tag
}
func (n node) Attrs() attributes {
	return n.attrs
}
func (n node) Children() ArrayNode {
	return n.children
}
func (n *node) Append(ns ...Node) Node {
	n.children = n.children.Concat(NewArrayNode(ns...))
	return n
}
func (n *node) SetAttr(k string, v string) Node {
	n.attrs.Set(k, v)
	return n
}
func renderNode(n node) string {
	// log.Printf("renderNode %v", n)
	var attrs []string
	for _, p := range n.attrs {
		kv := fmt.Sprintf(" %s=\"%s\"", p.k, p.v)
		attrs = append(attrs, kv)
	}

	children := n.children.MapString(func(n Node) string { return n.Render() })

	return fmt.Sprintf("<%s%s>%s</%s>",
		n.tag,
		strings.Join(attrs, " "),
		strings.Join(children.Slice(), "\n"),
		n.tag)
}
func (n node) Render() string {
	if textNodeTag == n.tag {
		return n.txt.Content()
	}
	return renderNode(n)
}

func createNode(tag string, attrs attributes, children ...Node) Node {
	n := new(node)
	n.tag = tag
	n.attrs = attrs
	n.children = NewArrayNode(children...)
	return n
}

type htmlFactoryFunc func(attrs attributes, children ...Node) Node

func htmlFactory(tag string) htmlFactoryFunc {
	f := func(attrs attributes, children ...Node) Node {
		return createNode(tag, attrs, children...)
	}
	return f
}

var (
	NoDisplay = htmlFactory("DIV")(NewAttr().Set("style", "display:none"))
	Div       = htmlFactory("DIV")
	P         = htmlFactory("P")
	Span      = htmlFactory("SPAN")
	H1        = htmlFactory("H1")
	H2        = htmlFactory("H2")
	H3        = htmlFactory("H3")
	A         = htmlFactory("A")
	Pre       = htmlFactory("PRE")
	Img       = htmlFactory("IMG")
	Style     = htmlFactory("STYLE")
	Script    = htmlFactory("SCRIPT")
	HeadLink  = htmlFactory("LINK")
	HeadMeta  = htmlFactory("META")
)

func Text(content string) Node {
	n := new(node)
	n.tag = textNodeTag
	n.txt = textNode{
		content: content,
	}
	return n
}

func Textf(f string, args ...interface{}) Node {
	return Text(fmt.Sprintf(f, args...))
}

type document struct {
	head Node
	body Node
}

func (doc document) Render() string {
	dc := "<!DOCTYPE html>"
	htmlNode := createNode("HTML", NewAttr(), doc.head, doc.body)
	return fmt.Sprintf("%s\n%s", dc, htmlNode.Render())
}

func NewDoc(attrs attributes) document {
	return document{
		head: createNode("HEAD", NewAttr()),
		body: createNode("BODY", attrs),
	}
}
