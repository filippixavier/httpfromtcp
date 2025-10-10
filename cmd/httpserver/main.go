package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"log"
	"os"
	"os/signal"
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

func test(w response.Writer, req *request.Request) {
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
	server, err := server.Serve(port, test)
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
