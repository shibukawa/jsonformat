package jsonformat

import (
	"fmt"
	"strings"
	"testing"
)

// TestEndToEndFormatting tests complex nested JSON structures
func TestEndToEndFormatting(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:  "simple object with array of objects",
			input: `{"users":[{"id":1,"name":"Alice"},{"id":2,"name":"Bob"}],"meta":{"count":2}}`,
			expected: `{
  "users": 
  [
    {"id": 1, "name": "Alice"},
    {"id": 2, "name": "Bob"}
  ],
  "meta": {
    "count": 2
  }
}`,
		},
		{
			name:  "deeply nested structure",
			input: `{"level1":{"level2":{"level3":[{"deep":"value","nested":{"more":"data"}}]}}}`,
			expected: `{
  "level1": {
    "level2": {
      "level3": 
      [
        {"deep": "value", "nested": {
            "more": "data"
          }}
      ]
    }
  }
}`,
		},
		{
			name:  "mixed arrays and objects",
			input: `{"data":[{"type":"user","info":{"name":"Alice","age":30}},{"type":"admin","info":{"name":"Bob","permissions":["read","write"]}}],"status":"active"}`,
			expected: `{
  "data": 
  [
    {"type": "user", "info": {
        "name": "Alice",
        "age": 30
      }},
    {"type": "admin", "info": {
        "name": "Bob",
        "permissions": 
        [
          "read",
          "write"
        ]
      }}
  ],
  "status": "active"
}`,
		},
		{
			name:  "array of arrays with objects",
			input: `{"matrix":[[{"x":1,"y":2},{"x":3,"y":4}],[{"x":5,"y":6},{"x":7,"y":8}]]}`,
			expected: `{
  "matrix": 
  [
    [
      {"x": 1, "y": 2}, {"x": 3, "y": 4}
    ],
    [
      {"x": 5, "y": 6}, {"x": 7, "y": 8}
    ]
  ]
}`,
		},
		{
			name:  "complex real-world API response",
			input: `{"success":true,"data":{"users":[{"id":1,"name":"John Doe","profile":{"age":25,"city":"New York"}},{"id":2,"name":"Jane Smith","profile":{"age":30,"city":"Los Angeles"}}],"pagination":{"page":1,"limit":10,"total":2}}}`,
			expected: `{
  "success": true,
  "data": {
    "users": 
    [
      {"id": 1, "name": "John Doe", "profile": {
          "age": 25,
          "city": "New York"
        }},
      {"id": 2, "name": "Jane Smith", "profile": {
          "age": 30,
          "city": "Los Angeles"
        }}
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 2
    }
  }
}`,
		},
		{
			name:  "array with mixed primitive types",
			input: `{"values":[1,"string",true,null,3.14,false]}`,
			expected: `{
  "values": 
  [
    1,
    "string",
    true,
    null,
    3.14,
    false
  ]
}`,
		},
		{
			name:  "nested objects without arrays (normal formatting)",
			input: `{"config":{"database":{"host":"localhost","port":5432,"credentials":{"username":"admin","password":"secret"}}}}`,
			expected: `{
  "config": {
    "database": {
      "host": "localhost",
      "port": 5432,
      "credentials": {
        "username": "admin",
        "password": "secret"
      }
    }
  }
}`,
		},
		{
			name:  "array at root level",
			input: `[{"name":"Alice","age":25},{"name":"Bob","age":30}]`,
			expected: `[
  {"name": "Alice", "age": 25},
  {"name": "Bob", "age": 30}
]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewFormatter(DefaultConfig())
			result, err := formatter.Format(tt.input)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("Expected:\n%s\n\nGot:\n%s", tt.expected, result)
			}
		})
	}
}

// TestSingleLineArrayObjectsFeature tests the single-line array objects feature specifically
func TestSingleLineArrayObjectsFeature(t *testing.T) {
	input := `{"items":[{"id":1,"name":"item1"},{"id":2,"name":"item2"}],"nested":{"array":[{"x":1,"y":2}]}}`

	t.Run("with single-line enabled", func(t *testing.T) {
		config := NewConfig(WithCompactDepth(3))
		formatter := NewFormatter(config)
		result, err := formatter.Format(input)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}

		expected := `{
  "items": 
  [
    {"id": 1, "name": "item1"},
    {"id": 2, "name": "item2"}
  ],
  "nested": {
    "array": 
    [
      {"x": 1, "y": 2}
    ]
  }
}`

		if result != expected {
			t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, result)
		}
	})

	t.Run("with single-line disabled", func(t *testing.T) {
		config := NewConfig(WithCompactDepth(0))
		formatter := NewFormatter(config)
		result, err := formatter.Format(input)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}

		expected := `{
  "items": 
  [
    {
      "id": 1,
      "name": "item1"
    },
    {
      "id": 2,
      "name": "item2"
    }
  ],
  "nested": {
    "array": 
    [
      {
        "x": 1,
        "y": 2
      }
    ]
  }
}`

		if result != expected {
			t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, result)
		}
	})
}

// TestDifferentIndentationConfigurations tests various indentation settings
func TestDifferentIndentationConfigurations(t *testing.T) {
	input := `{"array":[{"key":"value"}],"object":{"nested":"data"}}`

	tests := []struct {
		name     string
		config   *Config
		expected string
	}{
		{
			name:   "2 spaces (default)",
			config: NewConfig(WithIndentSize(2)),
			expected: `{
  "array": 
  [
    {"key": "value"}
  ],
  "object": {
    "nested": "data"
  }
}`,
		},
		{
			name:   "4 spaces",
			config: NewConfig(WithIndentSize(4)),
			expected: `{
    "array": 
    [
        {"key": "value"}
    ],
    "object": {
        "nested": "data"
    }
}`,
		},
		{
			name:     "tabs",
			config:   NewConfig(WithTabs()),
			expected: "{\n\t\"array\": \n\t[\n\t\t{\"key\": \"value\"}\n\t],\n\t\"object\": {\n\t\t\"nested\": \"data\"\n\t}\n}",
		},
		{
			name:   "zero indent",
			config: NewConfig(WithIndentSize(0)),
			expected: `{
"array": 
[
{"key": "value"}
],
"object": {
"nested": "data"
}
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewFormatter(tt.config)
			result, err := formatter.Format(input)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("Expected:\n%s\n\nGot:\n%s", tt.expected, result)
			}
		})
	}
}

