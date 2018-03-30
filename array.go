// Code generated by go generate - 2018-03-30 10:15:37.438388468 +0000 UTC
package main

// imports 
import (
 )


// interfaces

// Int
type ArrayInt interface {
	First() ResultInt
	Slice() arrayInt
	Each(func(int))
	Concat(ArrayInt) ArrayInt

	MapInt(func(int) int) ArrayInt
	ReduceInt(func(int, int, ArrayInt) int, int) int

	MapNode(func(int) Node) ArrayNode
	ReduceNode(func(int, int, ArrayInt) Node, Node) Node

	MapString(func(int) string) ArrayString
	ReduceString(func(int, int, ArrayInt) string, string) string
 
}

// Node
type ArrayNode interface {
	First() ResultNode
	Slice() arrayNode
	Each(func(Node))
	Concat(ArrayNode) ArrayNode

	MapInt(func(Node) int) ArrayInt
	ReduceInt(func(Node, int, ArrayNode) int, int) int

	MapNode(func(Node) Node) ArrayNode
	ReduceNode(func(Node, int, ArrayNode) Node, Node) Node

	MapString(func(Node) string) ArrayString
	ReduceString(func(Node, int, ArrayNode) string, string) string
 
}

// String
type ArrayString interface {
	First() ResultString
	Slice() arrayString
	Each(func(string))
	Concat(ArrayString) ArrayString

	MapInt(func(string) int) ArrayInt
	ReduceInt(func(string, int, ArrayString) int, int) int

	MapNode(func(string) Node) ArrayNode
	ReduceNode(func(string, int, ArrayString) Node, Node) Node

	MapString(func(string) string) ArrayString
	ReduceString(func(string, int, ArrayString) string, string) string
 
}
 // end of interfaces



// implements


type arrayInt []int

func NewArrayInt(a ...int)ArrayInt {
	return arrayInt(a)
}

func (a arrayInt) First() ResultInt {
	if len(a) > 0 {
		return OkInt(a[0])
	}
	return ErrInt("Out Of Bound Array Access")
}


func (a arrayInt) Slice() arrayInt {
	return a
}

func (a arrayInt) Each(f func(int)) {
	for _, e := range a { f(e) }
}

func (a arrayInt) Concat(xs ArrayInt) ArrayInt {
	return arrayInt(append(a.Slice(), xs.Slice()...))
}



func (a arrayInt) MapInt(f func(int) int) ArrayInt {
	var r = make(arrayInt, len(a))
	for i, e := range a { r[i] = f(e) }
	return arrayInt(r)
}

func (a arrayInt) ReduceInt(f func(int, int, ArrayInt) int, initial int) int {
	var r int = initial
	for i, e := range a { 
		r = f(e, i, a)
	}
	return r
}



func (a arrayInt) MapNode(f func(int) Node) ArrayNode {
	var r = make(arrayNode, len(a))
	for i, e := range a { r[i] = f(e) }
	return arrayNode(r)
}

func (a arrayInt) ReduceNode(f func(int, int, ArrayInt) Node, initial Node) Node {
	var r Node = initial
	for i, e := range a { 
		r = f(e, i, a)
	}
	return r
}



func (a arrayInt) MapString(f func(int) string) ArrayString {
	var r = make(arrayString, len(a))
	for i, e := range a { r[i] = f(e) }
	return arrayString(r)
}

func (a arrayInt) ReduceString(f func(int, int, ArrayInt) string, initial string) string {
	var r string = initial
	for i, e := range a { 
		r = f(e, i, a)
	}
	return r
}

 // end of Int



type arrayNode []Node

func NewArrayNode(a ...Node)ArrayNode {
	return arrayNode(a)
}

func (a arrayNode) First() ResultNode {
	if len(a) > 0 {
		return OkNode(a[0])
	}
	return ErrNode("Out Of Bound Array Access")
}


func (a arrayNode) Slice() arrayNode {
	return a
}

func (a arrayNode) Each(f func(Node)) {
	for _, e := range a { f(e) }
}

func (a arrayNode) Concat(xs ArrayNode) ArrayNode {
	return arrayNode(append(a.Slice(), xs.Slice()...))
}



func (a arrayNode) MapInt(f func(Node) int) ArrayInt {
	var r = make(arrayInt, len(a))
	for i, e := range a { r[i] = f(e) }
	return arrayInt(r)
}

func (a arrayNode) ReduceInt(f func(Node, int, ArrayNode) int, initial int) int {
	var r int = initial
	for i, e := range a { 
		r = f(e, i, a)
	}
	return r
}



func (a arrayNode) MapNode(f func(Node) Node) ArrayNode {
	var r = make(arrayNode, len(a))
	for i, e := range a { r[i] = f(e) }
	return arrayNode(r)
}

func (a arrayNode) ReduceNode(f func(Node, int, ArrayNode) Node, initial Node) Node {
	var r Node = initial
	for i, e := range a { 
		r = f(e, i, a)
	}
	return r
}



func (a arrayNode) MapString(f func(Node) string) ArrayString {
	var r = make(arrayString, len(a))
	for i, e := range a { r[i] = f(e) }
	return arrayString(r)
}

func (a arrayNode) ReduceString(f func(Node, int, ArrayNode) string, initial string) string {
	var r string = initial
	for i, e := range a { 
		r = f(e, i, a)
	}
	return r
}

 // end of Node



type arrayString []string

func NewArrayString(a ...string)ArrayString {
	return arrayString(a)
}

func (a arrayString) First() ResultString {
	if len(a) > 0 {
		return OkString(a[0])
	}
	return ErrString("Out Of Bound Array Access")
}


func (a arrayString) Slice() arrayString {
	return a
}

func (a arrayString) Each(f func(string)) {
	for _, e := range a { f(e) }
}

func (a arrayString) Concat(xs ArrayString) ArrayString {
	return arrayString(append(a.Slice(), xs.Slice()...))
}



func (a arrayString) MapInt(f func(string) int) ArrayInt {
	var r = make(arrayInt, len(a))
	for i, e := range a { r[i] = f(e) }
	return arrayInt(r)
}

func (a arrayString) ReduceInt(f func(string, int, ArrayString) int, initial int) int {
	var r int = initial
	for i, e := range a { 
		r = f(e, i, a)
	}
	return r
}



func (a arrayString) MapNode(f func(string) Node) ArrayNode {
	var r = make(arrayNode, len(a))
	for i, e := range a { r[i] = f(e) }
	return arrayNode(r)
}

func (a arrayString) ReduceNode(f func(string, int, ArrayString) Node, initial Node) Node {
	var r Node = initial
	for i, e := range a { 
		r = f(e, i, a)
	}
	return r
}



func (a arrayString) MapString(f func(string) string) ArrayString {
	var r = make(arrayString, len(a))
	for i, e := range a { r[i] = f(e) }
	return arrayString(r)
}

func (a arrayString) ReduceString(f func(string, int, ArrayString) string, initial string) string {
	var r string = initial
	for i, e := range a { 
		r = f(e, i, a)
	}
	return r
}

 // end of String

 // end of implements


