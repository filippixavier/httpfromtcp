package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")

	return h
}

type Writer struct {
	Writer io.Writer
	state  responseState
}

type responseState int

const (
	stateWriteStatusLine = iota
	stateWriteHeaders
	stateWriteBody
)

func NewWriter(w io.Writer) Writer {
	return Writer{
		Writer: w,
		state:  stateWriteStatusLine,
	}
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != stateWriteHeaders {
		return fmt.Errorf("must call WriteStatusLines first")
	}

	for name, value := range headers {
		_, err := fmt.Fprintf(w.Writer, "%s: %s\r\n", name, value)
		if err != nil {
			return err
		}
	}

	_, err := fmt.Fprintf(w.Writer, "\r\n")

	if err == nil {
		w.state = stateWriteBody
	}

	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != stateWriteBody {
		return 0, fmt.Errorf("must call WriteHeaders first")
	}

	return w.Writer.Write(p)
}
