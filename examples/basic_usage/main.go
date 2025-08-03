// Package main demonstrates basic usage of the go-json-formatter library.
// This example shows how to format JSON with default configuration.
//
// To run this example from the project root:
//
//	go run examples/basic_usage/main.go
package main

import (
	"fmt"
	"log"

	"github.com/shibukawa/jsonformat"
)

func main() {
	// Example JSON string with nested structures
	jsonStr := `{"users":[{"id":1,"name":"Alice","email":"alice@example.com"},{"id":2,"name":"Bob","email":"bob@example.com"}],"meta":{"count":2,"version":"1.0"}}`

	// Create formatter with default configuration (CompactDepth=3)
	config := jsonformat.DefaultConfig()
	f := jsonformat.NewFormatter(config)

	// Format the JSON
	formatted, err := f.Format(jsonStr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== Basic Usage Example ===")
	fmt.Println("Original JSON:")
	fmt.Println(jsonStr)
	fmt.Println("\nFormatted JSON (CompactDepth=3):")
	fmt.Println(formatted)
}
