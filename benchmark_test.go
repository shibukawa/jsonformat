package jsonformat

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// BenchmarkFormatterSmallJSON benchmarks formatting of small JSON structures
func BenchmarkFormatterSmallJSON(b *testing.B) {
	input := `{"users":[{"id":1,"name":"Alice"},{"id":2,"name":"Bob"}],"meta":{"count":2}}`
	formatter := NewFormatter(DefaultConfig())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := formatter.Format(input)
		if err != nil {
			b.Fatalf("Formatting failed: %v", err)
		}
	}
}

// BenchmarkFormatterMediumJSON benchmarks formatting of medium-sized JSON structures
func BenchmarkFormatterMediumJSON(b *testing.B) {
	// Create a medium-sized JSON with nested structures
	input := `{
		"users": [
			{"id": 1, "name": "Alice", "profile": {"age": 25, "city": "NYC", "preferences": {"theme": "dark", "notifications": true}}},
			{"id": 2, "name": "Bob", "profile": {"age": 30, "city": "LA", "preferences": {"theme": "light", "notifications": false}}},
			{"id": 3, "name": "Charlie", "profile": {"age": 35, "city": "Chicago", "preferences": {"theme": "auto", "notifications": true}}}
		],
		"meta": {
			"count": 3,
			"page": 1,
			"total_pages": 1,
			"filters": {
				"active": true,
				"roles": ["user", "admin"],
				"created_after": "2023-01-01"
			}
		},
		"config": {
			"api_version": "v1",
			"features": {
				"auth": true,
				"logging": true,
				"metrics": false
			}
		}
	}`
	formatter := NewFormatter(DefaultConfig())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := formatter.Format(input)
		if err != nil {
			b.Fatalf("Formatting failed: %v", err)
		}
	}
}

// BenchmarkFormatterLargeJSON benchmarks formatting of large JSON structures
func BenchmarkFormatterLargeJSON(b *testing.B) {
	// Create a large JSON with many array elements
	var builder strings.Builder
	builder.WriteString(`{"items":[`)
	for i := 0; i < 1000; i++ {
		if i > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(fmt.Sprintf(`{"id":%d,"name":"item%d","data":{"value":%d,"active":true}}`, i, i, i*10))
	}
	builder.WriteString(`],"meta":{"count":1000,"generated":true}}`)

	input := builder.String()
	formatter := NewFormatter(DefaultConfig())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := formatter.Format(input)
		if err != nil {
			b.Fatalf("Formatting failed: %v", err)
		}
	}
}

// BenchmarkFormatterDeeplyNested benchmarks formatting of deeply nested JSON structures
func BenchmarkFormatterDeeplyNested(b *testing.B) {
	// Create a deeply nested structure
	var builder strings.Builder
	depth := 20

	// Build opening braces
	for i := 0; i < depth; i++ {
		builder.WriteString(fmt.Sprintf(`{"level%d":`, i))
	}
	builder.WriteString(`[{"deep":"value","nested":true}]`)

	// Build closing braces
	for i := 0; i < depth; i++ {
		builder.WriteString("}")
	}

	input := builder.String()
	formatter := NewFormatter(DefaultConfig())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := formatter.Format(input)
		if err != nil {
			b.Fatalf("Formatting failed: %v", err)
		}
	}
}

// BenchmarkFormatterArrayHeavy benchmarks formatting of array-heavy JSON structures
func BenchmarkFormatterArrayHeavy(b *testing.B) {
	// Create JSON with many nested arrays
	input := `{
		"matrix": [
			[{"x":1,"y":2},{"x":3,"y":4},{"x":5,"y":6}],
			[{"x":7,"y":8},{"x":9,"y":10},{"x":11,"y":12}],
			[{"x":13,"y":14},{"x":15,"y":16},{"x":17,"y":18}]
		],
		"vectors": [
			[1,2,3,4,5],
			[6,7,8,9,10],
			[11,12,13,14,15]
		],
		"data": [
			{"items":[{"a":1},{"b":2},{"c":3}]},
			{"items":[{"d":4},{"e":5},{"f":6}]},
			{"items":[{"g":7},{"h":8},{"i":9}]}
		]
	}`
	formatter := NewFormatter(DefaultConfig())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := formatter.Format(input)
		if err != nil {
			b.Fatalf("Formatting failed: %v", err)
		}
	}
}

