// Code generated by go generate - 2018-03-27 19:57:33.548777591 +0000 UTC
package main

// imports 
import (
	"time"
)


// interfaces

// Error
type OptionError interface {
	Map(func(error))
	FoldF(func(), func(error))

	MapError(func(error) error) OptionError
	FoldError(error, func(error) error) error
	FoldErrorF(func() error, func(error) error) error

	MapNode(func(error) Node) OptionNode
	FoldNode(Node, func(error) Node) Node
	FoldNodeF(func() Node, func(error) Node) Node

	MapSerializedPart(func(error) SerializedPart) OptionSerializedPart
	FoldSerializedPart(SerializedPart, func(error) SerializedPart) SerializedPart
	FoldSerializedPartF(func() SerializedPart, func(error) SerializedPart) SerializedPart

	MapString(func(error) string) OptionString
	FoldString(string, func(error) string) string
	FoldStringF(func() string, func(error) string) string

	MapTime(func(error) time.Time) OptionTime
	FoldTime(time.Time, func(error) time.Time) time.Time
	FoldTimeF(func() time.Time, func(error) time.Time) time.Time

	MapUInt64(func(error) uint64) OptionUInt64
	FoldUInt64(uint64, func(error) uint64) uint64
	FoldUInt64F(func() uint64, func(error) uint64) uint64
 
}

// Node
type OptionNode interface {
	Map(func(Node))
	FoldF(func(), func(Node))

	MapError(func(Node) error) OptionError
	FoldError(error, func(Node) error) error
	FoldErrorF(func() error, func(Node) error) error

	MapNode(func(Node) Node) OptionNode
	FoldNode(Node, func(Node) Node) Node
	FoldNodeF(func() Node, func(Node) Node) Node

	MapSerializedPart(func(Node) SerializedPart) OptionSerializedPart
	FoldSerializedPart(SerializedPart, func(Node) SerializedPart) SerializedPart
	FoldSerializedPartF(func() SerializedPart, func(Node) SerializedPart) SerializedPart

	MapString(func(Node) string) OptionString
	FoldString(string, func(Node) string) string
	FoldStringF(func() string, func(Node) string) string

	MapTime(func(Node) time.Time) OptionTime
	FoldTime(time.Time, func(Node) time.Time) time.Time
	FoldTimeF(func() time.Time, func(Node) time.Time) time.Time

	MapUInt64(func(Node) uint64) OptionUInt64
	FoldUInt64(uint64, func(Node) uint64) uint64
	FoldUInt64F(func() uint64, func(Node) uint64) uint64
 
}

// SerializedPart
type OptionSerializedPart interface {
	Map(func(SerializedPart))
	FoldF(func(), func(SerializedPart))

	MapError(func(SerializedPart) error) OptionError
	FoldError(error, func(SerializedPart) error) error
	FoldErrorF(func() error, func(SerializedPart) error) error

	MapNode(func(SerializedPart) Node) OptionNode
	FoldNode(Node, func(SerializedPart) Node) Node
	FoldNodeF(func() Node, func(SerializedPart) Node) Node

	MapSerializedPart(func(SerializedPart) SerializedPart) OptionSerializedPart
	FoldSerializedPart(SerializedPart, func(SerializedPart) SerializedPart) SerializedPart
	FoldSerializedPartF(func() SerializedPart, func(SerializedPart) SerializedPart) SerializedPart

	MapString(func(SerializedPart) string) OptionString
	FoldString(string, func(SerializedPart) string) string
	FoldStringF(func() string, func(SerializedPart) string) string

	MapTime(func(SerializedPart) time.Time) OptionTime
	FoldTime(time.Time, func(SerializedPart) time.Time) time.Time
	FoldTimeF(func() time.Time, func(SerializedPart) time.Time) time.Time

	MapUInt64(func(SerializedPart) uint64) OptionUInt64
	FoldUInt64(uint64, func(SerializedPart) uint64) uint64
	FoldUInt64F(func() uint64, func(SerializedPart) uint64) uint64
 
}

// String
type OptionString interface {
	Map(func(string))
	FoldF(func(), func(string))

	MapError(func(string) error) OptionError
	FoldError(error, func(string) error) error
	FoldErrorF(func() error, func(string) error) error

	MapNode(func(string) Node) OptionNode
	FoldNode(Node, func(string) Node) Node
	FoldNodeF(func() Node, func(string) Node) Node

	MapSerializedPart(func(string) SerializedPart) OptionSerializedPart
	FoldSerializedPart(SerializedPart, func(string) SerializedPart) SerializedPart
	FoldSerializedPartF(func() SerializedPart, func(string) SerializedPart) SerializedPart

	MapString(func(string) string) OptionString
	FoldString(string, func(string) string) string
	FoldStringF(func() string, func(string) string) string

	MapTime(func(string) time.Time) OptionTime
	FoldTime(time.Time, func(string) time.Time) time.Time
	FoldTimeF(func() time.Time, func(string) time.Time) time.Time

	MapUInt64(func(string) uint64) OptionUInt64
	FoldUInt64(uint64, func(string) uint64) uint64
	FoldUInt64F(func() uint64, func(string) uint64) uint64
 
}

