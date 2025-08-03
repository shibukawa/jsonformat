// Package main demonstrates error handling patterns for the go-json-formatter library.
// This example shows how to handle various error conditions gracefully.
//
// To run this example from the project root:
//
//	go run examples/error_handling/main.go
package main

import (
	"errors"
	"fmt"

	"github.com/shibukawa/jsonformat"
)

func main() {
	fmt.Println("=== Error Handling Examples ===")

	config := jsonformat.DefaultConfig()
	f := jsonformat.NewFormatter(config)

	// Example 1: Empty input
	fmt.Println("--- Example 1: Empty input ---")
	_, err := f.Format("")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Example 2: Invalid JSON - missing closing brace
	fmt.Println("\n--- Example 2: Invalid JSON (missing closing brace) ---")
	invalidJSON1 := `{"name": "Alice", "age": 30`
	_, err = f.Format(invalidJSON1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		// Check if it's a FormatError and unwrap if needed
		var formatErr *jsonformat.FormatError
		if errors.As(err, &formatErr) {
			fmt.Printf("Format error details: %s\n", formatErr.Error())
			if formatErr.Unwrap() != nil {
				fmt.Printf("Underlying error: %v\n", formatErr.Unwrap())
			}
		}
	}

	// Example 3: Invalid JSON - malformed structure
	fmt.Println("\n--- Example 3: Invalid JSON (malformed structure) ---")
	invalidJSON2 := `{"name": "Alice",, "age": 30}`
	_, err = f.Format(invalidJSON2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Example 4: Valid JSON that formats successfully
	fmt.Println("\n--- Example 4: Valid JSON (successful formatting) ---")
	validJSON := `{"name": "Alice", "age": 30}`
	formatted, err := f.Format(validJSON)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Successfully formatted:")
		fmt.Println(formatted)
	}

	// Example 5: Using FormatBytes method
	fmt.Println("\n--- Example 5: Using FormatBytes method ---")
	jsonBytes := []byte(`{"items": [{"id": 1}, {"id": 2}]}`)
	formattedBytes, err := f.FormatBytes(jsonBytes)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Successfully formatted bytes:")
		fmt.Println(string(formattedBytes))
	}
}
