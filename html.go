package main

import (
	"fmt"
	"strings"
)

const textNodeTag = "__TEXT_NODE_TAG__"

type pair struct {
	k string
	v string
}

type Attributes []pair

func NewAttr() Attributes {
	return Attributes{}
}

func (a Attributes) Add(k string, v string) Attributes {
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
	Attrs() Attributes
	Children() []Node
	Append(n ...Node) *node
	Render() string
}

type node struct {
	txt      TextNode
	tag      string
	attrs    Attributes
	children []Node
}

func (n node) Text() TextNode {
	return n.txt
}
func (n node) Tag() string {
	return n.tag
}
func (n node) Attrs() Attributes {
	return n.attrs
}
func (n node) Children() []Node {
	return n.children
}

func (n *node) Append(ns ...Node) *node {
	n.children = append(n.children, ns...)
	return n
}

func renderNode(n node) string {
	var children []string
	var attrs []string
	for _, c := range n.children {
		children = append(children, c.Render())
	}
	for _, p := range n.attrs {
		kv := fmt.Sprintf("%s=\"%s\"", p.k, p.v)
		attrs = append(attrs, kv)
	}

	return fmt.Sprintf("<%s %s>%s</%s>",
		n.tag,
		strings.Join(attrs, " "),
		strings.Join(children, "\n"),
		n.tag)
}

func (n node) Render() string {
	if textNodeTag == n.tag {
		return n.txt.Content()
	}
	return renderNode(n)
}

func createNode(tag string, attrs Attributes, children ...Node) Node {
	n := new(node)
	n.tag = tag
	n.attrs = attrs
	n.children = children
	return n
}

type factoryFunc func(attrs Attributes, children ...Node) Node

func factory(tag string) factoryFunc {
	f := func(attrs Attributes, children ...Node) Node {
		return createNode(tag, attrs, children...)
	}
	return f
}

var (
	Div  = factory("DIV")
	P    = factory("P")
	Span = factory("SPAN")
	H1   = factory("H1")
	H2   = factory("H2")
	H3   = factory("H3")
	A    = factory("A")
	Pre  = factory("PRE")
	Img  = factory("IMG")
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

func NewDoc() document {
	return document{
		head: createNode("HEAD", NewAttr()),
		body: createNode("BODY", NewAttr()),
	}
}