// Time
type OptionTime interface {
	Map(func(time.Time))
	FoldF(func(), func(time.Time))

	MapError(func(time.Time) error) OptionError
	FoldError(error, func(time.Time) error) error
	FoldErrorF(func() error, func(time.Time) error) error

	MapNode(func(time.Time) Node) OptionNode
	FoldNode(Node, func(time.Time) Node) Node
	FoldNodeF(func() Node, func(time.Time) Node) Node

	MapSerializedPart(func(time.Time) SerializedPart) OptionSerializedPart
	FoldSerializedPart(SerializedPart, func(time.Time) SerializedPart) SerializedPart
	FoldSerializedPartF(func() SerializedPart, func(time.Time) SerializedPart) SerializedPart

	MapString(func(time.Time) string) OptionString
	FoldString(string, func(time.Time) string) string
	FoldStringF(func() string, func(time.Time) string) string

	MapTime(func(time.Time) time.Time) OptionTime
	FoldTime(time.Time, func(time.Time) time.Time) time.Time
	FoldTimeF(func() time.Time, func(time.Time) time.Time) time.Time

	MapUInt64(func(time.Time) uint64) OptionUInt64
	FoldUInt64(uint64, func(time.Time) uint64) uint64
	FoldUInt64F(func() uint64, func(time.Time) uint64) uint64
 
}

// UInt64
type OptionUInt64 interface {
	Map(func(uint64))
	FoldF(func(), func(uint64))

	MapError(func(uint64) error) OptionError
	FoldError(error, func(uint64) error) error
	FoldErrorF(func() error, func(uint64) error) error

	MapNode(func(uint64) Node) OptionNode
	FoldNode(Node, func(uint64) Node) Node
	FoldNodeF(func() Node, func(uint64) Node) Node

	MapSerializedPart(func(uint64) SerializedPart) OptionSerializedPart
	FoldSerializedPart(SerializedPart, func(uint64) SerializedPart) SerializedPart
	FoldSerializedPartF(func() SerializedPart, func(uint64) SerializedPart) SerializedPart

	MapString(func(uint64) string) OptionString
	FoldString(string, func(uint64) string) string
	FoldStringF(func() string, func(uint64) string) string

	MapTime(func(uint64) time.Time) OptionTime
	FoldTime(time.Time, func(uint64) time.Time) time.Time
	FoldTimeF(func() time.Time, func(uint64) time.Time) time.Time

	MapUInt64(func(uint64) uint64) OptionUInt64
	FoldUInt64(uint64, func(uint64) uint64) uint64
	FoldUInt64F(func() uint64, func(uint64) uint64) uint64
 
}


// functions 

func IdError(v error) error {return v}
func OptionErrorFrom(v error, err error) OptionError {
	if err != nil {
		return NoneError()
	}
	return SomeError(v)
}

func IdNode(v Node) Node {return v}
func OptionNodeFrom(v Node, err error) OptionNode {
	if err != nil {
		return NoneNode()
	}
	return SomeNode(v)
}

func IdSerializedPart(v SerializedPart) SerializedPart {return v}
func OptionSerializedPartFrom(v SerializedPart, err error) OptionSerializedPart {
	if err != nil {
		return NoneSerializedPart()
	}
	return SomeSerializedPart(v)
}

func IdString(v string) string {return v}
func OptionStringFrom(v string, err error) OptionString {
	if err != nil {
		return NoneString()
	}
	return SomeString(v)
}

func IdTime(v time.Time) time.Time {return v}
func OptionTimeFrom(v time.Time, err error) OptionTime {
	if err != nil {
		return NoneTime()
	}
	return SomeTime(v)
}

func IdUInt64(v uint64) uint64 {return v}
func OptionUInt64From(v uint64, err error) OptionUInt64 {
	if err != nil {
		return NoneUInt64()
	}
	return SomeUInt64(v)
}



// none
type noneError struct{}

func NoneError() OptionError {
	return noneError{}
}

// map NoneError
func (n noneError) Map(f func(error)) {}
func (n noneError) FoldF(l func(), r func(error)) { l() }

 // map NoneError => OptionError
func (n noneError) MapError(f func(v error) error) OptionError {
	return noneError{}
}
// fold NoneError => Optionerror
func (n noneError) FoldError(a error, f func(v error) error) error {
	return a
}
func (n noneError) FoldErrorF(a func() error, f func(v error) error) error {
	return a()
}

 // map NoneError => OptionNode
func (n noneError) MapNode(f func(v error) Node) OptionNode {
	return noneNode{}
}
// fold NoneError => OptionNode
func (n noneError) FoldNode(a Node, f func(v error) Node) Node {
	return a
}
func (n noneError) FoldNodeF(a func() Node, f func(v error) Node) Node {
	return a()
}

 // map NoneError => OptionSerializedPart
func (n noneError) MapSerializedPart(f func(v error) SerializedPart) OptionSerializedPart {
	return noneSerializedPart{}
}
// fold NoneError => OptionSerializedPart
func (n noneError) FoldSerializedPart(a SerializedPart, f func(v error) SerializedPart) SerializedPart {
	return a
}
func (n noneError) FoldSerializedPartF(a func() SerializedPart, f func(v error) SerializedPart) SerializedPart {
	return a()
}

 // map NoneError => OptionString
func (n noneError) MapString(f func(v error) string) OptionString {
	return noneString{}
}
// fold NoneError => Optionstring
func (n noneError) FoldString(a string, f func(v error) string) string {
	return a
}
func (n noneError) FoldStringF(a func() string, f func(v error) string) string {
	return a()
}

 // map NoneError => OptionTime
