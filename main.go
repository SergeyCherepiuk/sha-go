package main

import (
	"bytes"
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

	sha.Hash(message)
}
