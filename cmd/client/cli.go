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
	fmt.Println("Welcome to DoggoDB CLI!")
	fmt.Println("Type 'exit' or 'quit' to leave.")

	for {
		fmt.Print("doggodb> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}

		// Remove the newline character
		input = input[:len(input)-1]

		// Handle exit commands
		if input == "exit" || input == "exit;" || input == "quit" || input == "quit;" {
			fmt.Println("Exiting...")
			break
		}

		// Parse the input
		tokens, err := query.Tokenize(input)
		fmt.Println("DEBUG tokens:", tokens)
		stmt, err := query.Parse(tokens)
		if err != nil {
			fmt.Printf("Invalid command: %v\n", err)
			continue
		}

		// Execute the parsed statement
		switch cmd := stmt.(type) {
		case *query.CreateTableStatement:
			err = client.CreateTable(cmd.Table)
			if err != nil {
				fmt.Printf("Failed to create table: %v\n", err)
			} else {
				fmt.Println("Table '" + cmd.Table + "' created successfully.")
			}

		case *query.InsertStatement:
			if len(cmd.Columns) != len(cmd.Values) {
				fmt.Println("Error: Columns and Values length mismatch")
				return
			}
			// TODO: fix this shit upstream! Is Insert struct values map[string]string or two string slices for columns and values?
			values := make(map[string]string)
			for i, column := range cmd.Columns {
				values[column] = cmd.Values[i]
			}
			err = client.Insert(cmd.Table, values)
			if err != nil {
				fmt.Printf("Failed to insert row: %v\n", err)
			} else {
				fmt.Println("Row inserted successfully.")
			}

		case *query.SelectStatement:
			rows, err := client.Select(cmd.Table, cmd.Columns, cmd.Conditions)
			if err != nil {
				fmt.Printf("Failed to execute SELECT: %v\n", err)
			} else {
				for _, row := range rows {
					fmt.Println(row)
				}
			}

		case *query.UpdateStatement:
			rowsUpdated, err := client.Update(cmd.Table, cmd.Assignments, cmd.Conditions)
			if err != nil {
				fmt.Printf("Failed to execute UPDATE: %v\n", err)
			} else {
				fmt.Printf("Row(s) updated: %d\n", rowsUpdated)
			}

		case *query.DeleteStatement:
			rowsDeleted, err := client.Delete(cmd.Table, cmd.Conditions)
			if err != nil {
				fmt.Printf("Failed to execute DELETE: %v\n", err)
			} else {
				fmt.Printf("Row(s) deleted: %d\n", rowsDeleted)
			}

		default:
			fmt.Println("Unsupported command.")
		}
	}
}