func (n noneError) MapTime(f func(v error) time.Time) OptionTime {
	return noneTime{}
}
// fold NoneError => Optiontime.Time
func (n noneError) FoldTime(a time.Time, f func(v error) time.Time) time.Time {
	return a
}
func (n noneError) FoldTimeF(a func() time.Time, f func(v error) time.Time) time.Time {
	return a()
}

 // map NoneError => OptionUInt64
func (n noneError) MapUInt64(f func(v error) uint64) OptionUInt64 {
	return noneUInt64{}
}
// fold NoneError => Optionuint64
func (n noneError) FoldUInt64(a uint64, f func(v error) uint64) uint64 {
	return a
}
func (n noneError) FoldUInt64F(a func() uint64, f func(v error) uint64) uint64 {
	return a()
}
 // end of nones

// some
type someError struct {
	v error
}

func SomeError(v error) someError {
	return someError{v}
}
// map NoneError
func (s someError) Map(f func(error)) { f(s.v) }
func (s someError) FoldF(l func(), r func(error)) { r(s.v) }

// map SoneError => OptionError
func (s someError) MapError(f func(v error) error) OptionError {
	return SomeError(f(s.v))
}
// fold SomeError => error
func (s someError) FoldError(a error, f func(v error) error) error {
	return f(s.v)
}
func (s someError) FoldErrorF(a func() error, f func(v error) error) error {
	return f(s.v)
}

// map SoneError => OptionNode
func (s someError) MapNode(f func(v error) Node) OptionNode {
	return SomeNode(f(s.v))
}
// fold SomeError => Node
func (s someError) FoldNode(a Node, f func(v error) Node) Node {
	return f(s.v)
}
func (s someError) FoldNodeF(a func() Node, f func(v error) Node) Node {
	return f(s.v)
}

// map SoneError => OptionSerializedPart
func (s someError) MapSerializedPart(f func(v error) SerializedPart) OptionSerializedPart {
	return SomeSerializedPart(f(s.v))
}
// fold SomeError => SerializedPart
func (s someError) FoldSerializedPart(a SerializedPart, f func(v error) SerializedPart) SerializedPart {
	return f(s.v)
}
func (s someError) FoldSerializedPartF(a func() SerializedPart, f func(v error) SerializedPart) SerializedPart {
	return f(s.v)
}

// map SoneError => OptionString
func (s someError) MapString(f func(v error) string) OptionString {
	return SomeString(f(s.v))
}
// fold SomeError => string
func (s someError) FoldString(a string, f func(v error) string) string {
	return f(s.v)
}
func (s someError) FoldStringF(a func() string, f func(v error) string) string {
	return f(s.v)
}

// map SoneError => OptionTime
func (s someError) MapTime(f func(v error) time.Time) OptionTime {
	return SomeTime(f(s.v))
}
// fold SomeError => time.Time
func (s someError) FoldTime(a time.Time, f func(v error) time.Time) time.Time {
	return f(s.v)
}
func (s someError) FoldTimeF(a func() time.Time, f func(v error) time.Time) time.Time {
	return f(s.v)
}

// map SoneError => OptionUInt64
func (s someError) MapUInt64(f func(v error) uint64) OptionUInt64 {
	return SomeUInt64(f(s.v))
}
// fold SomeError => uint64
func (s someError) FoldUInt64(a uint64, f func(v error) uint64) uint64 {
	return f(s.v)
}
func (s someError) FoldUInt64F(a func() uint64, f func(v error) uint64) uint64 {
	return f(s.v)
}
 // end of somes


// none
type noneNode struct{}

func NoneNode() OptionNode {
	return noneNode{}
}

// map NoneNode
func (n noneNode) Map(f func(Node)) {}
func (n noneNode) FoldF(l func(), r func(Node)) { l() }

 // map NoneNode => OptionError
func (n noneNode) MapError(f func(v Node) error) OptionError {
	return noneError{}
}
// fold NoneNode => Optionerror
func (n noneNode) FoldError(a error, f func(v Node) error) error {
	return a
}
func (n noneNode) FoldErrorF(a func() error, f func(v Node) error) error {
	return a()
}

 // map NoneNode => OptionNode
func (n noneNode) MapNode(f func(v Node) Node) OptionNode {
	return noneNode{}
}
// fold NoneNode => OptionNode
func (n noneNode) FoldNode(a Node, f func(v Node) Node) Node {
	return a
}
func (n noneNode) FoldNodeF(a func() Node, f func(v Node) Node) Node {
	return a()
}

 // map NoneNode => OptionSerializedPart
func (n noneNode) MapSerializedPart(f func(v Node) SerializedPart) OptionSerializedPart {
	return noneSerializedPart{}
}
// fold NoneNode => OptionSerializedPart
func (n noneNode) FoldSerializedPart(a SerializedPart, f func(v Node) SerializedPart) SerializedPart {
	return a
}
func (n noneNode) FoldSerializedPartF(a func() SerializedPart, f func(v Node) SerializedPart) SerializedPart {
	return a()
}

 // map NoneNode => OptionString
func (n noneNode) MapString(f func(v Node) string) OptionString {
	return noneString{}
}
// fold NoneNode => Optionstring
func (n noneNode) FoldString(a string, f func(v Node) string) string {
	return a
}
func (n noneNode) FoldStringF(a func() string, f func(v Node) string) string {
	return a()
}

 // map NoneNode => OptionTime
func (n noneNode) MapTime(f func(v Node) time.Time) OptionTime {
	return noneTime{}
}
// fold NoneNode => Optiontime.Time
func (n noneNode) FoldTime(a time.Time, f func(v Node) time.Time) time.Time {
	return a
}
func (n noneNode) FoldTimeF(a func() time.Time, f func(v Node) time.Time) time.Time {
	return a()
}

 // map NoneNode => OptionUInt64
