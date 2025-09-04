package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

const inputFilePath = "./message.txt"

func main() {
	file, err := os.Open(inputFilePath)

	if err != nil {
		log.Fatalf("could not open %s: %s\n", inputFilePath, err)
	}

	defer file.Close()

	b := make([]byte, 8)

	for {
		_, err := file.Read(b)

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Printf("error: %s\n", err.Error())
			break
		}
		// Could have done that, but useless since the Printf make the casting,
		// still, it looks better
		// str := string(b[:n])
		// fmt.Printf("read: %s\n", str)
		fmt.Printf("read: %s\n", b)
	}
}
