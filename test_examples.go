package jsonformat

import (
	"fmt"
	"log"
)

func testBasicUsage() {
	fmt.Println("=== Testing Basic Usage ===")

	// Example JSON string with nested structures
	jsonStr := `{"users":[{"id":1,"name":"Alice","email":"alice@example.com"},{"id":2,"name":"Bob","email":"bob@example.com"}],"meta":{"count":2,"version":"1.0"}}`

	// Create formatter with default configuration
	config := DefaultConfig()
	f := NewFormatter(config)

	// Format the JSON
	formatted, err := f.Format(jsonStr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Original JSON:")
	fmt.Println(jsonStr)
	fmt.Println("\nFormatted JSON:")
	fmt.Println(formatted)
}

func testCustomConfig() {
	fmt.Println("\n=== Testing Custom Configuration ===")

	jsonStr := `{"data":[{"type":"user","attributes":{"name":"Alice","age":30}}]}`

	// Test 4-space indentation
	config := NewConfig(WithIndentSize(4))
	f := NewFormatter(config)
	formatted, err := f.Format(jsonStr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("4-space indentation:")
	fmt.Println(formatted)

	// Test compact at depth 2
	config2 := NewConfig(WithCompactDepth(2))
	f2 := NewFormatter(config2)
	formatted2, err := f2.Format(jsonStr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nCompact at depth 2:")
	fmt.Println(formatted2)

	// Test no compact formatting
	config3 := NewConfig(WithCompactDepth(0))
	f3 := NewFormatter(config3)
	formatted3, err := f3.Format(jsonStr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nNo compact formatting:")
	fmt.Println(formatted3)
}

func testErrorHandling() {
	fmt.Println("\n=== Testing Error Handling ===")

	config := DefaultConfig()
	f := NewFormatter(config)

	// Test empty input
	_, err := f.Format("")
	if err != nil {
		fmt.Printf("Empty input error: %v\n", err)
	}

	// Test invalid JSON
	_, err = f.Format(`{"name": "Alice", "age": 30`)
	if err != nil {
		fmt.Printf("Invalid JSON error: %v\n", err)
	}

	// Test valid JSON
	formatted, err := f.Format(`{"name": "Alice", "age": 30}`)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Valid JSON formatted successfully:")
		fmt.Println(formatted)
	}
}

func main() {
	testBasicUsage()
	testCustomConfig()
	testErrorHandling()

	fmt.Println("\n=== All examples tested successfully! ===")
}