func (n noneNode) MapUInt64(f func(v Node) uint64) OptionUInt64 {
	return noneUInt64{}
}
// fold NoneNode => Optionuint64
func (n noneNode) FoldUInt64(a uint64, f func(v Node) uint64) uint64 {
	return a
}
func (n noneNode) FoldUInt64F(a func() uint64, f func(v Node) uint64) uint64 {
	return a()
}
 // end of nones

// some
type someNode struct {
	v Node
}

func SomeNode(v Node) someNode {
	return someNode{v}
}
// map NoneNode
func (s someNode) Map(f func(Node)) { f(s.v) }
func (s someNode) FoldF(l func(), r func(Node)) { r(s.v) }

// map SoneNode => OptionError
func (s someNode) MapError(f func(v Node) error) OptionError {
	return SomeError(f(s.v))
}
// fold SomeNode => error
func (s someNode) FoldError(a error, f func(v Node) error) error {
	return f(s.v)
}
func (s someNode) FoldErrorF(a func() error, f func(v Node) error) error {
	return f(s.v)
}

// map SoneNode => OptionNode
func (s someNode) MapNode(f func(v Node) Node) OptionNode {
	return SomeNode(f(s.v))
}
// fold SomeNode => Node
func (s someNode) FoldNode(a Node, f func(v Node) Node) Node {
	return f(s.v)
}
func (s someNode) FoldNodeF(a func() Node, f func(v Node) Node) Node {
	return f(s.v)
}

// map SoneNode => OptionSerializedPart
func (s someNode) MapSerializedPart(f func(v Node) SerializedPart) OptionSerializedPart {
	return SomeSerializedPart(f(s.v))
}
// fold SomeNode => SerializedPart
func (s someNode) FoldSerializedPart(a SerializedPart, f func(v Node) SerializedPart) SerializedPart {
	return f(s.v)
}
func (s someNode) FoldSerializedPartF(a func() SerializedPart, f func(v Node) SerializedPart) SerializedPart {
	return f(s.v)
}

// map SoneNode => OptionString
func (s someNode) MapString(f func(v Node) string) OptionString {
	return SomeString(f(s.v))
}
// fold SomeNode => string
func (s someNode) FoldString(a string, f func(v Node) string) string {
	return f(s.v)
}
func (s someNode) FoldStringF(a func() string, f func(v Node) string) string {
	return f(s.v)
}

// map SoneNode => OptionTime
func (s someNode) MapTime(f func(v Node) time.Time) OptionTime {
	return SomeTime(f(s.v))
}
// fold SomeNode => time.Time
func (s someNode) FoldTime(a time.Time, f func(v Node) time.Time) time.Time {
	return f(s.v)
}
func (s someNode) FoldTimeF(a func() time.Time, f func(v Node) time.Time) time.Time {
	return f(s.v)
}

// map SoneNode => OptionUInt64
func (s someNode) MapUInt64(f func(v Node) uint64) OptionUInt64 {
	return SomeUInt64(f(s.v))
}
// fold SomeNode => uint64
func (s someNode) FoldUInt64(a uint64, f func(v Node) uint64) uint64 {
	return f(s.v)
}
func (s someNode) FoldUInt64F(a func() uint64, f func(v Node) uint64) uint64 {
	return f(s.v)
}
 // end of somes


// none
type noneSerializedPart struct{}

func NoneSerializedPart() OptionSerializedPart {
	return noneSerializedPart{}
}

// map NoneSerializedPart
func (n noneSerializedPart) Map(f func(SerializedPart)) {}
func (n noneSerializedPart) FoldF(l func(), r func(SerializedPart)) { l() }

 // map NoneSerializedPart => OptionError
func (n noneSerializedPart) MapError(f func(v SerializedPart) error) OptionError {
	return noneError{}
}
// fold NoneSerializedPart => Optionerror
func (n noneSerializedPart) FoldError(a error, f func(v SerializedPart) error) error {
	return a
}
func (n noneSerializedPart) FoldErrorF(a func() error, f func(v SerializedPart) error) error {
	return a()
}

 // map NoneSerializedPart => OptionNode
func (n noneSerializedPart) MapNode(f func(v SerializedPart) Node) OptionNode {
	return noneNode{}
}
// fold NoneSerializedPart => OptionNode
func (n noneSerializedPart) FoldNode(a Node, f func(v SerializedPart) Node) Node {
	return a
}
func (n noneSerializedPart) FoldNodeF(a func() Node, f func(v SerializedPart) Node) Node {
	return a()
}

 // map NoneSerializedPart => OptionSerializedPart
func (n noneSerializedPart) MapSerializedPart(f func(v SerializedPart) SerializedPart) OptionSerializedPart {
	return noneSerializedPart{}
}
// fold NoneSerializedPart => OptionSerializedPart
func (n noneSerializedPart) FoldSerializedPart(a SerializedPart, f func(v SerializedPart) SerializedPart) SerializedPart {
	return a
}
func (n noneSerializedPart) FoldSerializedPartF(a func() SerializedPart, f func(v SerializedPart) SerializedPart) SerializedPart {
	return a()
}

 // map NoneSerializedPart => OptionString