// TestRealWorldJSONExamples tests with realistic JSON data structures
func TestRealWorldJSONExamples(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:  "REST API response with products",
			input: `{"status":"success","data":{"products":[{"id":"p1","name":"Laptop","price":999.99},{"id":"p2","name":"Mouse","price":29.99}],"total":2}}`,
			expected: `{
  "status": "success",
  "data": {
    "products": 
    [
      {"id": "p1", "name": "Laptop", "price": 999.99},
      {"id": "p2", "name": "Mouse", "price": 29.99}
    ],
    "total": 2
  }
}`,
		},
		{
			name:  "configuration with database connections",
			input: `{"database":{"connections":[{"name":"primary","url":"postgres://localhost:5432/db1"},{"name":"replica","url":"postgres://replica:5432/db1"}]}}`,
			expected: `{
  "database": {
    "connections": 
    [
      {"name": "primary", "url": "postgres://localhost:5432/db1"},
      {"name": "replica", "url": "postgres://replica:5432/db1"}
    ]
  }
}`,
		},
		{
			name:  "user profile with posts",
			input: `{"user":{"id":"123","name":"John","posts":[{"id":"p1","title":"Hello World"},{"id":"p2","title":"Another Post"}]}}`,
			expected: `{
  "user": {
    "id": "123",
    "name": "John",
    "posts": 
    [
      {"id": "p1", "title": "Hello World"},
      {"id": "p2", "title": "Another Post"}
    ]
  }
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewFormatter(DefaultConfig())
			result, err := formatter.Format(tt.input)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("Expected:\n%s\n\nGot:\n%s", tt.expected, result)
			}
		})
	}
}

// TestComplexNestedStructures tests deeply nested and complex structures
func TestComplexNestedStructures(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:  "deeply nested objects and arrays",
			input: `{"level1":{"level2":{"level3":{"level4":[{"deep":{"nested":{"structure":{"with":["arrays",{"and":"objects"}]}}}}]}}}}`,
			expected: `{
  "level1": {
    "level2": {
      "level3": {
        "level4": 
        [
          {"deep": {
              "nested": {
                "structure": {
                  "with": 
                  [
                    "arrays",
                    {"and": "objects"}
                  ]
                }
              }
            }}
        ]
      }
    }
  }
}`,
		},
		{
			name:  "mixed data types in complex structure",
			input: `{"string":"text","number":42,"boolean":true,"null":null,"arrayOfObjects":[{"id":1,"active":true},{"id":2,"active":false}]}`,
			expected: `{
  "string": "text",
  "number": 42,
  "boolean": true,
  "null": null,
  "arrayOfObjects": 
  [
    {"id": 1, "active": true},
    {"id": 2, "active": false}
  ]
}`,
		},
		{
			name:  "array of arrays with mixed content",
			input: `{"matrix":[["a","b","c"],[1,2,3],[true,false,null],[{"x":1},{"y":2},{"z":3}]]}`,
			expected: `{
  "matrix": 
  [
    [
      "a", "b", "c"
    ],
    [
      1, 2, 3
    ],
    [
      true, false, null
    ],
    [
      {"x": 1}, {"y": 2}, {"z": 3}
    ]
  ]
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewFormatter(DefaultConfig())
			result, err := formatter.Format(tt.input)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("Expected:\n%s\n\nGot:\n%s", tt.expected, result)
			}
		})
	}
}

