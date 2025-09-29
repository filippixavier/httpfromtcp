package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"log"
	"net"
)

const port = ":42069"

func main() {
	listener, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("could not open %s: %s\n", port, err)
	}

	fmt.Println("Listening for TCP traffic on", port)
	defer listener.Close()

	for {
		connection, err := listener.Accept()

		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}

		fmt.Println("Connection accepted")

		req, err := request.RequestFromReader(connection)

		if err != nil {
			log.Fatalf("error when reading from reader: %s", err)
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %v\n", req.RequestLine.Method)
		fmt.Printf("- Target: %v\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %v\n", req.RequestLine.HttpVersion)
		fmt.Printf("- Request Line: %v\n", req.RequestLine.HttpVersion)
		fmt.Println("Headers:")

		for key := range req.Headers {
			fmt.Printf("- %v: %v\n", key, req.Headers[key])
		}

		fmt.Println("Body:")
		fmt.Printf("%s\n", req.Body)

		fmt.Println("Connection to ", connection.RemoteAddr(), "closed")
	}
}
