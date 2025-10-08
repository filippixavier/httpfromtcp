package response

import (
	"fmt"
	"io"
)

type StatusCode int

const (
	Ok            StatusCode = 200
	BadRequest    StatusCode = 400
	InternalError StatusCode = 500
)

func getStatusLine(statusCode StatusCode) []byte {
	reason := ""
	switch statusCode {
	case Ok:
		reason = "OK"
	case BadRequest:
		reason = "Bad Request"
	case InternalError:
		reason = "Internal Server Error"
	}

	return []byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reason))
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	_, err := w.Write(getStatusLine(statusCode))

	return err
}