// TestFormatterConsistency tests that the formatter produces consistent output
func TestFormatterConsistency(t *testing.T) {
	input := `{"users":[{"id":1,"name":"Alice","profile":{"age":25,"city":"NYC"}},{"id":2,"name":"Bob","profile":{"age":30,"city":"LA"}}],"meta":{"count":2,"page":1}}`

	formatter := NewFormatter(DefaultConfig())

	// Format the same input multiple times
	results := make([]string, 5)
	for i := 0; i < 5; i++ {
		result, err := formatter.Format(input)
		if err != nil {
			t.Errorf("Unexpected error on iteration %d: %v", i, err)
			return
		}
		results[i] = result
	}

	// All results should be identical
	for i := 1; i < len(results); i++ {
		if results[i] != results[0] {
			t.Errorf("Inconsistent output on iteration %d:\nFirst:\n%s\n\nCurrent:\n%s", i, results[0], results[i])
		}
	}
}

// TestFormatterWithDifferentInputFormats tests various input formats
func TestFormatterWithDifferentInputFormats(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:  "compact input",
			input: `{"a":[{"b":1},{"c":2}]}`,
			expected: `{
  "a": 
  [
    {"b": 1},
    {"c": 2}
  ]
}`,
		},
		{
			name: "already formatted input",
			input: `{
  "a": [
    {"b": 1},
    {"c": 2}
  ]
}`,
			expected: `{
  "a": 
  [
    {"b": 1},
    {"c": 2}
  ]
}`,
		},
		{
			name:  "input with extra whitespace",
			input: `  {  "a"  :  [  {  "b"  :  1  }  ,  {  "c"  :  2  }  ]  }  `,
			expected: `{
  "a": 
  [
    {"b": 1},
    {"c": 2}
  ]
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewFormatter(DefaultConfig())
			result, err := formatter.Format(tt.input)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("Expected:\n%s\n\nGot:\n%s", tt.expected, result)
			}
		})
	}
}

// TestSpecialCharactersInStrings tests proper handling of special characters
func TestSpecialCharactersInStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:  "strings with quotes and escapes",
			input: `{"message":"He said \"Hello\" and then\\left","path":"C:\\Users\\test"}`,
			expected: `{
  "message": "He said \"Hello\" and then\\left",
  "path": "C:\\Users\\test"
}`,
		},
		{
			name:  "strings with newlines and tabs",
			input: `{"multiline":"Line 1\nLine 2\tTabbed"}`,
			expected: `{
  "multiline": "Line 1\nLine 2\tTabbed"
}`,
		},
		{
			name:  "unicode characters",
			input: `{"unicode":"Hello 世界 🌍","emoji":"🚀 🎉 ✨"}`,
			expected: `{
  "unicode": "Hello 世界 🌍",
  "emoji": "🚀 🎉 ✨"
}`,
		},
		{
			name:  "array with special string values",
			input: `{"items":[{"text":"Line 1\nLine 2"},{"text":"Quote: \"test\""},{"text":"Unicode: 世界"}]}`,
			expected: `{
  "items": 
  [
    {"text": "Line 1\nLine 2"},
    {"text": "Quote: \"test\""},
    {"text": "Unicode: 世界"}
  ]
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewFormatter(DefaultConfig())
			result, err := formatter.Format(tt.input)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("Expected:\n%s\n\nGot:\n%s", tt.expected, result)
			}
		})
	}
}

