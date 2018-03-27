package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"strings"

	"github.com/labstack/echo"
)

const (
	HeaderContentTranferEncoding = "Content-Transfer-Encoding"
	HeaderContentType            = "Content-Type"
)

func mainContentType(p *multipart.Part) string {
	ct := p.Header.Get("Content-Type")
	ctParts := strings.Split(ct, "/")
	return ctParts[0]
}

type SerializedPart interface {
	ContentType() string
	ContentEncoding() string
	Content() []byte
	DecodedContent() []byte
	ContentString() string
	FileName() string
	MainType() string
	SubType() string
	Walk(func(SerializedPart))
}

type serPart struct {
	contentType     string
	contentEncoding string
	content         []byte
	filename        string
	children        []SerializedPart
}

type SerializedMessage interface {
	Get(string) string
	Root() SerializedPart
	Parse() SerializedMessage
}

type serMsg struct {
	orig *mail.Message
	root *serPart
}

func MakeSerializedMsg(orig *[]byte) ResultSerializedMessage {
	reader := bytes.NewReader(*orig)
	return ResultMessageFrom(mail.ReadMessage(reader)).
		MapSerializedMessage(
			func(msg *mail.Message) SerializedMessage {
				return &serMsg{
					orig: msg,
					root: &serPart{},
				}
			})
}

func (s *serMsg) Root() SerializedPart {
	return s.root
}

func (s *serMsg) Get(key string) string {
	if s.orig != nil {
		return s.orig.Header.Get(key)
	}
	return ""
}

func (s *serMsg) Parse() SerializedMessage {
	if s.orig == nil {
		return s
	}
	cte := s.orig.Header.Get(HeaderContentTranferEncoding)
	mediatype, params, err := mime.ParseMediaType(
		s.orig.Header.Get(HeaderContentType))

	if err != nil {
		return s
	}

	s.root = &serPart{
		contentType:     mediatype,
		contentEncoding: cte,
	}

	if isMultipart(mediatype) {
		boundary := params["boundary"]
		walkPartSync(s.orig.Body, boundary, &s.root.children)
	} else {
		ResultSByteFrom(ioutil.ReadAll(s.orig.Body)).
			Map(func(body []byte) {
				s.root.content = body
			})
	}

	return s
}

func (s *serPart) ContentType() string {
	return s.contentType
}
func (s *serPart) ContentEncoding() string {
	return s.contentEncoding
}
func (s *serPart) Content() []byte {
	return s.content
}
func (s *serPart) DecodedContent() []byte {
	return decodeContent(s.content, s.contentEncoding)
}

func (s *serPart) ContentString() string {
	return fmt.Sprintf("%s", decodeContent(s.content, s.contentEncoding))
}
func (s *serPart) FileName() string {
	return s.filename
}

func (s *serPart) MainType() string {
	mediaType, _, err := mime.ParseMediaType(s.contentType)
	if err != nil {
		return "any"
	}
	ctParts := strings.Split(mediaType, "/")
	return ctParts[0]
}

func (s *serPart) SubType() string {
	mediaType, _, err := mime.ParseMediaType(s.contentType)
	if err != nil {
		return "any"
	}
	ctParts := strings.Split(mediaType, "/")
	return ctParts[1]
}

// func walker(s *serPart, c chan SerializedPart, depth int) {
// 	c <- s
// 	for _, p := range s.children {
// 		sp := p.(*serPart)
// 		walker(sp, c, depth+1)
// 	}
// 	if 0 == depth {
// 		c <- nil
// 	}
// }

func (s *serPart) Walk(f func(SerializedPart)) {
	f(s)
	for _, p := range s.children {
		p.Walk(f)
	}
}

func serPartF(p *multipart.Part) serPart {
	ct := p.Header.Get(echo.HeaderContentType)
	cte := p.Header.Get(HeaderContentTranferEncoding)
	fn := p.FileName()
	var ret = serPart{}
	log.Printf("serPartF %s %s", ct, fn)
	input, err := ioutil.ReadAll(p)
	if err != nil {
		return ret
	}

	return serPart{
		contentType:     ct,
		contentEncoding: cte,
		content:         input,
		filename:        fn,
	}
}

