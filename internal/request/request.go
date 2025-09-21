package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
	Status      int
	buffer      []byte
}

func (r *Request) parse(data []byte) (int, error) {
	r.buffer = append(r.buffer, data...)

	req, n, err := parseRequestLine(r.buffer)

	if n != 0 {
		r.RequestLine = req
		r.Status = 1

		return n, err
	}

	return 0, nil
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func IsLetter(s string) bool {
	return !strings.ContainsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r)
	})
}

func parseRequestLine(input []byte) (RequestLine, int, error) {
	fullRequest := string(input)

	if !strings.Contains(fullRequest, "\r\n") {
		return RequestLine{}, 0, nil
	}

	requestLine := strings.Split(fullRequest, "\r\n")[0]

	requestParts := strings.Split(requestLine, " ")

	if len(requestParts) != 3 {
		return RequestLine{}, len(requestLine), fmt.Errorf("request doesn't have 3 parts")
	}

	method := requestParts[0]

	if method != strings.ToUpper(method) && !IsLetter(method) {
		return RequestLine{}, len(requestLine), fmt.Errorf("method must be in all caps and only contain alphabetic characters")
	}

	httpVersion := strings.Split(requestParts[2], "/")[1]

	if httpVersion != "1.1" {
		return RequestLine{}, len(requestLine), fmt.Errorf("only HTTP/1.1 methods allowed")
	}

	return RequestLine{
		HttpVersion:   httpVersion,
		RequestTarget: requestParts[1],
		Method:        method,
	}, len(requestLine), nil

}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req := Request{
		Status: 0,
	}

	for {
		data := make([]byte, 8)
		n, err := reader.Read(data)

		if err != nil || errors.Is(err, io.EOF) {
			return &Request{}, fmt.Errorf("error when reading from input: %s", err)
		}

		_, err = req.parse(data[:n])

		if err != nil {
			return &Request{}, fmt.Errorf("request error: %s", err)
		}

		if req.Status == 1 {
			return &req, nil
		}
	}
}
