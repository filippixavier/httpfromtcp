package request

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
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

func parseRequestLine(input []byte) (RequestLine, error) {
	fullRequest := string(input)
	requestLine := strings.Split(fullRequest, "\r\n")[0]

	requestParts := strings.Split(requestLine, " ")

	if len(requestParts) != 3 {
		return RequestLine{}, fmt.Errorf("request doesn't have 3 parts")
	}

	method := requestParts[0]

	if method != strings.ToUpper(method) && !IsLetter(method) {
		return RequestLine{}, fmt.Errorf("method must be in all caps and only contain alphabetic characters")
	}

	httpVersion := strings.Split(requestParts[2], "/")[1]

	if httpVersion != "1.1" {
		return RequestLine{}, fmt.Errorf("only HTTP/1.1 methods allowed")
	}

	return RequestLine{
		HttpVersion:   httpVersion,
		RequestTarget: requestParts[1],
		Method:        method,
	}, nil

}

func RequestFromReader(reader io.Reader) (*Request, error) {
	input, err := io.ReadAll(reader)

	req := Request{}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading line: %v\n", err)
		return nil, err
	}

	requestLine, err := parseRequestLine(input)

	if err != nil {
		return nil, err
	}

	req.RequestLine = requestLine

	return &req, nil
}