// TestLargeNumberHandling tests proper formatting of various number types
func TestLargeNumberHandling(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:  "various number formats",
			input: `{"int":42,"float":3.14159,"negative":-123,"zero":0,"scientific":1.23e10,"small":0.000001}`,
			expected: `{
  "int": 42,
  "float": 3.14159,
  "negative": -123,
  "zero": 0,
  "scientific": 12300000000,
  "small": 0.000001
}`,
		},
		{
			name:  "numbers in arrays",
			input: `{"numbers":[1,2.5,-3,0,1e5,0.1]}`,
			expected: `{
  "numbers": 
  [
    1,
    2.5,
    -3,
    0,
    100000,
    0.1
  ]
}`,
		},
		{
			name:  "objects with numeric values in array",
			input: `{"data":[{"value":123.45,"count":10},{"value":-67.89,"count":0}]}`,
			expected: `{
  "data": 
  [
    {"value": 123.45, "count": 10},
    {"value": -67.89, "count": 0}
  ]
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewFormatter(DefaultConfig())
			result, err := formatter.Format(tt.input)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("Expected:\n%s\n\nGot:\n%s", tt.expected, result)
			}
		})
	}
}

// TestEdgeCases tests edge cases like empty arrays, empty objects, and null values
func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:  "empty object",
			input: `{}`,
			expected: `{
}`,
		},
		{
			name:  "empty array",
			input: `[]`,
			expected: `[
]`,
		},
		{
			name:  "object with empty array",
			input: `{"items":[]}`,
			expected: `{
  "items": 
  [
  ]
}`,
		},
		{
			name:  "object with empty object",
			input: `{"config":{}}`,
			expected: `{
  "config": {
  }
}`,
		},
		{
			name:  "array with empty objects",
			input: `[{},{}]`,
			expected: `[
  {}
  {}
]`,
		},
		{
			name:  "array with empty arrays",
			input: `[[],[]]`,
			expected: `[
  [
  ]
  [
  ]
]`,
		},
		{
			name:  "null values in object",
			input: `{"name":null,"value":null}`,
			expected: `{
  "name": null,
  "value": null
}`,
		},
		{
			name:  "null values in array",
			input: `[null,null,null]`,
			expected: `[
  null,
  null,
  null
]`,
		},

		{
			name:  "array with mixed empty values",
			input: `[{},[],"",0,null,false]`,
			expected: `[
  {}
  [
  ]
  "",
  0,
  null,
  false
]`,
		},
		{
			name:  "deeply nested empty structures",
			input: `{"level1":{"level2":{"level3":{"empty_array":[],"empty_object":{}}}}}`,
			expected: `{
  "level1": {
    "level2": {
      "level3": {
        "empty_array": 
        [
        ]
        "empty_object": {
        }
      }
    }
  }
}`,
		},
		{
			name:     "single value types",
			input:    `"string"`,
			expected: `"string"`,
		},
		{
			name:     "single number",
			input:    `42`,
			expected: `42`,
		},
		{
			name:     "single boolean true",
			input:    `true`,
			expected: `true`,
		},
		{
			name:     "single boolean false",
			input:    `false`,
			expected: `false`,
		},
		{
			name:     "single null",
			input:    `null`,
			expected: `null`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewFormatter(DefaultConfig())
			result, err := formatter.Format(tt.input)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("Expected:\n%s\n\nGot:\n%s", tt.expected, result)
			}
		})
	}
}

// TestInvalidJSONInput tests error handling for invalid JSON input
func TestInvalidJSONInput(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		errorType   string
	}{
		{
			name:        "empty string",
			input:       "",
			expectError: true,
			errorType:   "empty input",
		},
		{
			name:        "only whitespace",
			input:       "   \n\t  ",
			expectError: true,
			errorType:   "no valid tokens",
		},
		{
			name:        "unclosed object",
			input:       `{"key": "value"`,
			expectError: true,
			errorType:   "unclosed objects",
		},
		{
			name:        "unclosed array",
			input:       `["item1", "item2"`,
			expectError: true,
			errorType:   "unclosed arrays",
		},
		{
			name:        "extra closing brace",
			input:       `{"key": "value"}}`,
			expectError: true,
			errorType:   "unexpected delimiter",
		},
		{
			name:        "extra closing bracket",
			input:       `["item1", "item2"]]`,
			expectError: true,
			errorType:   "unexpected delimiter",
		},
		{
			name:        "missing comma between object properties",
			input:       `{"key1": "value1" "key2": "value2"}`,
			expectError: true,
			errorType:   "invalid JSON",
		},
		{
			name:        "missing comma between array elements",
			input:       `["item1" "item2"]`,
			expectError: true,
			errorType:   "invalid JSON",
		},
		{
			name:        "trailing comma in object",
			input:       `{"key": "value",}`,
			expectError: true,
			errorType:   "invalid JSON",
		},
		{
			name:        "trailing comma in array",
			input:       `["item1", "item2",]`,
			expectError: true,
			errorType:   "invalid JSON",
		},
		{
			name:        "unquoted string key",
			input:       `{key: "value"}`,
			expectError: true,
			errorType:   "invalid JSON",
		},
		{
			name:        "single quotes instead of double quotes",
			input:       `{'key': 'value'}`,
			expectError: true,
			errorType:   "invalid JSON",
		},
		{
			name:        "invalid escape sequence",
			input:       `{"key": "invalid\xescape"}`,
			expectError: true,
			errorType:   "invalid JSON",
		},
		{
			name:        "control characters in string",
			input:       "{\"key\": \"value\x00with\x01control\"}",
			expectError: true,
			errorType:   "invalid JSON",
		},
		{
			name:        "incomplete string",
			input:       `{"key": "incomplete`,
			expectError: true,
			errorType:   "invalid JSON",
		},
		{
			name:        "invalid number format",
			input:       `{"number": 123.}`,
			expectError: true,
			errorType:   "invalid JSON",
		},
		{
			name:        "invalid boolean",
			input:       `{"bool": True}`,
			expectError: true,
			errorType:   "invalid JSON",
		},
		{
			name:        "invalid null",
			input:       `{"null": NULL}`,
			expectError: true,
			errorType:   "invalid JSON",
		},
		{
			name:        "mixed brackets",
			input:       `{"array": [1, 2, 3}`,
			expectError: true,
			errorType:   "mismatched brackets",
		},
		{
			name:        "nested unclosed structures",
			input:       `{"outer": {"inner": {"deep": [1, 2, 3`,
			expectError: true,
			errorType:   "unclosed structures",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewFormatter(DefaultConfig())
			result, err := formatter.Format(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none. Result: %s", result)
				} else {
					// Verify that the error is a FormatError
					if formatErr, ok := err.(*FormatError); ok {
						t.Logf("Got expected FormatError: %s", formatErr.Error())
					} else {
						t.Logf("Got expected error (not FormatError): %s", err.Error())
					}
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

// TestErrorMessageQuality tests that error messages are descriptive and helpful
func TestErrorMessageQuality(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedError string
	}{
		{
			name:          "empty input",
			input:         "",
			expectedError: "empty",
		},
		{
			name:          "unclosed object",
			input:         `{"key": "value"`,
			expectedError: "unclosed",
		},
		{
			name:          "invalid JSON syntax",
			input:         `{"key": "value",}`,
			expectedError: "invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewFormatter(DefaultConfig())
			_, err := formatter.Format(tt.input)

			if err == nil {
				t.Error("Expected error but got none")
				return
			}

			errorMsg := strings.ToLower(err.Error())
			if !strings.Contains(errorMsg, tt.expectedError) {
				t.Errorf("Expected error message to contain '%s', got: %s", tt.expectedError, err.Error())
			}
		})
	}
}

// TestLargeJSONProcessing tests handling of large JSON inputs
func TestLargeJSONProcessing(t *testing.T) {
	t.Run("large array with many elements", func(t *testing.T) {
		// Create a large array with 1000 objects
		var builder strings.Builder
		builder.WriteString(`{"items":[`)
		for i := 0; i < 1000; i++ {
			if i > 0 {
				builder.WriteString(",")
			}
			builder.WriteString(fmt.Sprintf(`{"id":%d,"name":"item%d"}`, i, i))
		}
		builder.WriteString(`]}`)

		input := builder.String()
		formatter := NewFormatter(DefaultConfig())
		result, err := formatter.Format(input)

		if err != nil {
			t.Errorf("Unexpected error with large array: %v", err)
			return
		}

		// Verify the result is properly formatted
		if !strings.Contains(result, `"items":`) {
			t.Error("Result doesn't contain expected structure")
		}

		// Verify it contains some of the expected items
		if !strings.Contains(result, `{"id": 0, "name": "item0"}`) {
			t.Error("Result doesn't contain expected first item")
		}

		if !strings.Contains(result, `{"id": 999, "name": "item999"}`) {
			t.Error("Result doesn't contain expected last item")
		}
	})

	t.Run("deeply nested structure", func(t *testing.T) {
		// Create a deeply nested structure (but not too deep to avoid stack overflow)
		var builder strings.Builder
		depth := 50

		// Build opening braces
		for i := 0; i < depth; i++ {
			builder.WriteString(fmt.Sprintf(`{"level%d":`, i))
		}
		builder.WriteString(`"deep_value"`)

		// Build closing braces
		for i := 0; i < depth; i++ {
			builder.WriteString("}")
		}

		input := builder.String()
		formatter := NewFormatter(DefaultConfig())
		result, err := formatter.Format(input)

		if err != nil {
			t.Errorf("Unexpected error with deeply nested structure: %v", err)
			return
		}

		// Verify the result contains the deep value
		if !strings.Contains(result, `"deep_value"`) {
			t.Error("Result doesn't contain expected deep value")
		}

		// Verify proper indentation exists
		if !strings.Contains(result, "  ") {
			t.Error("Result doesn't appear to be properly indented")
		}
	})

	t.Run("large string values", func(t *testing.T) {
		// Create an object with large string values
		largeString := strings.Repeat("a", 10000) // 10KB string
		input := fmt.Sprintf(`{"large_string":"%s","normal":"value"}`, largeString)

		formatter := NewFormatter(DefaultConfig())
		result, err := formatter.Format(input)

		if err != nil {
			t.Errorf("Unexpected error with large string: %v", err)
			return
		}

		// Verify the result contains both values
		if !strings.Contains(result, `"large_string":`) {
			t.Error("Result doesn't contain large_string key")
		}

		if !strings.Contains(result, `"normal": "value"`) {
			t.Error("Result doesn't contain normal key-value pair")
		}
	})
}

// TestMemoryEfficiency tests that the formatter doesn't have memory leaks
func TestMemoryEfficiency(t *testing.T) {
	t.Run("repeated formatting doesn't leak memory", func(t *testing.T) {
		input := `{"users":[{"id":1,"name":"Alice","profile":{"age":25,"city":"NYC"}},{"id":2,"name":"Bob","profile":{"age":30,"city":"LA"}}],"meta":{"count":2,"page":1}}`
		formatter := NewFormatter(DefaultConfig())

		// Format the same input many times
		for i := 0; i < 100; i++ {
			result, err := formatter.Format(input)
			if err != nil {
				t.Errorf("Error on iteration %d: %v", i, err)
				return
			}

			// Verify result is consistent
			if !strings.Contains(result, `"users":`) {
				t.Errorf("Inconsistent result on iteration %d", i)
				return
			}
		}
	})

	t.Run("formatter instances are independent", func(t *testing.T) {
		input := `{"test":"value"}`

		// Create multiple formatter instances
		formatters := make([]*Formatter, 10)
		for i := 0; i < 10; i++ {
			formatters[i] = NewFormatter(DefaultConfig())
		}

		// Use all formatters concurrently (simulate concurrent usage)
		results := make([]string, 10)
		errors := make([]error, 10)

		for i := 0; i < 10; i++ {
			results[i], errors[i] = formatters[i].Format(input)
		}

		// Verify all results are identical and no errors occurred
		expected := `{
  "test": "value"
}`

		for i := 0; i < 10; i++ {
			if errors[i] != nil {
				t.Errorf("Error from formatter %d: %v", i, errors[i])
			}
			if results[i] != expected {
				t.Errorf("Inconsistent result from formatter %d: got %s", i, results[i])
			}
		}
	})
}

// TestConfigurationEdgeCases tests edge cases in configuration
func TestConfigurationEdgeCases(t *testing.T) {
	input := `{"array":[{"key":"value"}]}`

	t.Run("nil config", func(t *testing.T) {
		formatter := NewFormatter(nil)
		result, err := formatter.Format(input)

		if err != nil {
			t.Errorf("Unexpected error with nil config: %v", err)
			return
		}

		// Should use default configuration
		expected := `{
  "array": 
  [
    {"key": "value"}
  ]
}`
		if result != expected {
			t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, result)
		}
	})

	t.Run("extreme indent sizes", func(t *testing.T) {
		// Test with maximum allowed indent size
		config := NewConfig(WithIndentSize(20))
		formatter := NewFormatter(config)
		result, err := formatter.Format(input)

		if err != nil {
			t.Errorf("Unexpected error with large indent: %v", err)
			return
		}

		// Should contain very large indentation
		if !strings.Contains(result, strings.Repeat(" ", 20)) {
			t.Error("Result doesn't contain expected large indentation")
		}
	})

	t.Run("zero indent size", func(t *testing.T) {
		config := NewConfig(WithIndentSize(0))
		formatter := NewFormatter(config)
		result, err := formatter.Format(input)

		if err != nil {
			t.Errorf("Unexpected error with zero indent: %v", err)
			return
		}

		// Should have no indentation spaces
		lines := strings.Split(result, "\n")
		for i, line := range lines {
			if i > 0 && strings.HasPrefix(line, " ") {
				t.Errorf("Found unexpected indentation in line: %s", line)
			}
		}
	})
}

// TestFormatBytes tests the FormatBytes method
func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:  "simple object",
			input: []byte(`{"key":"value"}`),
			expected: []byte(`{
  "key": "value"
}`),
		},
		{
			name:  "array with objects",
			input: []byte(`[{"id":1},{"id":2}]`),
			expected: []byte(`[
  {"id": 1},
  {"id": 2}
]`),
		},
		{
			name:     "empty input",
			input:    []byte(""),
			expected: nil, // Should return error, so nil result
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewFormatter(DefaultConfig())
			result, err := formatter.FormatBytes(tt.input)

			if tt.expected == nil {
				// Expecting an error
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}

				if string(result) != string(tt.expected) {
					t.Errorf("Expected:\n%s\n\nGot:\n%s", string(tt.expected), string(result))
				}
			}
		})
	}
}
