package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const port = 42069

const badRequest = `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`

const internalError = `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`

const ok = `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`

func handleProxy(w response.Writer, req *request.Request) {
	destination := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	headers := response.GetDefaultHeaders(0)

	resp, err := http.Get(fmt.Sprintf("https://httpbin.org/%s", destination))

	if err != nil {
		handleError(w, req, err)
	}

	defer resp.Body.Close()

	headers.Delete("Content-Length")
	headers.Set("Transfer-Encoding", "chunked")
	headers.Set("Trailer", "X-Content-SHA256,X-Content-Length")

	buf := make([]byte, 1024)
	readToIndex := 0

	err = w.WriteStatusLine(response.Ok)

	if err != nil {
		handleError(w, req, err)
		return
	}

	err = w.WriteHeaders(headers)

	if err != nil {
		handleError(w, req, err)
		return
	}

	for {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		numBytesRead, err := resp.Body.Read(buf[readToIndex:])

		if err != nil {
			if errors.Is(err, io.EOF) {
				w.WriteChunkedBodyDone()
				break
			}
			handleError(w, req, err)
			return
		}

		readToIndex += numBytesRead

		_, err = w.WriteChunkedBody(buf[:readToIndex])

		if err != nil {
			fmt.Println(err)
			handleError(w, req, err)
			return
		}
	}

	checksum := sha256.Sum256(buf[:readToIndex])

	headers.Set("X-Content-SHA256", fmt.Sprintf("%x", checksum))
	headers.Set("X-Content-Length", fmt.Sprintf("%d", len(buf[:readToIndex])))

	err = w.WriteTrailers(headers)

	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

func handleError(w response.Writer, _ *request.Request, e error) {
	headers := response.GetDefaultHeaders(0)
	errorString := fmt.Sprintf("%e", e)
	headers.Override("Content-Length", fmt.Sprintf("%d", len(errorString)))

	err := w.WriteStatusLine(response.InternalError)
	if err != nil {
		fmt.Println(err)
	}

	err = w.WriteHeaders(headers)

	if err != nil {
		fmt.Println(err)
	}

	_, err = w.WriteBody([]byte(errorString))

	if err != nil {
		fmt.Println(err)
	}
}

func handler(w response.Writer, req *request.Request) {
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		handleProxy(w, req)
		return
	}

	headers := response.GetDefaultHeaders(0)
	headers.Override("Content-Type", "text/html")

	var body []byte
	var status response.StatusCode

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		status = response.BadRequest
		body = []byte(badRequest)
	case "/myproblem":
		status = response.InternalError
		body = []byte(internalError)
	default:
		status = response.Ok
		body = []byte(ok)
	}

	headers.Override("Content-Length", fmt.Sprintf("%d", len(body)))

	err := w.WriteStatusLine(status)
	if err != nil {
		fmt.Println(err)
	}

	err = w.WriteHeaders(headers)

	if err != nil {
		fmt.Println(err)
	}

	_, err = w.WriteBody(body)

	if err != nil {
		fmt.Println(err)
	}

}

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