func (n noneSerializedPart) MapString(f func(v SerializedPart) string) OptionString {
	return noneString{}
}
// fold NoneSerializedPart => Optionstring
func (n noneSerializedPart) FoldString(a string, f func(v SerializedPart) string) string {
	return a
}
func (n noneSerializedPart) FoldStringF(a func() string, f func(v SerializedPart) string) string {
	return a()
}

 // map NoneSerializedPart => OptionTime
func (n noneSerializedPart) MapTime(f func(v SerializedPart) time.Time) OptionTime {
	return noneTime{}
}
// fold NoneSerializedPart => Optiontime.Time
func (n noneSerializedPart) FoldTime(a time.Time, f func(v SerializedPart) time.Time) time.Time {
	return a
}
func (n noneSerializedPart) FoldTimeF(a func() time.Time, f func(v SerializedPart) time.Time) time.Time {
	return a()
}

 // map NoneSerializedPart => OptionUInt64
func (n noneSerializedPart) MapUInt64(f func(v SerializedPart) uint64) OptionUInt64 {
	return noneUInt64{}
}
// fold NoneSerializedPart => Optionuint64
func (n noneSerializedPart) FoldUInt64(a uint64, f func(v SerializedPart) uint64) uint64 {
	return a
}
func (n noneSerializedPart) FoldUInt64F(a func() uint64, f func(v SerializedPart) uint64) uint64 {
	return a()
}
 // end of nones

// some
type someSerializedPart struct {
	v SerializedPart
}

func SomeSerializedPart(v SerializedPart) someSerializedPart {
	return someSerializedPart{v}
}
// map NoneSerializedPart
func (s someSerializedPart) Map(f func(SerializedPart)) { f(s.v) }
func (s someSerializedPart) FoldF(l func(), r func(SerializedPart)) { r(s.v) }

// map SoneSerializedPart => OptionError
func (s someSerializedPart) MapError(f func(v SerializedPart) error) OptionError {
	return SomeError(f(s.v))
}
// fold SomeSerializedPart => error
func (s someSerializedPart) FoldError(a error, f func(v SerializedPart) error) error {
	return f(s.v)
}
func (s someSerializedPart) FoldErrorF(a func() error, f func(v SerializedPart) error) error {
	return f(s.v)
}

// map SoneSerializedPart => OptionNode
func (s someSerializedPart) MapNode(f func(v SerializedPart) Node) OptionNode {
	return SomeNode(f(s.v))
}
// fold SomeSerializedPart => Node
func (s someSerializedPart) FoldNode(a Node, f func(v SerializedPart) Node) Node {
	return f(s.v)
}
func (s someSerializedPart) FoldNodeF(a func() Node, f func(v SerializedPart) Node) Node {
	return f(s.v)
}

// map SoneSerializedPart => OptionSerializedPart
func (s someSerializedPart) MapSerializedPart(f func(v SerializedPart) SerializedPart) OptionSerializedPart {
	return SomeSerializedPart(f(s.v))
}
// fold SomeSerializedPart => SerializedPart
func (s someSerializedPart) FoldSerializedPart(a SerializedPart, f func(v SerializedPart) SerializedPart) SerializedPart {
	return f(s.v)
}
func (s someSerializedPart) FoldSerializedPartF(a func() SerializedPart, f func(v SerializedPart) SerializedPart) SerializedPart {
	return f(s.v)
}

// map SoneSerializedPart => OptionString
func (s someSerializedPart) MapString(f func(v SerializedPart) string) OptionString {
	return SomeString(f(s.v))
}
// fold SomeSerializedPart => string
func (s someSerializedPart) FoldString(a string, f func(v SerializedPart) string) string {
	return f(s.v)
}
func (s someSerializedPart) FoldStringF(a func() string, f func(v SerializedPart) string) string {
	return f(s.v)
}

// map SoneSerializedPart => OptionTime
func (s someSerializedPart) MapTime(f func(v SerializedPart) time.Time) OptionTime {
	return SomeTime(f(s.v))
}
// fold SomeSerializedPart => time.Time
func (s someSerializedPart) FoldTime(a time.Time, f func(v SerializedPart) time.Time) time.Time {
	return f(s.v)
}
func (s someSerializedPart) FoldTimeF(a func() time.Time, f func(v SerializedPart) time.Time) time.Time {
	return f(s.v)
}

// map SoneSerializedPart => OptionUInt64
func (s someSerializedPart) MapUInt64(f func(v SerializedPart) uint64) OptionUInt64 {
	return SomeUInt64(f(s.v))
}
// fold SomeSerializedPart => uint64
func (s someSerializedPart) FoldUInt64(a uint64, f func(v SerializedPart) uint64) uint64 {
	return f(s.v)
}
func (s someSerializedPart) FoldUInt64F(a func() uint64, f func(v SerializedPart) uint64) uint64 {
	return f(s.v)
}
 // end of somes


// none
type noneString struct{}

func NoneString() OptionString {
	return noneString{}
}

// map NoneString
func (n noneString) Map(f func(string)) {}
func (n noneString) FoldF(l func(), r func(string)) { l() }

 // map NoneString => OptionError
func (n noneString) MapError(f func(v string) error) OptionError {
	return noneError{}
}
// fold NoneString => Optionerror
func (n noneString) FoldError(a error, f func(v string) error) error {
	return a
}
func (n noneString) FoldErrorF(a func() error, f func(v string) error) error {
	return a()
}

 // map NoneString => OptionNode
func (n noneString) MapNode(f func(v string) Node) OptionNode {
	return noneNode{}
}
// fold NoneString => OptionNode
func (n noneString) FoldNode(a Node, f func(v string) Node) Node {
	return a
}
func (n noneString) FoldNodeF(a func() Node, f func(v string) Node) Node {
	return a()
}

 // map NoneString => OptionSerializedPart
