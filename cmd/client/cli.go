package main

import (
	"fmt"
	"log"

	"github.com/H3199/doggodb/internal/api"
)

func main() {
	// TODO: No hardcoded values
	client, err := api.NewClient("localhost:50051")
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer client.Close()

	for {
		var input string
		fmt.Print("doggodb> ")
		fmt.Scanln(&input)

		switch input {
		case "exit", "quit":
			fmt.Println("Exiting CLI.")
			return
		case "create":
			err := client.CreateTable("users")
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("Table created successfully.")
			}
		case "insert":
			err := client.Insert("users", map[string]string{
				"id":   "1",
				"name": "John Doe",
				"age":  "30",
			})
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("Data inserted successfully.")
			}
		case "select":
			rows, err := client.Select("users", []string{"id", "name", "age"}, "age > 20")
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				for _, row := range rows {
					fmt.Printf("%s %s %s\n", row.Values["id"], row.Values["name"], row.Values["age"])
				}
			}
		default:
			fmt.Println("Unknown command.")
		}
	}
}
