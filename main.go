package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/SergeyCherepiuk/sha-go/internal/sha"
)

func main() {
	message, err := os.ReadFile("message.txt")
	if err != nil {
		log.Fatal(err)
	}

	message = bytes.TrimSpace(message)
	fmt.Println(string(message))

	hash := sha.Hash(message)
	fmt.Println(string(hash))
}