// BenchmarkStandardLibraryComparison benchmarks standard library json.Marshal with indent
func BenchmarkStandardLibrarySmall(b *testing.B) {
	data := map[string]interface{}{
		"users": []map[string]interface{}{
			{"id": 1, "name": "Alice"},
			{"id": 2, "name": "Bob"},
		},
		"meta": map[string]interface{}{
			"count": 2,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			b.Fatalf("Standard library formatting failed: %v", err)
		}
	}
}

// BenchmarkStandardLibraryMedium benchmarks standard library with medium JSON
func BenchmarkStandardLibraryMedium(b *testing.B) {
	data := map[string]interface{}{
		"users": []map[string]interface{}{
			{
				"id":   1,
				"name": "Alice",
				"profile": map[string]interface{}{
					"age":  25,
					"city": "NYC",
					"preferences": map[string]interface{}{
						"theme":         "dark",
						"notifications": true,
					},
				},
			},
			{
				"id":   2,
				"name": "Bob",
				"profile": map[string]interface{}{
					"age":  30,
					"city": "LA",
					"preferences": map[string]interface{}{
						"theme":         "light",
						"notifications": false,
					},
				},
			},
		},
		"meta": map[string]interface{}{
			"count":       2,
			"page":        1,
			"total_pages": 1,
			"filters": map[string]interface{}{
				"active":        true,
				"roles":         []string{"user", "admin"},
				"created_after": "2023-01-01",
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			b.Fatalf("Standard library formatting failed: %v", err)
		}
	}
}

// BenchmarkStandardLibraryLarge benchmarks standard library with large JSON
func BenchmarkStandardLibraryLarge(b *testing.B) {
	items := make([]map[string]interface{}, 1000)
	for i := 0; i < 1000; i++ {
		items[i] = map[string]interface{}{
			"id":   i,
			"name": fmt.Sprintf("item%d", i),
			"data": map[string]interface{}{
				"value":  i * 10,
				"active": true,
			},
		}
	}

	data := map[string]interface{}{
		"items": items,
		"meta": map[string]interface{}{
			"count":     1000,
			"generated": true,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			b.Fatalf("Standard library formatting failed: %v", err)
		}
	}
}

// BenchmarkFormatterWithDifferentConfigs benchmarks different configuration options
func BenchmarkFormatterSpaces2(b *testing.B) {
	input := `{"users":[{"id":1,"name":"Alice"},{"id":2,"name":"Bob"}],"meta":{"count":2}}`
	config := NewConfig(WithIndentSize(2), WithSpaces())
	formatter := NewFormatter(config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := formatter.Format(input)
		if err != nil {
			b.Fatalf("Formatting failed: %v", err)
		}
	}
}

func BenchmarkFormatterSpaces4(b *testing.B) {
	input := `{"users":[{"id":1,"name":"Alice"},{"id":2,"name":"Bob"}],"meta":{"count":2}}`
	config := NewConfig(WithIndentSize(4), WithSpaces())
	formatter := NewFormatter(config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := formatter.Format(input)
		if err != nil {
			b.Fatalf("Formatting failed: %v", err)
		}
	}
}

func BenchmarkFormatterTabs(b *testing.B) {
	input := `{"users":[{"id":1,"name":"Alice"},{"id":2,"name":"Bob"}],"meta":{"count":2}}`
	config := NewConfig(WithTabs())
	formatter := NewFormatter(config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := formatter.Format(input)
		if err != nil {
			b.Fatalf("Formatting failed: %v", err)
		}
	}
}

func BenchmarkFormatterSingleLineDisabled(b *testing.B) {
	input := `{"users":[{"id":1,"name":"Alice"},{"id":2,"name":"Bob"}],"meta":{"count":2}}`
	config := NewConfig(WithCompactDepth(0))
	formatter := NewFormatter(config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := formatter.Format(input)
		if err != nil {
			b.Fatalf("Formatting failed: %v", err)
		}
	}
}

// BenchmarkFormatterMemoryAllocation benchmarks memory allocation patterns
func BenchmarkFormatterMemorySmall(b *testing.B) {
	input := `{"users":[{"id":1,"name":"Alice"}]}`
	formatter := NewFormatter(DefaultConfig())

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := formatter.Format(input)
		if err != nil {
			b.Fatalf("Formatting failed: %v", err)
		}
	}
}

