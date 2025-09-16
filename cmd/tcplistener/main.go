package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const port = ":42069"

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)

	go func() {
		b := make([]byte, 8)
		currentLineContents := ""

		defer f.Close()
		defer close(lines)

		for {
			n, err := f.Read(b)

			if err != nil {
				if currentLineContents != "" {
					lines <- currentLineContents
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				break
			}

			str := string(b[:n])
			parts := strings.Split(str, "\n")

			for i := 0; i < len(parts)-1; i++ {
				lines <- fmt.Sprintf("%s%s", currentLineContents, parts[i])
				currentLineContents = ""
			}

			currentLineContents += parts[len(parts)-1]
		}
	}()

	return lines
}

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

		ch := getLinesChannel(connection)

		for str := range ch {
			fmt.Println(str)
		}

		fmt.Println("Connection to ", connection.RemoteAddr(), "closed")
	}
}
