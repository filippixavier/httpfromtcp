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

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	// Teacher use a for ... range loop to access content, which *should* work for small maps, but according to this
	// https://medium.com/@AlexanderObregon/go-map-internals-and-why-ordering-isnt-stable-69551a7582c8
	// article, maps doesn't have a stable order.
	// Then again, headers order doesn't matter so...
	_, err := fmt.Fprintf(w, "Content-Length: %s\r\n", headers.Get("Content-Length"))
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "Connection: %s\r\n", headers.Get("Connection"))
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "Content-Type: %s\r\n", headers.Get("Content-Type"))

	fmt.Fprint(w, "\r\n")

	return err
}