func BenchmarkFormatterMemoryMedium(b *testing.B) {
	input := `{"users":[{"id":1,"name":"Alice","profile":{"age":25,"city":"NYC"}},{"id":2,"name":"Bob","profile":{"age":30,"city":"LA"}}],"meta":{"count":2}}`
	formatter := NewFormatter(DefaultConfig())

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := formatter.Format(input)
		if err != nil {
			b.Fatalf("Formatting failed: %v", err)
		}
	}
}

func BenchmarkFormatterMemoryLarge(b *testing.B) {
	// Create a large JSON structure
	var builder strings.Builder
	builder.WriteString(`{"items":[`)
	for i := 0; i < 100; i++ {
		if i > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(fmt.Sprintf(`{"id":%d,"name":"item%d"}`, i, i))
	}
	builder.WriteString(`]}`)

	input := builder.String()
	formatter := NewFormatter(DefaultConfig())

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := formatter.Format(input)
		if err != nil {
			b.Fatalf("Formatting failed: %v", err)
		}
	}
}

// BenchmarkFormatterStringOperations benchmarks string building operations
func BenchmarkFormatterStringEscaping(b *testing.B) {
	input := `{"message":"He said \"Hello\" and then\\left","path":"C:\\Users\\test","multiline":"Line 1\nLine 2\tTabbed"}`
	formatter := NewFormatter(DefaultConfig())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := formatter.Format(input)
		if err != nil {
			b.Fatalf("Formatting failed: %v", err)
		}
	}
}

func BenchmarkFormatterLargeStrings(b *testing.B) {
	largeString := strings.Repeat("a", 1000) // 1KB string
	input := fmt.Sprintf(`{"large_string":"%s","normal":"value"}`, largeString)
	formatter := NewFormatter(DefaultConfig())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := formatter.Format(input)
		if err != nil {
			b.Fatalf("Formatting failed: %v", err)
		}
	}
}

// BenchmarkFormatterEdgeCases benchmarks edge cases
func BenchmarkFormatterEmptyStructures(b *testing.B) {
	input := `{"empty_object":{},"empty_array":[],"nested":{"empty":{}}}`
	formatter := NewFormatter(DefaultConfig())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := formatter.Format(input)
		if err != nil {
			b.Fatalf("Formatting failed: %v", err)
		}
	}
}

func BenchmarkFormatterNullValues(b *testing.B) {
	input := `{"null_value":null,"array_with_nulls":[null,null,null],"mixed":[1,null,"string",null,true]}`
	formatter := NewFormatter(DefaultConfig())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := formatter.Format(input)
		if err != nil {
			b.Fatalf("Formatting failed: %v", err)
		}
	}
}

func BenchmarkFormatterNumbers(b *testing.B) {
	input := `{"int":42,"float":3.14159,"negative":-123,"zero":0,"scientific":1.23e10,"small":0.000001,"array":[1,2.5,-3,0,1e5,0.1]}`
	formatter := NewFormatter(DefaultConfig())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := formatter.Format(input)
		if err != nil {
			b.Fatalf("Formatting failed: %v", err)
		}
	}
}

// BenchmarkFormatterReusability benchmarks formatter reuse
func BenchmarkFormatterReuse(b *testing.B) {
	inputs := []string{
		`{"users":[{"id":1,"name":"Alice"}]}`,
		`{"products":[{"id":"p1","name":"Laptop"}]}`,
		`{"orders":[{"id":"o1","total":99.99}]}`,
	}
	formatter := NewFormatter(DefaultConfig())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		input := inputs[i%len(inputs)]
		_, err := formatter.Format(input)
		if err != nil {
			b.Fatalf("Formatting failed: %v", err)
		}
	}
}

// BenchmarkFormatterConcurrency benchmarks concurrent usage
func BenchmarkFormatterConcurrent(b *testing.B) {
	input := `{"users":[{"id":1,"name":"Alice"},{"id":2,"name":"Bob"}]}`

	b.RunParallel(func(pb *testing.PB) {
		formatter := NewFormatter(DefaultConfig())
		for pb.Next() {
			_, err := formatter.Format(input)
			if err != nil {
				b.Fatalf("Formatting failed: %v", err)
			}
		}
	})
}
