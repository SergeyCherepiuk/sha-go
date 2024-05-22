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

	hash := sha.Sum(message)
	fmt.Println(hash.String())
	fmt.Println(hash.Bits())
}
