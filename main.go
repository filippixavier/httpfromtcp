package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const inputFilePath = "./message.txt"

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
	file, err := os.Open(inputFilePath)

	if err != nil {
		log.Fatalf("could not open %s: %s\n", inputFilePath, err)
	}

	ch := getLinesChannel(file)

	for str := range ch {
		fmt.Printf("read: %s\n", str)
	}
}
