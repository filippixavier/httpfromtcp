package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
)

type StatusCode int

const (
	Ok            StatusCode = 200
	BadRequest    StatusCode = 400
	InternalError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	reason := ""
	switch statusCode {
	case Ok:
		reason = "HTTP/1.1 200 OK\r\n"
	case BadRequest:
		reason = "HTTP/1.1 400 Bad Request\r\n"
	case InternalError:
		reason = "HTTP/1.1 500 Internal Server Error\r\n"
	}

	_, err := w.Write([]byte(reason))

	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")

	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	_, err := fmt.Fprintf(w, "Content-Length: %s\r\n", headers.Get("Content-Length"))
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "Connection: %s\r\n", headers.Get("Connection"))
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "Content-Type: %s\r\n", headers.Get("Content-Type"))

	return err
}
