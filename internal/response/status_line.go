package response

import (
	"fmt"
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

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != stateWriteStatusLine {
		return fmt.Errorf("cannot call WriteStatusline a second time")
	}
	_, err := w.Writer.Write(getStatusLine(statusCode))

	if err == nil {
		w.state = stateWriteHeaders
	}

	return err
}