func (n noneString) MapSerializedPart(f func(v string) SerializedPart) OptionSerializedPart {
	return noneSerializedPart{}
}
// fold NoneString => OptionSerializedPart
func (n noneString) FoldSerializedPart(a SerializedPart, f func(v string) SerializedPart) SerializedPart {
	return a
}
func (n noneString) FoldSerializedPartF(a func() SerializedPart, f func(v string) SerializedPart) SerializedPart {
	return a()
}

 // map NoneString => OptionString
func (n noneString) MapString(f func(v string) string) OptionString {
	return noneString{}
}
// fold NoneString => Optionstring
func (n noneString) FoldString(a string, f func(v string) string) string {
	return a
}
func (n noneString) FoldStringF(a func() string, f func(v string) string) string {
	return a()
}

 // map NoneString => OptionTime
func (n noneString) MapTime(f func(v string) time.Time) OptionTime {
	return noneTime{}
}
// fold NoneString => Optiontime.Time
func (n noneString) FoldTime(a time.Time, f func(v string) time.Time) time.Time {
	return a
}
func (n noneString) FoldTimeF(a func() time.Time, f func(v string) time.Time) time.Time {
	return a()
}

 // map NoneString => OptionUInt64
func (n noneString) MapUInt64(f func(v string) uint64) OptionUInt64 {
	return noneUInt64{}
}
// fold NoneString => Optionuint64
func (n noneString) FoldUInt64(a uint64, f func(v string) uint64) uint64 {
	return a
}
func (n noneString) FoldUInt64F(a func() uint64, f func(v string) uint64) uint64 {
	return a()
}
 // end of nones

// some
type someString struct {
	v string
}

func SomeString(v string) someString {
	return someString{v}
}
// map NoneString
func (s someString) Map(f func(string)) { f(s.v) }
func (s someString) FoldF(l func(), r func(string)) { r(s.v) }

// map SoneString => OptionError
func (s someString) MapError(f func(v string) error) OptionError {
	return SomeError(f(s.v))
}
// fold SomeString => error
func (s someString) FoldError(a error, f func(v string) error) error {
	return f(s.v)
}
func (s someString) FoldErrorF(a func() error, f func(v string) error) error {
	return f(s.v)
}

// map SoneString => OptionNode
func (s someString) MapNode(f func(v string) Node) OptionNode {
	return SomeNode(f(s.v))
}
// fold SomeString => Node
func (s someString) FoldNode(a Node, f func(v string) Node) Node {
	return f(s.v)
}
func (s someString) FoldNodeF(a func() Node, f func(v string) Node) Node {
	return f(s.v)
}

// map SoneString => OptionSerializedPart
func (s someString) MapSerializedPart(f func(v string) SerializedPart) OptionSerializedPart {
	return SomeSerializedPart(f(s.v))
}
// fold SomeString => SerializedPart
func (s someString) FoldSerializedPart(a SerializedPart, f func(v string) SerializedPart) SerializedPart {
	return f(s.v)
}
func (s someString) FoldSerializedPartF(a func() SerializedPart, f func(v string) SerializedPart) SerializedPart {
	return f(s.v)
}

// map SoneString => OptionString
func (s someString) MapString(f func(v string) string) OptionString {
	return SomeString(f(s.v))
}
// fold SomeString => string
func (s someString) FoldString(a string, f func(v string) string) string {
	return f(s.v)
}
func (s someString) FoldStringF(a func() string, f func(v string) string) string {
	return f(s.v)
}

// map SoneString => OptionTime
func (s someString) MapTime(f func(v string) time.Time) OptionTime {
	return SomeTime(f(s.v))
}
// fold SomeString => time.Time
func (s someString) FoldTime(a time.Time, f func(v string) time.Time) time.Time {
	return f(s.v)
}
func (s someString) FoldTimeF(a func() time.Time, f func(v string) time.Time) time.Time {
	return f(s.v)
}

// map SoneString => OptionUInt64
func (s someString) MapUInt64(f func(v string) uint64) OptionUInt64 {
	return SomeUInt64(f(s.v))
}
// fold SomeString => uint64
func (s someString) FoldUInt64(a uint64, f func(v string) uint64) uint64 {
	return f(s.v)
}
func (s someString) FoldUInt64F(a func() uint64, f func(v string) uint64) uint64 {
	return f(s.v)
}
 // end of somes


// none
type noneTime struct{}

func NoneTime() OptionTime {
	return noneTime{}
}

// map NoneTime
func (n noneTime) Map(f func(time.Time)) {}
func (n noneTime) FoldF(l func(), r func(time.Time)) { l() }

 // map NoneTime => OptionError
func (n noneTime) MapError(f func(v time.Time) error) OptionError {
	return noneError{}
}
// fold NoneTime => Optionerror
func (n noneTime) FoldError(a error, f func(v time.Time) error) error {
	return a
}
func (n noneTime) FoldErrorF(a func() error, f func(v time.Time) error) error {
	return a()
}

 // map NoneTime => OptionNode
func (n noneTime) MapNode(f func(v time.Time) Node) OptionNode {
	return noneNode{}
}
// fold NoneTime => OptionNode
func (n noneTime) FoldNode(a Node, f func(v time.Time) Node) Node {
	return a
}
func (n noneTime) FoldNodeF(a func() Node, f func(v time.Time) Node) Node {
	return a()
}

 // map NoneTime => OptionSerializedPart
