package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	message, err := os.ReadFile("message.txt")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(message)
}
