package main

import (
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func test(w io.Writer, req *request.Request) *server.HandlerError {
	var h *server.HandlerError

	if req.RequestLine.RequestTarget == "/yourproblem" {
		h = &server.HandlerError{
			StatusCode: response.BadRequest,
			ErrorMsg:   []byte("Your problem is not my problem\n"),
		}
	} else if req.RequestLine.RequestTarget == "/myproblem" {
		h = &server.HandlerError{
			StatusCode: response.InternalError,
			ErrorMsg:   []byte("Woopsie, my bad\n"),
		}
	} else {
		w.Write([]byte("All good, frfr\n"))
	}

	return h
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