func (n noneTime) MapSerializedPart(f func(v time.Time) SerializedPart) OptionSerializedPart {
	return noneSerializedPart{}
}
// fold NoneTime => OptionSerializedPart
func (n noneTime) FoldSerializedPart(a SerializedPart, f func(v time.Time) SerializedPart) SerializedPart {
	return a
}
func (n noneTime) FoldSerializedPartF(a func() SerializedPart, f func(v time.Time) SerializedPart) SerializedPart {
	return a()
}

 // map NoneTime => OptionString
func (n noneTime) MapString(f func(v time.Time) string) OptionString {
	return noneString{}
}
// fold NoneTime => Optionstring
func (n noneTime) FoldString(a string, f func(v time.Time) string) string {
	return a
}
func (n noneTime) FoldStringF(a func() string, f func(v time.Time) string) string {
	return a()
}

 // map NoneTime => OptionTime
func (n noneTime) MapTime(f func(v time.Time) time.Time) OptionTime {
	return noneTime{}
}
// fold NoneTime => Optiontime.Time
func (n noneTime) FoldTime(a time.Time, f func(v time.Time) time.Time) time.Time {
	return a
}
func (n noneTime) FoldTimeF(a func() time.Time, f func(v time.Time) time.Time) time.Time {
	return a()
}

 // map NoneTime => OptionUInt64
func (n noneTime) MapUInt64(f func(v time.Time) uint64) OptionUInt64 {
	return noneUInt64{}
}
// fold NoneTime => Optionuint64
func (n noneTime) FoldUInt64(a uint64, f func(v time.Time) uint64) uint64 {
	return a
}
func (n noneTime) FoldUInt64F(a func() uint64, f func(v time.Time) uint64) uint64 {
	return a()
}
 // end of nones

// some
type someTime struct {
	v time.Time
}

func SomeTime(v time.Time) someTime {
	return someTime{v}
}
// map NoneTime
func (s someTime) Map(f func(time.Time)) { f(s.v) }
func (s someTime) FoldF(l func(), r func(time.Time)) { r(s.v) }

// map SoneTime => OptionError
func (s someTime) MapError(f func(v time.Time) error) OptionError {
	return SomeError(f(s.v))
}
// fold SomeTime => error
func (s someTime) FoldError(a error, f func(v time.Time) error) error {
	return f(s.v)
}
func (s someTime) FoldErrorF(a func() error, f func(v time.Time) error) error {
	return f(s.v)
}

// map SoneTime => OptionNode
func (s someTime) MapNode(f func(v time.Time) Node) OptionNode {
	return SomeNode(f(s.v))
}
// fold SomeTime => Node
func (s someTime) FoldNode(a Node, f func(v time.Time) Node) Node {
	return f(s.v)
}
func (s someTime) FoldNodeF(a func() Node, f func(v time.Time) Node) Node {
	return f(s.v)
}

// map SoneTime => OptionSerializedPart
func (s someTime) MapSerializedPart(f func(v time.Time) SerializedPart) OptionSerializedPart {
	return SomeSerializedPart(f(s.v))
}
// fold SomeTime => SerializedPart
func (s someTime) FoldSerializedPart(a SerializedPart, f func(v time.Time) SerializedPart) SerializedPart {
	return f(s.v)
}
func (s someTime) FoldSerializedPartF(a func() SerializedPart, f func(v time.Time) SerializedPart) SerializedPart {
	return f(s.v)
}

// map SoneTime => OptionString
func (s someTime) MapString(f func(v time.Time) string) OptionString {
	return SomeString(f(s.v))
}
// fold SomeTime => string
func (s someTime) FoldString(a string, f func(v time.Time) string) string {
	return f(s.v)
}
func (s someTime) FoldStringF(a func() string, f func(v time.Time) string) string {
	return f(s.v)
}

// map SoneTime => OptionTime
func (s someTime) MapTime(f func(v time.Time) time.Time) OptionTime {
	return SomeTime(f(s.v))
}
// fold SomeTime => time.Time
func (s someTime) FoldTime(a time.Time, f func(v time.Time) time.Time) time.Time {
	return f(s.v)
}
func (s someTime) FoldTimeF(a func() time.Time, f func(v time.Time) time.Time) time.Time {
	return f(s.v)
}

// map SoneTime => OptionUInt64
func (s someTime) MapUInt64(f func(v time.Time) uint64) OptionUInt64 {
	return SomeUInt64(f(s.v))
}
// fold SomeTime => uint64
func (s someTime) FoldUInt64(a uint64, f func(v time.Time) uint64) uint64 {
	return f(s.v)
}
func (s someTime) FoldUInt64F(a func() uint64, f func(v time.Time) uint64) uint64 {
	return f(s.v)
}
 // end of somes


// none
type noneUInt64 struct{}

func NoneUInt64() OptionUInt64 {
	return noneUInt64{}
}

// map NoneUInt64
func (n noneUInt64) Map(f func(uint64)) {}
func (n noneUInt64) FoldF(l func(), r func(uint64)) { l() }

 // map NoneUInt64 => OptionError
func (n noneUInt64) MapError(f func(v uint64) error) OptionError {
	return noneError{}
}
// fold NoneUInt64 => Optionerror
func (n noneUInt64) FoldError(a error, f func(v uint64) error) error {
	return a
}
func (n noneUInt64) FoldErrorF(a func() error, f func(v uint64) error) error {
	return a()
}

 // map NoneUInt64 => OptionNode
