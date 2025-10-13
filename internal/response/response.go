package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strings"
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
	stateWriteTrailers
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

	n, err := w.Writer.Write(p)
	if err == nil {
		w.state = stateWriteTrailers
	}

	return n, err
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.state != stateWriteBody {
		return 0, fmt.Errorf("must call WriteHeaders first")
	}

	// Remember! The chunk length is in HEXADECIMAL!
	_, err := fmt.Fprintf(w.Writer, "%x\r\n", len(p))

	if err != nil {
		return 0, err
	}

	numBytesSent, err := fmt.Fprintf(w.Writer, "%s\r\n", p)

	if err != nil {
		return 0, err
	}

	return numBytesSent - 2, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.state != stateWriteBody {
		return 0, fmt.Errorf("must call WriteHeaders first")
	}
	w.state = stateWriteTrailers
	return w.Writer.Write([]byte("0\r\n"))
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	defer w.Writer.Write([]byte("\r\n"))
	trailers := h.Get("Trailer")

	splittedTrailers := strings.Split(trailers, ",")

	for _, trailer := range splittedTrailers {
		value := h.Get(trailer)
		_, err := fmt.Fprintf(w.Writer, "%s: %s\r\n", trailer, value)

		if err != nil {
			return err
		}
	}

	return nil
}
