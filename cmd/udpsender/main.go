package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const addr = "localhost:42069"

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)

	if err != nil {
		// log.Fatalf("Fatal error when resolving udp address: %s", err)
		fmt.Fprintf(os.Stderr, "Error resolving UDP address: %v\n", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)

	if err != nil {
		// log.Fatalf("Fatal error when establishing udp connection: %s", err)
		fmt.Fprintf(os.Stderr, "Error dialing UDP: %v\n", err)
		os.Exit(1)
	}

	defer conn.Close()

	fmt.Printf("Sending to %s. Type your message and press Enter to send. Press Ctrl+C to exit.\n", addr)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')

		if err != nil {
			// fmt.Printf("error: %s\n", err)
			// continue
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}

		_, err = conn.Write([]byte(line))

		if err != nil {
			// fmt.Printf("error: %s\n", err)
			// continue
			fmt.Fprintf(os.Stderr, "Error sending message: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Message sent: %s", line)
	}
}