func (n noneUInt64) MapNode(f func(v uint64) Node) OptionNode {
	return noneNode{}
}
// fold NoneUInt64 => OptionNode
func (n noneUInt64) FoldNode(a Node, f func(v uint64) Node) Node {
	return a
}
func (n noneUInt64) FoldNodeF(a func() Node, f func(v uint64) Node) Node {
	return a()
}

 // map NoneUInt64 => OptionSerializedPart
func (n noneUInt64) MapSerializedPart(f func(v uint64) SerializedPart) OptionSerializedPart {
	return noneSerializedPart{}
}
// fold NoneUInt64 => OptionSerializedPart
func (n noneUInt64) FoldSerializedPart(a SerializedPart, f func(v uint64) SerializedPart) SerializedPart {
	return a
}
func (n noneUInt64) FoldSerializedPartF(a func() SerializedPart, f func(v uint64) SerializedPart) SerializedPart {
	return a()
}

 // map NoneUInt64 => OptionString
func (n noneUInt64) MapString(f func(v uint64) string) OptionString {
	return noneString{}
}
// fold NoneUInt64 => Optionstring
func (n noneUInt64) FoldString(a string, f func(v uint64) string) string {
	return a
}
func (n noneUInt64) FoldStringF(a func() string, f func(v uint64) string) string {
	return a()
}

 // map NoneUInt64 => OptionTime
func (n noneUInt64) MapTime(f func(v uint64) time.Time) OptionTime {
	return noneTime{}
}
// fold NoneUInt64 => Optiontime.Time
func (n noneUInt64) FoldTime(a time.Time, f func(v uint64) time.Time) time.Time {
	return a
}
func (n noneUInt64) FoldTimeF(a func() time.Time, f func(v uint64) time.Time) time.Time {
	return a()
}

 // map NoneUInt64 => OptionUInt64
func (n noneUInt64) MapUInt64(f func(v uint64) uint64) OptionUInt64 {
	return noneUInt64{}
}
// fold NoneUInt64 => Optionuint64
func (n noneUInt64) FoldUInt64(a uint64, f func(v uint64) uint64) uint64 {
	return a
}
func (n noneUInt64) FoldUInt64F(a func() uint64, f func(v uint64) uint64) uint64 {
	return a()
}
 // end of nones

// some
type someUInt64 struct {
	v uint64
}

func SomeUInt64(v uint64) someUInt64 {
	return someUInt64{v}
}
// map NoneUInt64
func (s someUInt64) Map(f func(uint64)) { f(s.v) }
func (s someUInt64) FoldF(l func(), r func(uint64)) { r(s.v) }

// map SoneUInt64 => OptionError
func (s someUInt64) MapError(f func(v uint64) error) OptionError {
	return SomeError(f(s.v))
}
// fold SomeUInt64 => error
func (s someUInt64) FoldError(a error, f func(v uint64) error) error {
	return f(s.v)
}
func (s someUInt64) FoldErrorF(a func() error, f func(v uint64) error) error {
	return f(s.v)
}

// map SoneUInt64 => OptionNode
func (s someUInt64) MapNode(f func(v uint64) Node) OptionNode {
	return SomeNode(f(s.v))
}
// fold SomeUInt64 => Node
func (s someUInt64) FoldNode(a Node, f func(v uint64) Node) Node {
	return f(s.v)
}
func (s someUInt64) FoldNodeF(a func() Node, f func(v uint64) Node) Node {
	return f(s.v)
}

// map SoneUInt64 => OptionSerializedPart
func (s someUInt64) MapSerializedPart(f func(v uint64) SerializedPart) OptionSerializedPart {
	return SomeSerializedPart(f(s.v))
}
// fold SomeUInt64 => SerializedPart
func (s someUInt64) FoldSerializedPart(a SerializedPart, f func(v uint64) SerializedPart) SerializedPart {
	return f(s.v)
}
func (s someUInt64) FoldSerializedPartF(a func() SerializedPart, f func(v uint64) SerializedPart) SerializedPart {
	return f(s.v)
}

// map SoneUInt64 => OptionString
func (s someUInt64) MapString(f func(v uint64) string) OptionString {
	return SomeString(f(s.v))
}
// fold SomeUInt64 => string
func (s someUInt64) FoldString(a string, f func(v uint64) string) string {
	return f(s.v)
}
func (s someUInt64) FoldStringF(a func() string, f func(v uint64) string) string {
	return f(s.v)
}

// map SoneUInt64 => OptionTime
func (s someUInt64) MapTime(f func(v uint64) time.Time) OptionTime {
	return SomeTime(f(s.v))
}
// fold SomeUInt64 => time.Time
func (s someUInt64) FoldTime(a time.Time, f func(v uint64) time.Time) time.Time {
	return f(s.v)
}
func (s someUInt64) FoldTimeF(a func() time.Time, f func(v uint64) time.Time) time.Time {
	return f(s.v)
}

// map SoneUInt64 => OptionUInt64
func (s someUInt64) MapUInt64(f func(v uint64) uint64) OptionUInt64 {
	return SomeUInt64(f(s.v))
}
// fold SomeUInt64 => uint64
func (s someUInt64) FoldUInt64(a uint64, f func(v uint64) uint64) uint64 {
	return f(s.v)
}
func (s someUInt64) FoldUInt64F(a func() uint64, f func(v uint64) uint64) uint64 {
	return f(s.v)
}
 // end of somes

 // end of everything

