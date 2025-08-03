package jsonformat

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestTokenParserProcessToken(t *testing.T) {
	tests := []struct {
		name        string
		token       json.Token
		expectError bool
		errorMsg    string
	}{
		{
			name:        "string token",
			token:       "test",
			expectError: false,
		},
		{
			name:        "number token",
			token:       float64(42),
			expectError: false,
		},
		{
			name:        "boolean true token",
			token:       true,
			expectError: false,
		},
		{
			name:        "boolean false token",
			token:       false,
			expectError: false,
		},
		{
			name:        "null token",
			token:       nil,
			expectError: false,
		},
		{
			name:        "object start delimiter",
			token:       json.Delim('{'),
			expectError: false,
		},
		{
			name:        "object end delimiter",
			token:       json.Delim('}'),
			expectError: true, // Will error because we're not in an object
		},
		{
			name:        "array start delimiter",
			token:       json.Delim('['),
			expectError: false,
		},
		{
			name:        "array end delimiter",
			token:       json.Delim(']'),
			expectError: true, // Will error because we're not in an array
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var builder strings.Builder
			parser := &TokenParser{
				decoder:        nil, // Not needed for this test
				depth:          0,
				inArray:        make([]bool, 0),
				builder:        &builder,
				config:         DefaultConfig(),
				isFirstElement: true,
				expectingKey:   false,
				inputLength:    0,
			}

			err := parser.processToken(tt.token)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestTokenParserHandleDelimiter(t *testing.T) {
	tests := []struct {
		name        string
		delim       json.Delim
		setupParser func(*TokenParser)
		expectError bool
		errorMsg    string
	}{
		{
			name:  "start object",
			delim: '{',
			setupParser: func(p *TokenParser) {
				// Default setup is fine
			},
			expectError: false,
		},
		{
			name:  "end object - valid",
			delim: '}',
			setupParser: func(p *TokenParser) {
				// Set up parser to be in an object
				p.depth = 1
				p.inArray = []bool{false}
				p.expectingKey = true
			},
			expectError: false,
		},
		{
			name:  "end object - invalid (not in object)",
			delim: '}',
			setupParser: func(p *TokenParser) {
				// Default setup - not in object
			},
			expectError: true,
		},
		{
			name:  "start array",
			delim: '[',
			setupParser: func(p *TokenParser) {
				// Default setup is fine
			},
			expectError: false,
		},
		{
			name:  "end array - valid",
			delim: ']',
			setupParser: func(p *TokenParser) {
				// Set up parser to be in an array
				p.depth = 1
				p.inArray = []bool{true}
			},
			expectError: false,
		},
		{
			name:  "end array - invalid (not in array)",
			delim: ']',
			setupParser: func(p *TokenParser) {
				// Default setup - not in array
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var builder strings.Builder
			parser := &TokenParser{
				decoder:        nil,
				depth:          0,
				inArray:        make([]bool, 0),
				builder:        &builder,
				config:         DefaultConfig(),
				isFirstElement: true,
				expectingKey:   false,
				inputLength:    0,
			}

			tt.setupParser(parser)

			err := parser.handleDelimiter(tt.delim)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestTokenParserHandleString(t *testing.T) {
	tests := []struct {
		name           string
		value          string
		setupParser    func(*TokenParser)
		expectError    bool
		expectedOutput string
	}{
		{
			name:  "object key",
			value: "name",
			setupParser: func(p *TokenParser) {
				p.depth = 1
				p.inArray = []bool{false}
				p.expectingKey = true
				p.isFirstElement = true
			},
			expectError:    false,
			expectedOutput: "\n  \"name\": ",
		},
		{
			name:  "string value in object",
			value: "Alice",
			setupParser: func(p *TokenParser) {
				p.depth = 1
				p.inArray = []bool{false}
				p.expectingKey = false
				p.isFirstElement = false
			},
			expectError:    false,
			expectedOutput: `"Alice"`,
		},
		{
			name:  "string value in array",
			value: "test",
			setupParser: func(p *TokenParser) {
				p.depth = 1
				p.inArray = []bool{true}
				p.expectingKey = false
				p.isFirstElement = true
			},
			expectError:    false,
			expectedOutput: "\n  \"test\"",
		},
		{
			name:  "string with escaping",
			value: "test\nwith\ttabs",
			setupParser: func(p *TokenParser) {
				p.depth = 0
				p.expectingKey = false
				p.isFirstElement = true
			},
			expectError:    false,
			expectedOutput: `"test\nwith\ttabs"`,
		},
		{
			name:  "very long string",
			value: strings.Repeat("a", 1000001), // Over 1MB limit
			setupParser: func(p *TokenParser) {
				p.expectingKey = false
				p.isFirstElement = true
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var builder strings.Builder
			parser := &TokenParser{
				decoder:        nil,
				depth:          0,
				inArray:        make([]bool, 0),
				builder:        &builder,
				config:         DefaultConfig(),
				isFirstElement: true,
				expectingKey:   false,
				inputLength:    0,
			}

			tt.setupParser(parser)

			err := parser.handleString(tt.value)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			if !tt.expectError && builder.String() != tt.expectedOutput {
				t.Errorf("Expected output %q, got %q", tt.expectedOutput, builder.String())
			}
		})
	}
}

func TestTokenParserHandleNumber(t *testing.T) {
	tests := []struct {
		name           string
		value          float64
		setupParser    func(*TokenParser)
		expectError    bool
		expectedOutput string
	}{
		{
			name:  "integer number",
			value: 42,
			setupParser: func(p *TokenParser) {
				p.expectingKey = false
				p.isFirstElement = true
			},
			expectError:    false,
			expectedOutput: "42",
		},
		{
			name:  "float number",
			value: 3.14,
			setupParser: func(p *TokenParser) {
				p.expectingKey = false
				p.isFirstElement = true
			},
			expectError:    false,
			expectedOutput: "3.14",
		},
		{
			name:  "zero",
			value: 0,
			setupParser: func(p *TokenParser) {
				p.expectingKey = false
				p.isFirstElement = true
			},
			expectError:    false,
			expectedOutput: "0",
		},
		{
			name:  "negative number",
			value: -123.45,
			setupParser: func(p *TokenParser) {
				p.expectingKey = false
				p.isFirstElement = true
			},
			expectError:    false,
			expectedOutput: "-123.45",
		},
		{
			name:  "number as object key - should error",
			value: 42,
			setupParser: func(p *TokenParser) {
				p.expectingKey = true
				p.isFirstElement = true
			},
			expectError: true,
		},
		{
			name:  "number in array",
			value: 123,
			setupParser: func(p *TokenParser) {
				p.depth = 1
				p.inArray = []bool{true}
				p.expectingKey = false
				p.isFirstElement = true
			},
			expectError:    false,
			expectedOutput: "\n  123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var builder strings.Builder
			parser := &TokenParser{
				decoder:        nil,
				depth:          0,
				inArray:        make([]bool, 0),
				builder:        &builder,
				config:         DefaultConfig(),
				isFirstElement: true,
				expectingKey:   false,
				inputLength:    0,
			}

			tt.setupParser(parser)

			err := parser.handleNumber(tt.value)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			if !tt.expectError && builder.String() != tt.expectedOutput {
				t.Errorf("Expected output %q, got %q", tt.expectedOutput, builder.String())
			}
		})
	}
}

func TestTokenParserHandleBoolean(t *testing.T) {
	tests := []struct {
		name           string
		value          bool
		setupParser    func(*TokenParser)
		expectError    bool
		expectedOutput string
	}{
		{
			name:  "true value",
			value: true,
			setupParser: func(p *TokenParser) {
				p.expectingKey = false
				p.isFirstElement = true
			},
			expectError:    false,
			expectedOutput: "true",
		},
		{
			name:  "false value",
			value: false,
			setupParser: func(p *TokenParser) {
				p.expectingKey = false
				p.isFirstElement = true
			},
			expectError:    false,
			expectedOutput: "false",
		},
		{
			name:  "boolean as object key - should error",
			value: true,
			setupParser: func(p *TokenParser) {
				p.expectingKey = true
				p.isFirstElement = true
			},
			expectError: true,
		},
		{
			name:  "boolean in array",
			value: true,
			setupParser: func(p *TokenParser) {
				p.depth = 1
				p.inArray = []bool{true}
				p.expectingKey = false
				p.isFirstElement = true
			},
			expectError:    false,
			expectedOutput: "\n  true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var builder strings.Builder
			parser := &TokenParser{
				decoder:        nil,
				depth:          0,
				inArray:        make([]bool, 0),
				builder:        &builder,
				config:         DefaultConfig(),
				isFirstElement: true,
				expectingKey:   false,
				inputLength:    0,
			}

			tt.setupParser(parser)

			err := parser.handleBoolean(tt.value)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			if !tt.expectError && builder.String() != tt.expectedOutput {
				t.Errorf("Expected output %q, got %q", tt.expectedOutput, builder.String())
			}
		})
	}
}

func TestTokenParserHandleNull(t *testing.T) {
	tests := []struct {
		name           string
		setupParser    func(*TokenParser)
		expectError    bool
		expectedOutput string
	}{
		{
			name: "null value",
			setupParser: func(p *TokenParser) {
				p.expectingKey = false
				p.isFirstElement = true
			},
			expectError:    false,
			expectedOutput: "null",
		},
		{
			name: "null as object key - should error",
			setupParser: func(p *TokenParser) {
				p.expectingKey = true
				p.isFirstElement = true
			},
			expectError: true,
		},
		{
			name: "null in array",
			setupParser: func(p *TokenParser) {
				p.depth = 1
				p.inArray = []bool{true}
				p.expectingKey = false
				p.isFirstElement = true
			},
			expectError:    false,
			expectedOutput: "\n  null",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var builder strings.Builder
			parser := &TokenParser{
				decoder:        nil,
				depth:          0,
				inArray:        make([]bool, 0),
				builder:        &builder,
				config:         DefaultConfig(),
				isFirstElement: true,
				expectingKey:   false,
				inputLength:    0,
			}

			tt.setupParser(parser)

			err := parser.handleNull()

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			if !tt.expectError && builder.String() != tt.expectedOutput {
				t.Errorf("Expected output %q, got %q", tt.expectedOutput, builder.String())
			}
		})
	}
}

func TestTokenParserStateTransitions(t *testing.T) {
	t.Run("enterArray", func(t *testing.T) {
		var builder strings.Builder
		parser := &TokenParser{
			decoder:        nil,
			depth:          0,
			inArray:        make([]bool, 0),
			builder:        &builder,
			config:         DefaultConfig(),
			isFirstElement: true,
			expectingKey:   false,
			inputLength:    0,
		}

		err := parser.enterArray()
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if parser.depth != 1 {
			t.Errorf("Expected depth 1, got %d", parser.depth)
		}

		if len(parser.inArray) != 1 || !parser.inArray[0] {
			t.Error("Expected inArray to contain [true]")
		}
	})

	t.Run("exitArray", func(t *testing.T) {
		var builder strings.Builder
		parser := &TokenParser{
			decoder:        nil,
			depth:          1,
			inArray:        []bool{true},
			builder:        &builder,
			config:         DefaultConfig(),
			isFirstElement: true,
			expectingKey:   false,
			inputLength:    0,
		}

		err := parser.exitArray()
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if parser.depth != 0 {
			t.Errorf("Expected depth 0, got %d", parser.depth)
		}

		if len(parser.inArray) != 0 {
			t.Error("Expected inArray to be empty")
		}
	})

	t.Run("enterObject", func(t *testing.T) {
		var builder strings.Builder
		parser := &TokenParser{
			decoder:        nil,
			depth:          0,
			inArray:        make([]bool, 0),
			builder:        &builder,
			config:         DefaultConfig(),
			isFirstElement: true,
			expectingKey:   false,
			inputLength:    0,
		}

		err := parser.enterObject()
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if parser.depth != 1 {
			t.Errorf("Expected depth 1, got %d", parser.depth)
		}

		if len(parser.inArray) != 1 || parser.inArray[0] {
			t.Error("Expected inArray to contain [false]")
		}
	})

	t.Run("exitObject", func(t *testing.T) {
		var builder strings.Builder
		parser := &TokenParser{
			decoder:        nil,
			depth:          1,
			inArray:        []bool{false},
			builder:        &builder,
			config:         DefaultConfig(),
			isFirstElement: true,
			expectingKey:   false,
			inputLength:    0,
		}

		err := parser.exitObject()
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if parser.depth != 0 {
			t.Errorf("Expected depth 0, got %d", parser.depth)
		}

		if len(parser.inArray) != 0 {
			t.Error("Expected inArray to be empty")
		}
	})
}

func TestTokenParserDepthTracking(t *testing.T) {
	var builder strings.Builder
	parser := &TokenParser{
		decoder:        nil,
		depth:          0,
		inArray:        make([]bool, 0),
		builder:        &builder,
		config:         DefaultConfig(),
		isFirstElement: true,
		expectingKey:   false,
		inputLength:    0,
	}

	// Test nested structures
	// Start with object
	err := parser.enterObject()
	if err != nil {
		t.Fatalf("Failed to enter object: %v", err)
	}
	if parser.depth != 1 {
		t.Errorf("Expected depth 1, got %d", parser.depth)
	}

	// Enter array within object
	err = parser.enterArray()
	if err != nil {
		t.Fatalf("Failed to enter array: %v", err)
	}
	if parser.depth != 2 {
		t.Errorf("Expected depth 2, got %d", parser.depth)
	}

	// Enter object within array
	err = parser.enterObject()
	if err != nil {
		t.Fatalf("Failed to enter nested object: %v", err)
	}
	if parser.depth != 3 {
		t.Errorf("Expected depth 3, got %d", parser.depth)
	}

	// Exit nested object
	err = parser.exitObject()
	if err != nil {
		t.Fatalf("Failed to exit nested object: %v", err)
	}
	if parser.depth != 2 {
		t.Errorf("Expected depth 2, got %d", parser.depth)
	}

	// Exit array
	err = parser.exitArray()
	if err != nil {
		t.Fatalf("Failed to exit array: %v", err)
	}
	if parser.depth != 1 {
		t.Errorf("Expected depth 1, got %d", parser.depth)
	}

	// Exit object
	err = parser.exitObject()
	if err != nil {
		t.Fatalf("Failed to exit object: %v", err)
	}
	if parser.depth != 0 {
		t.Errorf("Expected depth 0, got %d", parser.depth)
	}
}

func TestTokenParserArrayContextTracking(t *testing.T) {
	var builder strings.Builder
	parser := &TokenParser{
		decoder:        nil,
		depth:          0,
		inArray:        make([]bool, 0),
		builder:        &builder,
		config:         DefaultConfig(),
		isFirstElement: true,
		expectingKey:   false,
		inputLength:    0,
	}

	// Initially not in array
	if parser.isInArray() {
		t.Error("Expected not to be in array initially")
	}

	// Enter object - still not in array
	err := parser.enterObject()
	if err != nil {
		t.Fatalf("Failed to enter object: %v", err)
	}
	if parser.isInArray() {
		t.Error("Expected not to be in array after entering object")
	}

	// Enter array - now in array
	err = parser.enterArray()
	if err != nil {
		t.Fatalf("Failed to enter array: %v", err)
	}
	if !parser.isInArray() {
		t.Error("Expected to be in array after entering array")
	}

	// Enter object within array - still in array context (parent)
	err = parser.enterObject()
	if err != nil {
		t.Fatalf("Failed to enter object within array: %v", err)
	}
	if parser.isInArray() {
		t.Error("Expected not to be in array when in object (even if parent is array)")
	}

	// Exit object - back to array context
	err = parser.exitObject()
	if err != nil {
		t.Fatalf("Failed to exit object: %v", err)
	}
	if !parser.isInArray() {
		t.Error("Expected to be back in array context")
	}

	// Exit array - back to object context
	err = parser.exitArray()
	if err != nil {
		t.Fatalf("Failed to exit array: %v", err)
	}
	if parser.isInArray() {
		t.Error("Expected not to be in array after exiting array")
	}
}

func TestTokenParserSingleLineObjectDetection(t *testing.T) {
	tests := []struct {
		name           string
		setupParser    func(*TokenParser)
		expectedInline bool
	}{
		{
			name: "root level object - not inline",
			setupParser: func(p *TokenParser) {
				p.depth = 1
				p.inArray = []bool{false}
			},
			expectedInline: false,
		},
		{
			name: "object in array - should be inline",
			setupParser: func(p *TokenParser) {
				p.depth = 2
				p.inArray = []bool{true, false} // Parent is array, current is object
			},
			expectedInline: true,
		},
		{
			name: "nested object not in array - not inline",
			setupParser: func(p *TokenParser) {
				p.depth = 2
				p.inArray = []bool{false, false} // Both parent and current are objects
			},
			expectedInline: false,
		},
		{
			name: "deeply nested object in array - should be inline",
			setupParser: func(p *TokenParser) {
				p.depth = 3
				p.inArray = []bool{false, true, false} // Grandparent object, parent array, current object
			},
			expectedInline: true,
		},
		{
			name: "single line disabled - not inline",
			setupParser: func(p *TokenParser) {
				p.depth = 2
				p.inArray = []bool{true, false}
				p.config.CompactDepth = 0
			},
			expectedInline: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var builder strings.Builder
			parser := &TokenParser{
				decoder:        nil,
				depth:          0,
				inArray:        make([]bool, 0),
				builder:        &builder,
				config:         DefaultConfig(),
				isFirstElement: true,
				expectingKey:   false,
				inputLength:    0,
			}

			tt.setupParser(parser)

			result := parser.shouldFormatCompact()
			if result != tt.expectedInline {
				t.Errorf("Expected shouldFormatCompact() to return %t, got %t", tt.expectedInline, result)
			}
		})
	}
}

func TestTokenParserStateValidation(t *testing.T) {
	t.Run("exitArray with invalid state", func(t *testing.T) {
		tests := []struct {
			name        string
			setupParser func(*TokenParser)
			expectError bool
		}{
			{
				name: "exit array at depth 0",
				setupParser: func(p *TokenParser) {
					p.depth = 0
					p.inArray = []bool{}
				},
				expectError: true,
			},
			{
				name: "exit array with empty inArray stack",
				setupParser: func(p *TokenParser) {
					p.depth = 1
					p.inArray = []bool{}
				},
				expectError: true,
			},
			{
				name: "exit array when in object context",
				setupParser: func(p *TokenParser) {
					p.depth = 1
					p.inArray = []bool{false} // In object, not array
				},
				expectError: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var builder strings.Builder
				parser := &TokenParser{
					decoder:        nil,
					depth:          0,
					inArray:        make([]bool, 0),
					builder:        &builder,
					config:         DefaultConfig(),
					isFirstElement: true,
					expectingKey:   false,
					inputLength:    0,
				}

				tt.setupParser(parser)

				err := parser.exitArray()
				if tt.expectError && err == nil {
					t.Error("Expected error but got none")
				}
				if !tt.expectError && err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			})
		}
	})

	t.Run("exitObject with invalid state", func(t *testing.T) {
		tests := []struct {
			name        string
			setupParser func(*TokenParser)
			expectError bool
		}{
			{
				name: "exit object at depth 0",
				setupParser: func(p *TokenParser) {
					p.depth = 0
					p.inArray = []bool{}
				},
				expectError: true,
			},
			{
				name: "exit object with empty inArray stack",
				setupParser: func(p *TokenParser) {
					p.depth = 1
					p.inArray = []bool{}
				},
				expectError: true,
			},
			{
				name: "exit object when in array context",
				setupParser: func(p *TokenParser) {
					p.depth = 1
					p.inArray = []bool{true} // In array, not object
				},
				expectError: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var builder strings.Builder
				parser := &TokenParser{
					decoder:        nil,
					depth:          0,
					inArray:        make([]bool, 0),
					builder:        &builder,
					config:         DefaultConfig(),
					isFirstElement: true,
					expectingKey:   false,
					inputLength:    0,
				}

				tt.setupParser(parser)

				err := parser.exitObject()
				if tt.expectError && err == nil {
					t.Error("Expected error but got none")
				}
				if !tt.expectError && err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			})
		}
	})
}

func TestTokenParserIndentationHandling(t *testing.T) {
	t.Run("writeIndent with spaces", func(t *testing.T) {
		var builder strings.Builder
		parser := &TokenParser{
			decoder:        nil,
			depth:          2,
			inArray:        make([]bool, 0),
			builder:        &builder,
			config:         &Config{IndentSize: 4, UseTab: false, CompactDepth: 3},
			isFirstElement: true,
			expectingKey:   false,
			inputLength:    0,
		}

		err := parser.writeIndent()
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		expected := "        " // 8 spaces (2 * 4)
		if builder.String() != expected {
			t.Errorf("Expected %q, got %q", expected, builder.String())
		}
	})

	t.Run("writeIndent with tabs", func(t *testing.T) {
		var builder strings.Builder
		parser := &TokenParser{
			decoder:        nil,
			depth:          3,
			inArray:        make([]bool, 0),
			builder:        &builder,
			config:         &Config{IndentSize: 2, UseTab: true, CompactDepth: 3},
			isFirstElement: true,
			expectingKey:   false,
			inputLength:    0,
		}

		err := parser.writeIndent()
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		expected := "\t\t\t" // 3 tabs
		if builder.String() != expected {
			t.Errorf("Expected %q, got %q", expected, builder.String())
		}
	})

	t.Run("writeNewlineAndIndent", func(t *testing.T) {
		var builder strings.Builder
		parser := &TokenParser{
			decoder:        nil,
			depth:          1,
			inArray:        make([]bool, 0),
			builder:        &builder,
			config:         &Config{IndentSize: 2, UseTab: false, CompactDepth: 3},
			isFirstElement: true,
			expectingKey:   false,
			inputLength:    0,
		}

		err := parser.writeNewlineAndIndent()
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		expected := "\n  " // newline + 2 spaces
		if builder.String() != expected {
			t.Errorf("Expected %q, got %q", expected, builder.String())
		}
	})
}

func TestTokenParserUtilityFunctions(t *testing.T) {
	t.Run("escapeString", func(t *testing.T) {
		var builder strings.Builder
		parser := &TokenParser{
			decoder:        nil,
			depth:          0,
			inArray:        make([]bool, 0),
			builder:        &builder,
			config:         DefaultConfig(),
			isFirstElement: true,
			expectingKey:   false,
			inputLength:    0,
		}

		tests := []struct {
			input    string
			expected string
		}{
			{"simple", "simple"},
			{"with\nnewline", "with\\nnewline"},
			{"with\ttab", "with\\ttab"},
			{"with\"quote", "with\\\"quote"},
			{"with\\backslash", "with\\\\backslash"},
		}

		for _, tt := range tests {
			result, err := parser.escapeString(tt.input)
			if err != nil {
				t.Errorf("Expected no error for input %q, got: %v", tt.input, err)
			}
			if result != tt.expected {
				t.Errorf("For input %q, expected %q, got %q", tt.input, tt.expected, result)
			}
		}
	})

	t.Run("formatNumber", func(t *testing.T) {
		var builder strings.Builder
		parser := &TokenParser{
			decoder:        nil,
			depth:          0,
			inArray:        make([]bool, 0),
			builder:        &builder,
			config:         DefaultConfig(),
			isFirstElement: true,
			expectingKey:   false,
			inputLength:    0,
		}

		tests := []struct {
			input    float64
			expected string
		}{
			{42, "42"},
			{3.14, "3.14"},
			{0, "0"},
			{-123.45, "-123.45"},
		}

		for _, tt := range tests {
			result, err := parser.formatNumber(tt.input)
			if err != nil {
				t.Errorf("Expected no error for input %f, got: %v", tt.input, err)
			}
			if result != tt.expected {
				t.Errorf("For input %f, expected %q, got %q", tt.input, tt.expected, result)
			}
		}
	})
}
