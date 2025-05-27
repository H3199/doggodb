package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/H3199/doggodb/internal/api"
	"github.com/H3199/doggodb/internal/query"
)

func main() {
	client, err := api.NewClient("localhost:50051")
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer client.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("doggodb> ")
		input, err := reader.ReadString('\n') // Read the whole line
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}

		// Trim whitespace and handle tokenization
		input = input[:len(input)-1] // Remove the trailing newline character
		tokens, err := query.Tokenize(input)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		fmt.Println("DEBUG tokens:", tokens)
	}
}
