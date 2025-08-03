// Package main demonstrates custom configuration options for the go-json-formatter library.
// This example shows various configuration options and their effects.
//
// To run this example from the project root:
//
//	go run examples/custom_config/main.go
package main

import (
	"fmt"
	"log"

	"github.com/shibukawa/jsonformat"
)

func main() {
	// Example JSON string
	jsonStr := `{"data":[{"type":"user","attributes":{"name":"Alice","age":30}},{"type":"user","attributes":{"name":"Bob","age":25}}],"included":[{"type":"profile","id":"1","attributes":{"bio":"Developer"}}]}`

	fmt.Println("=== Custom Configuration Examples ===")
	fmt.Println("Original JSON:")
	fmt.Println(jsonStr)

	// Example 1: Custom indent size (4 spaces)
	fmt.Println("\n--- Example 1: 4-space indentation ---")
	config1 := jsonformat.NewConfig(
		jsonformat.WithIndentSize(4),
	)
	f1 := jsonformat.NewFormatter(config1)
	formatted1, err := f1.Format(jsonStr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(formatted1)

	// Example 2: Tab indentation
	fmt.Println("\n--- Example 2: Tab indentation ---")
	config2 := jsonformat.NewConfig(
		jsonformat.WithTabs(),
	)
	f2 := jsonformat.NewFormatter(config2)
	formatted2, err := f2.Format(jsonStr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(formatted2)

	// Example 3: Compact at depth 2
	fmt.Println("\n--- Example 3: Compact at depth 2 ---")
	config3 := jsonformat.NewConfig(
		jsonformat.WithCompactDepth(2),
	)
	f3 := jsonformat.NewFormatter(config3)
	formatted3, err := f3.Format(jsonStr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(formatted3)

	// Example 4: No compact formatting
	fmt.Println("\n--- Example 4: No compact formatting ---")
	config4 := jsonformat.NewConfig(
		jsonformat.WithCompactDepth(0),
	)
	f4 := jsonformat.NewFormatter(config4)
	formatted4, err := f4.Format(jsonStr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(formatted4)
}