func walkParts(r io.Reader, boundary string, c chan SerializedPart, depth int) {
	log.Printf("walkParts %v", depth)
	// reader := multipart.NewReader(r, boundary)
	// counter := 0
	// for {
	// 	counter = counter + 1
	// 	part, err := reader.NextPart()
	// 	if err != nil {
	// 		log.Printf("Err on NextPart %s", err.Error())
	// 		break
	// 	}

	// 	contentType := part.Header.Get("Content-Type")
	// 	mediaType, params, err := mime.ParseMediaType(contentType)
	// 	if err != nil {
	// 		log.Printf("Err on ParseMediaType %s", err.Error())
	// 		break
	// 	}

	// 	if !isMultipart(mediaType) {
	// 		c <- serPartF(part)
	// 	} else {
	// 		newBoundary := params["boundary"]
	// 		walkParts(part, newBoundary, c, depth+1)
	// 	}
	// 	part.Close()
	// }

	// if 0 == depth {
	// 	c <- nil
	// }
}

func walkPartSync(r io.Reader, boundary string, parts *[]SerializedPart) {
	// log.Printf("walkPartSync %v", r)
	reader := multipart.NewReader(r, boundary)
	for {
		part, err := reader.NextPart()
		if err != nil {
			log.Printf("Err on NextPart %s", err.Error())
			break
		}

		contentType := part.Header.Get("Content-Type")
		mediaType, params, err := mime.ParseMediaType(contentType)
		if err != nil {
			log.Printf("Err on ParseMediaType %s", err.Error())
			break
		}

		sp := serPartF(part)
		*parts = append(*parts, &sp)
		if isMultipart(mediaType) {
			newBoundary := params["boundary"]
			walkPartSync(part, newBoundary, &sp.children)
		}
		part.Close()
	}
}

func decodeContent(input []byte, cte string) []byte {
	log.Printf("decodeContent (%s)", cte)
	switch cte {

	case "base64":
		{
			content, err := base64.StdEncoding.DecodeString(string(input))
			if err != nil {
				return nil
			}
			return content
		}

	case "7bit":
		{
			r := quotedprintable.NewReader(bytes.NewReader(input))
			content, err := ioutil.ReadAll(r)
			if err != nil {
				return nil
			}
			return content
		}
	}

	return input
}

func isMultipart(mediaType string) bool {
	const rfc822 = "message/rfc822"
	const mp = "multipart"
	mainType := strings.Split(mediaType, "/")[0]
	if mediaType == rfc822 || mainType == mp {
		return true
	}
	return false
}

type attachment struct {
	r         io.Reader
	name      string
	mediaType string
}

func errAttach(r io.Reader) attachment {
	return attachment{
		r:         r,
		name:      "Error.txt",
		mediaType: "text/plain",
	}
}

// func SerializeMessage(
// 	data *[]byte, paramSender string, paramTopic string, paramId int) SerializedMessage {
// 	sm := newMsg(data)
// 	sm.Parse()
// 	return sm
// }

// func getAttachment(name string, data *[]byte) attachment {
// 	reader := bytes.NewReader(*data)
// 	msg, err := mail.ReadMessage(reader)
// 	if err != nil {
// 		return errAttach(sReader("Error Reading Message: ", err.Error()))
// 	}
// 	_, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
// 	if err != nil {
// 		return errAttach(sReader("Error MediaType: ", err.Error()))
// 	}
// 	mc := make(chan SerializedPart)
// 	boundary := params["boundary"]
// 	go walkParts(msg.Body, boundary, mc, 0)

// 	for {
// 		part := <-mc
// 		if part == nil {
// 			return errAttach(sReader("Attachment Not Found ", name))
// 		}
// 		mp := part.MainType()
// 		fn := part.FileName()
// 		content := decodeContent(part.Content(), part.ContentEncoding())
// 		if mp != "text" && fn == name {
// 			return attachment{
// 				r:         bytes.NewReader(content),
// 				name:      fn,
// 				mediaType: part.ContentType(),
// 			}
// 		}
// 	}

// }
