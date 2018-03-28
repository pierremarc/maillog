package main

import (
	"bytes"
	"encoding/base64"
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
	MediaType() string
	MainType() string
	SubType() string
	ContentEncoding() string
	Content() []byte
	DecodedContent() []byte
	ContentString() string
	FileName() string
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

func (s *serPart) ContentType() string {
	return s.contentType
}
func (s *serPart) MainType() string {
	mediaType, _, err := mime.ParseMediaType(s.contentType)
	if err != nil {
		return "*"
	}
	ctParts := strings.Split(mediaType, "/")
	return ctParts[0]
}

func (s *serPart) MediaType() string {
	mediaType, _, err := mime.ParseMediaType(s.contentType)
	if err != nil {
		return "*/*"
	}
	return mediaType
}

func (s *serPart) SubType() string {
	mediaType, _, err := mime.ParseMediaType(s.contentType)
	if err != nil {
		return "*"
	}
	ctParts := strings.Split(mediaType, "/")
	return ctParts[1]
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
	// return fmt.Sprintf("%s", decodeContent(s.content, s.contentEncoding))
	return string(decodeContent(s.content, s.contentEncoding))
}
func (s *serPart) FileName() string {
	return s.filename
}

func (s *serPart) Walk(f func(SerializedPart)) {
	// log.Printf("serPart.Walk() %v", s)
	f(s)
	for _, p := range s.children {
		p.Walk(f)
	}
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

		contentType := part.Header.Get(HeaderContentType)
		// cte := part.Header.Get(HeaderContentTranferEncoding)
		mediatype, params, err := mime.ParseMediaType(contentType)
		if err != nil {
			log.Printf("Err on ParseMediaType %s", err.Error())
			break
		}

		sp := serPartF(part)
		*parts = append(*parts, &sp)
		if isMultipart(mediatype) {
			newBoundary := params["boundary"]
			walkPartSync(part, newBoundary, &sp.children)
		}
		// } else {
		// 	bs, _ := ioutil.ReadAll(part)
		// 	log.Println("BS", bs)
		// 	ResultSByteFrom(ioutil.ReadAll(part)).
		// 		FoldF(
		// 			func(error) { log.Printf("Failed to read body (%s) (%s)", mediatype, cte) },
		// 			func(body []byte) {
		// 				log.Printf("Success to read body (%s) (%s)\n<<\n%s\n>>", mediatype, cte, body)
		// 				sp.content = body
		// 			})
		// }
		part.Close()
	}
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
	log.Println("Parse")
	if s.orig == nil {
		return s
	}
	cte := s.orig.Header.Get(HeaderContentTranferEncoding)
	mediatype, params, err := mime.ParseMediaType(
		s.orig.Header.Get(HeaderContentType))

	if err != nil {
		log.Printf("Failed ParseMediaType (%v)", err)
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
			FoldF(
				func(error) {
					log.Printf("Failed to read body (%s) (%s)", mediatype, cte)
				},
				func(body []byte) {
					log.Printf("Success to read body (%s) (%s)", mediatype, cte)
					s.root.content = body
				})
	}

	return s
}

func serPartF(p *multipart.Part) serPart {
	ct := p.Header.Get(echo.HeaderContentType)
	cte := p.Header.Get(HeaderContentTranferEncoding)
	fn := p.FileName()
	var ret = serPart{}
	log.Printf("serPartF %s %s", ct, fn)
	input, err := ioutil.ReadAll(p)
	if err != nil {
		log.Printf("Failed to read part %v", err)
		return ret
	}

	return serPart{
		contentType:     ct,
		contentEncoding: cte,
		content:         input,
		filename:        fn,
	}
}

func decodeContent(input []byte, cte string) []byte {
	log.Printf("decodeContent (%s)", cte)
	switch cte {

	case "base64":
		{
			content, err := base64.StdEncoding.DecodeString(string(input))
			if err != nil {
				log.Printf("Error:base64.StdEncoding.DecodeString (%s)", err.Error())
				// return nil
				return input
			}
			return content
		}

	case "7bit":
		{
			r := quotedprintable.NewReader(bytes.NewReader(input))
			content, err := ioutil.ReadAll(r)
			if err != nil {
				log.Printf("Error:quotedprintable.NewReader (%s)", err.Error())
				// return nil
				return input
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

// type attachment struct {
// 	r         io.Reader
// 	name      string
// 	mediaType string
// }

// func errAttach(r io.Reader) attachment {
// 	return attachment{
// 		r:         r,
// 		name:      "Error.txt",
// 		mediaType: "text/plain",
// 	}
// }
