// Copyright 2024 Yoshiki Shibukawa
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package jsonformat provides a JSON formatter with specialized formatting rules.
//
// The go-json-formatter library formats JSON with custom indentation and
// the unique feature of displaying objects within arrays on a single line
// for improved readability of array-heavy JSON data.
//
// Key Features:
//   - Custom indentation (spaces or tabs with configurable size)
//   - Single-line formatting for objects within arrays
//   - Flexible configuration using functional options
//   - Comprehensive error handling with position information
//   - Memory-efficient streaming token-based parsing
//
// Basic Usage:
//
//	config := DefaultConfig()
//	formatter := NewFormatter(config)
//	formatted, err := formatter.Format(jsonString)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(formatted)
//
// Custom Configuration:
//
//	config := NewConfig(
//	    WithIndentSize(4),
//	    WithTabs(),
//	    WithCompactDepth(2),
//	)
//	formatter := NewFormatter(config)
//
// The formatter handles various JSON structures and provides detailed
// error information for invalid input. It supports both string and
// byte slice input through Format() and FormatBytes() methods.
package jsonformat

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Config holds configuration options for JSON formatting.
// It allows customization of indentation style and formatting behavior.
type Config struct {
	// IndentSize specifies the number of spaces to use for each level of indentation.
	// Must be between 0 and 20. Default is 2.
	IndentSize int

	// UseTab determines whether to use tabs instead of spaces for indentation.
	// When true, IndentSize is ignored. Default is false.
	UseTab bool

	// CompactDepth specifies the depth at which elements should be formatted on a single line.
	// Elements at this depth or deeper will be formatted compactly without line breaks.
	// A value of 0 disables compact formatting. Default is 3.
	CompactDepth int
}

// ConfigOption is a functional option for configuring the formatter.
// It allows for flexible configuration using the functional options pattern.
type ConfigOption func(*Config)

// DefaultConfig returns a Config with sensible default values.
// Default configuration uses 2-space indentation and compact formatting at depth 3.
//
// Returns:
//   - IndentSize: 2
//   - UseTab: false
//   - CompactDepth: 3
func DefaultConfig() *Config {
	return &Config{
		IndentSize:   2,
		UseTab:       false,
		CompactDepth: 3,
	}
}

// NewConfig creates a new Config with the provided options.
// It starts with default values and applies the given options in order.
// If any configuration validation fails, it returns the default configuration.
//
// Example:
//
//	config := NewConfig(
//	    WithIndentSize(4),
//	    WithTabs(),
//	    WithCompactDepth(2),
//	)
func NewConfig(options ...ConfigOption) *Config {
	config := DefaultConfig()
	for _, option := range options {
		option(config)
	}

	// Validate configuration parameters
	if err := validateConfig(config); err != nil {
		// Return default config if validation fails
		return DefaultConfig()
	}

	return config
}

// validateConfig validates configuration parameters
func validateConfig(config *Config) error {
	if config == nil {
		return NewFormatError("config cannot be nil")
	}

	if config.IndentSize < 0 {
		return NewFormatError("IndentSize must be non-negative")
	}

	if config.IndentSize > 20 {
		return NewFormatError("IndentSize must not exceed 20 spaces")
	}

	if config.CompactDepth < 0 {
		return NewFormatError("CompactDepth must be non-negative")
	}

	return nil
}

// WithIndentSize sets the number of spaces to use for indentation.
// The size must be between 0 and 20. Invalid values are ignored.
//
// Example:
//
//	config := NewConfig(WithIndentSize(4)) // Use 4 spaces per indent level
func WithIndentSize(size int) ConfigOption {
	return func(c *Config) {
		if size >= 0 && size <= 20 {
			c.IndentSize = size
		}
	}
}

// WithTabs enables tab indentation instead of spaces.
// When enabled, the IndentSize setting is ignored.
//
// Example:
//
//	config := NewConfig(WithTabs()) // Use tabs for indentation
func WithTabs() ConfigOption {
	return func(c *Config) {
		c.UseTab = true
	}
}

// WithSpaces disables tab indentation and uses spaces instead.
// This is the default behavior, but can be used to explicitly override
// tab indentation in configuration chains.
//
// Example:
//
//	config := NewConfig(WithSpaces()) // Explicitly use spaces
func WithSpaces() ConfigOption {
	return func(c *Config) {
		c.UseTab = false
	}
}

// WithCompactDepth sets the depth at which elements should be formatted compactly.
// Elements at this depth or deeper will be formatted on a single line without line breaks.
// A value of 0 disables compact formatting entirely.
//
// Example with CompactDepth=3:
//
//	{"users": [{"name": "Alice", "age": 30}]}
//	The object {"name": "Alice", "age": 30} is at depth 3 and will be compact.
//
// Example:
//
//	config := NewConfig(WithCompactDepth(2)) // Compact at depth 2 and deeper
func WithCompactDepth(depth int) ConfigOption {
	return func(c *Config) {
		if depth >= 0 {
			c.CompactDepth = depth
		}
	}
}

// Formatter handles JSON formatting with custom rules.
// It provides methods to format JSON strings and byte slices according
// to the configured formatting options.
type Formatter struct {
	config *Config
}

// NewFormatter creates a new Formatter with the given configuration.
// If config is nil, it uses the default configuration.
//
// Example:
//
//	config := DefaultConfig()
//	formatter := NewFormatter(config)
func NewFormatter(config *Config) *Formatter {
	if config == nil {
		config = DefaultConfig()
	}
	return &Formatter{
		config: config,
	}
}

// Format formats a JSON string according to the configured rules.
// It parses the input JSON and applies custom formatting, including
// single-line objects within arrays if enabled.
//
// The method handles various error conditions gracefully:
//   - Empty input strings
//   - Invalid JSON syntax
//   - Malformed JSON structures
//   - Deeply nested structures (max depth: 100)
//
// Returns the formatted JSON string and any error encountered.
//
// Example:
//
//	formatted, err := formatter.Format(`{"users":[{"id":1,"name":"Alice"}]}`)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(formatted)
func (f *Formatter) Format(jsonStr string) (result string, err error) {
	// Implement panic recovery to handle unexpected errors gracefully
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				err = WrapFormatError("panic during formatting", v)
			case string:
				err = NewFormatError("panic during formatting: " + v)
			default:
				err = NewFormatError("unexpected panic during formatting")
			}
			result = ""
		}
	}()

	// Validate input
	if jsonStr == "" {
		return "", NewFormatError("input JSON string is empty")
	}

	// Create a decoder from the input string
	reader := strings.NewReader(jsonStr)
	decoder := json.NewDecoder(reader)

	// Create a string builder for output
	var builder strings.Builder

	// Create token parser with decoder and configuration
	parser := &TokenParser{
		decoder:        decoder,
		depth:          0,
		inArray:        make([]bool, 0),
		builder:        &builder,
		config:         f.config,
		isFirstElement: true,
		expectingKey:   false,
		inputLength:    len(jsonStr),
	}

	// Process all tokens sequentially
	tokenCount := 0
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				// EOF indicates we've processed all tokens successfully
				break
			}
			// Calculate approximate position in input
			position := parser.calculatePosition(reader)
			return "", WrapFormatErrorWithPosition("invalid JSON input", position, err)
		}

		tokenCount++
		if tokenCount > 10000 { // Prevent infinite loops with malformed JSON
			return "", NewFormatError("JSON structure too complex or malformed (too many tokens)")
		}

		err = parser.processToken(token)
		if err != nil {
			return "", err
		}
	}

	// Validate that we ended in a valid state
	if parser.depth != 0 {
		return "", NewFormatError("malformed JSON: unclosed objects or arrays")
	}

	// Validate that we have at least one token (not just whitespace)
	if tokenCount == 0 {
		return "", NewFormatError("input contains no valid JSON tokens")
	}

	return builder.String(), nil
}

// FormatBytes formats JSON bytes according to the configured rules.
// It converts the byte slice to a string, formats it using Format(),
// and returns the result as bytes.
//
// This method is convenient when working with JSON data as byte slices,
// such as when reading from files or network responses.
//
// Returns the formatted JSON as bytes and any error encountered.
//
// Example:
//
//	jsonBytes := []byte(`{"users":[{"id":1,"name":"Alice"}]}`)
//	formatted, err := formatter.FormatBytes(jsonBytes)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(string(formatted))
func (f *Formatter) FormatBytes(jsonBytes []byte) ([]byte, error) {
	// Convert bytes to string, format, then convert back to bytes
	formatted, err := f.Format(string(jsonBytes))
	if err != nil {
		return nil, err
	}
	return []byte(formatted), nil
}

// TokenParser handles token-based JSON parsing and formatting
type TokenParser struct {
	decoder        *json.Decoder
	depth          int
	inArray        []bool // Stack to track array context at each depth
	builder        *strings.Builder
	config         *Config
	isFirstElement bool // Track if this is the first element in current context
	expectingKey   bool // Track if we're expecting an object key next
	inputLength    int  // Length of original input for position calculation
}

// processToken processes a single JSON token with type switching
func (p *TokenParser) processToken(token json.Token) error {
	// Validate parser state before processing
	if p.depth < 0 {
		return NewFormatError("invalid parser state: negative depth")
	}
	if p.depth > 100 { // Prevent stack overflow with deeply nested structures
		return NewFormatError("JSON structure too deeply nested (max depth: 100)")
	}

	switch v := token.(type) {
	case json.Delim:
		return p.handleDelimiter(v)
	case string:
		return p.handleString(v)
	case float64:
		return p.handleNumber(v)
	case bool:
		return p.handleBoolean(v)
	case nil:
		return p.handleNull()
	default:
		return NewFormatError(fmt.Sprintf("unknown token type: %T", token))
	}
}

// handleDelimiter processes JSON delimiters (brackets and braces)
func (p *TokenParser) handleDelimiter(delim json.Delim) error {
	switch delim {
	case '{':
		return p.startObject()
	case '}':
		return p.endObject()
	case '[':
		return p.startArray()
	case ']':
		return p.endArray()
	default:
		return NewFormatError(fmt.Sprintf("unknown delimiter: %c", delim))
	}
}

// startObject handles the start of a JSON object
func (p *TokenParser) startObject() error {
	// Validate parser state
	if p.builder == nil {
		return NewFormatError("invalid parser state: builder is nil")
	}
	if p.config == nil {
		return NewFormatError("invalid parser state: config is nil")
	}

	// Validate state: we shouldn't be expecting a key when starting an object
	// unless we're at the root level
	if p.expectingKey && p.depth > 0 {
		return NewFormatError("malformed JSON: unexpected object start, expected object key")
	}

	// Validate depth limits to prevent stack overflow
	if p.depth >= 100 {
		return NewFormatError("JSON structure too deeply nested (max depth: 100)")
	}

	// Add comma if not the first element and we're in an array
	if !p.isFirstElement && p.isInArray() {
		if _, err := p.builder.WriteString(","); err != nil {
			return WrapFormatError("failed to write comma separator", err)
		}
		if p.shouldFormatCompact() {
			if _, err := p.builder.WriteString(" "); err != nil {
				return WrapFormatError("failed to write space", err)
			}
		} else {
			if err := p.writeNewlineAndIndent(); err != nil {
				return WrapFormatError("failed to write newline and indent", err)
			}
		}
	} else if p.depth > 0 && p.isInArray() {
		if p.shouldFormatCompact() {
			// For compact formatting, don't add newline
		} else {
			if err := p.writeNewlineAndIndent(); err != nil {
				return WrapFormatError("failed to write newline and indent", err)
			}
		}
	}

	// Write opening brace with space if it's a value after a key
	if p.depth > 0 && !p.isInArray() {
		// This is an object value, add space after colon
		if _, err := p.builder.WriteString(" {"); err != nil {
			return WrapFormatError("failed to write opening brace", err)
		}
	} else {
		if _, err := p.builder.WriteString("{"); err != nil {
			return WrapFormatError("failed to write opening brace", err)
		}
	}

	// Update parser state
	if err := p.enterObject(); err != nil {
		return WrapFormatError("failed to enter object state", err)
	}
	p.isFirstElement = true
	p.expectingKey = true

	return nil
}

// endObject handles the end of a JSON object
func (p *TokenParser) endObject() error {
	// Validate parser state
	if p.builder == nil {
		return NewFormatError("invalid parser state: builder is nil")
	}
	if p.config == nil {
		return NewFormatError("invalid parser state: config is nil")
	}

	// Validate state: we should be in an object context
	if p.depth == 0 {
		return NewFormatError("malformed JSON: unexpected object end, no matching opening brace")
	}
	if len(p.inArray) == 0 {
		return NewFormatError("malformed JSON: unexpected object end, invalid parser state")
	}

	// Validate that we're actually in an object (not array)
	if p.isInArray() {
		return NewFormatError("malformed JSON: unexpected object end, currently in array context")
	}

	// Check if this object should be formatted compactly BEFORE updating state
	isCompact := p.shouldFormatCompact()

	// Update parser state
	if err := p.exitObject(); err != nil {
		return WrapFormatError("failed to exit object state", err)
	}
	p.expectingKey = false

	// Format closing brace based on compact status
	if isCompact {
		// For compact objects, just add the closing brace without newline
		if _, err := p.builder.WriteString("}"); err != nil {
			return WrapFormatError("failed to write closing brace", err)
		}
	} else {
		// For normal objects, add newline and indentation before closing brace
		if err := p.writeNewlineAndIndent(); err != nil {
			return WrapFormatError("failed to write newline and indent", err)
		}
		if _, err := p.builder.WriteString("}"); err != nil {
			return WrapFormatError("failed to write closing brace", err)
		}
	}

	return nil
}

// startArray handles the start of a JSON array
func (p *TokenParser) startArray() error {
	// Validate parser state
	if p.builder == nil {
		return NewFormatError("invalid parser state: builder is nil")
	}
	if p.config == nil {
		return NewFormatError("invalid parser state: config is nil")
	}

	// Validate state: we shouldn't be expecting a key when starting an array
	// unless we're at the root level
	if p.expectingKey && p.depth > 0 {
		return NewFormatError("malformed JSON: unexpected array start, expected object key")
	}

	// Validate depth limits to prevent stack overflow
	if p.depth >= 100 {
		return NewFormatError("JSON structure too deeply nested (max depth: 100)")
	}

	// Add comma if not the first element and we're in an array
	if !p.isFirstElement && p.isInArray() {
		if _, err := p.builder.WriteString(","); err != nil {
			return WrapFormatError("failed to write comma separator", err)
		}
		if p.shouldFormatCompact() {
			if _, err := p.builder.WriteString(" "); err != nil {
				return WrapFormatError("failed to write space", err)
			}
		} else {
			if err := p.writeNewlineAndIndent(); err != nil {
				return WrapFormatError("failed to write newline and indent", err)
			}
		}
	} else if p.depth > 0 {
		// For arrays, don't add newline - they should start right after the colon
		// This will be handled by the key-value separator logic
	}

	// Write opening bracket with space if it's a value after a key
	if p.depth > 0 && !p.isInArray() {
		// This is an array value, add space after colon
		if _, err := p.builder.WriteString(" ["); err != nil {
			return WrapFormatError("failed to write opening bracket", err)
		}
	} else {
		if _, err := p.builder.WriteString("["); err != nil {
			return WrapFormatError("failed to write opening bracket", err)
		}
	}

	// Update parser state
	if err := p.enterArray(); err != nil {
		return WrapFormatError("failed to enter array state", err)
	}
	p.isFirstElement = true
	p.expectingKey = false

	return nil
}

// endArray handles the end of a JSON array
func (p *TokenParser) endArray() error {
	// Validate parser state
	if p.builder == nil {
		return NewFormatError("invalid parser state: builder is nil")
	}
	if p.config == nil {
		return NewFormatError("invalid parser state: config is nil")
	}

	// Validate state: we should be in an array context
	if p.depth == 0 {
		return NewFormatError("malformed JSON: unexpected array end, no matching opening bracket")
	}
	if len(p.inArray) == 0 {
		return NewFormatError("malformed JSON: unexpected array end, invalid parser state")
	}

	// Validate that we're actually in an array (not object)
	if !p.isInArray() {
		return NewFormatError("malformed JSON: unexpected array end, currently in object context")
	}

	// Check if this array should be formatted compactly BEFORE updating state
	isCompact := p.shouldFormatCompact()

	// Update parser state first
	if err := p.exitArray(); err != nil {
		return WrapFormatError("failed to exit array state", err)
	}

	// Format closing bracket based on compact status
	if isCompact {
		// For compact arrays, just add the closing bracket without newline
		if _, err := p.builder.WriteString("]"); err != nil {
			return WrapFormatError("failed to write closing bracket", err)
		}
	} else {
		// For normal arrays, add newline and indentation before closing bracket
		if err := p.writeNewlineAndIndent(); err != nil {
			return WrapFormatError("failed to write newline and indent", err)
		}
		if _, err := p.builder.WriteString("]"); err != nil {
			return WrapFormatError("failed to write closing bracket", err)
		}
	}

	// If we're back in an object after the array, next string will be a key
	if !p.isInArray() {
		p.expectingKey = true
	}

	return nil
}

// handleString handles string tokens (both keys and values)
func (p *TokenParser) handleString(value string) error {
	// Validate parser state
	if p.builder == nil {
		return NewFormatError("invalid parser state: builder is nil")
	}
	if p.config == nil {
		return NewFormatError("invalid parser state: config is nil")
	}

	// Validate string length to prevent memory issues
	if len(value) > 1000000 { // 1MB limit for individual strings
		return NewFormatError("string value too large (exceeds 1MB limit)")
	}

	// Check if this is an object key
	if p.expectingKey {
		// Validate that we're in an object context when expecting a key
		if p.depth == 0 || p.isInArray() {
			return NewFormatError("malformed JSON: unexpected object key outside of object context")
		}

		// Add comma if not the first element
		if !p.isFirstElement {
			if _, err := p.builder.WriteString(","); err != nil {
				return WrapFormatError("failed to write comma separator", err)
			}
			if p.shouldFormatCompact() {
				if _, err := p.builder.WriteString(" "); err != nil {
					return WrapFormatError("failed to write space", err)
				}
			} else {
				if err := p.writeNewlineAndIndent(); err != nil {
					return WrapFormatError("failed to write newline and indent", err)
				}
			}
		} else if p.depth > 0 && !p.shouldFormatCompact() {
			if err := p.writeNewlineAndIndent(); err != nil {
				return WrapFormatError("failed to write newline and indent", err)
			}
		}

		// Write the key with quotes and colon
		if _, err := p.builder.WriteString(`"`); err != nil {
			return WrapFormatError("failed to write opening quote for key", err)
		}
		escapedKey, err := p.escapeString(value)
		if err != nil {
			return WrapFormatError("failed to escape object key", err)
		}
		if _, err := p.builder.WriteString(escapedKey); err != nil {
			return WrapFormatError("failed to write object key", err)
		}
		if _, err := p.builder.WriteString(`":`); err != nil {
			return WrapFormatError("failed to write key-value separator", err)
		}

		// Mark that we've processed an element and now expect a value
		p.isFirstElement = false
		p.expectingKey = false
	} else {
		// This is a value (either in array or object value)
		// Only add comma if we're in an array and not the first element
		if !p.isFirstElement && p.isInArray() {
			if _, err := p.builder.WriteString(","); err != nil {
				return WrapFormatError("failed to write comma separator", err)
			}
			if p.shouldFormatCompact() {
				if _, err := p.builder.WriteString(" "); err != nil {
					return WrapFormatError("failed to write space", err)
				}
			} else {
				if err := p.writeNewlineAndIndent(); err != nil {
					return WrapFormatError("failed to write newline and indent", err)
				}
			}
		} else if p.depth > 0 && p.isInArray() {
			if p.shouldFormatCompact() {
				// For compact formatting, don't add newline
			} else {
				if err := p.writeNewlineAndIndent(); err != nil {
					return WrapFormatError("failed to write newline and indent", err)
				}
			}
		}

		// Write the JSON-escaped string with quotes, add space if it's a value after a key
		if p.depth > 0 && !p.isInArray() {
			// This is an object value, add space after colon
			if _, err := p.builder.WriteString(` "`); err != nil {
				return WrapFormatError("failed to write opening quote for string value", err)
			}
		} else {
			if _, err := p.builder.WriteString(`"`); err != nil {
				return WrapFormatError("failed to write opening quote for string value", err)
			}
		}
		escapedValue, err := p.escapeString(value)
		if err != nil {
			return WrapFormatError("failed to escape string value", err)
		}
		if _, err := p.builder.WriteString(escapedValue); err != nil {
			return WrapFormatError("failed to write string value", err)
		}
		if _, err := p.builder.WriteString(`"`); err != nil {
			return WrapFormatError("failed to write closing quote for string value", err)
		}

		// Mark that we've processed an element
		p.isFirstElement = false

		// If we're in an object, next string will be a key
		if !p.isInArray() {
			p.expectingKey = true
		}
	}

	return nil
}

// handleNumber handles numeric values
func (p *TokenParser) handleNumber(value float64) error {
	// Validate parser state
	if p.builder == nil {
		return NewFormatError("invalid parser state: builder is nil")
	}
	if p.config == nil {
		return NewFormatError("invalid parser state: config is nil")
	}

	// Validate that we're not expecting a key (numbers can't be object keys)
	if p.expectingKey {
		return NewFormatError("malformed JSON: unexpected number, expected object key")
	}

	// Validate number value for special cases
	if value != value { // NaN check
		return NewFormatError("invalid JSON: NaN values are not allowed")
	}
	if value == value+1 && value == value*2 { // Infinity check
		return NewFormatError("invalid JSON: infinite values are not allowed")
	}

	// Only add comma if we're in an array and not the first element
	if !p.isFirstElement && p.isInArray() {
		if _, err := p.builder.WriteString(","); err != nil {
			return WrapFormatError("failed to write comma separator", err)
		}
		if p.shouldFormatCompact() {
			if _, err := p.builder.WriteString(" "); err != nil {
				return WrapFormatError("failed to write space", err)
			}
		} else {
			if err := p.writeNewlineAndIndent(); err != nil {
				return WrapFormatError("failed to write newline and indent", err)
			}
		}
	} else if p.depth > 0 && p.isInArray() {
		if p.shouldFormatCompact() {
			// For compact formatting, don't add newline
		} else {
			if err := p.writeNewlineAndIndent(); err != nil {
				return WrapFormatError("failed to write newline and indent", err)
			}
		}
	}

	// Write the number value, add space if it's a value after a key
	formattedNumber, err := p.formatNumber(value)
	if err != nil {
		return WrapFormatError("failed to format number", err)
	}
	if p.depth > 0 && !p.isInArray() {
		// This is an object value, add space after colon
		if _, err := p.builder.WriteString(" " + formattedNumber); err != nil {
			return WrapFormatError("failed to write number value", err)
		}
	} else {
		if _, err := p.builder.WriteString(formattedNumber); err != nil {
			return WrapFormatError("failed to write number value", err)
		}
	}

	// Mark that we've processed an element
	p.isFirstElement = false

	// If we're in an object, next string will be a key
	if !p.isInArray() {
		p.expectingKey = true
	}

	return nil
}

// handleBoolean handles boolean values
func (p *TokenParser) handleBoolean(value bool) error {
	// Validate parser state
	if p.builder == nil {
		return NewFormatError("invalid parser state: builder is nil")
	}
	if p.config == nil {
		return NewFormatError("invalid parser state: config is nil")
	}

	// Validate that we're not expecting a key (booleans can't be object keys)
	if p.expectingKey {
		return NewFormatError("malformed JSON: unexpected boolean, expected object key")
	}

	// Only add comma if we're in an array and not the first element
	if !p.isFirstElement && p.isInArray() {
		if _, err := p.builder.WriteString(","); err != nil {
			return WrapFormatError("failed to write comma separator", err)
		}
		if p.shouldFormatCompact() {
			if _, err := p.builder.WriteString(" "); err != nil {
				return WrapFormatError("failed to write space", err)
			}
		} else {
			if err := p.writeNewlineAndIndent(); err != nil {
				return WrapFormatError("failed to write newline and indent", err)
			}
		}
	} else if p.depth > 0 && p.isInArray() {
		if p.shouldFormatCompact() {
			// For compact formatting, don't add newline
		} else {
			if err := p.writeNewlineAndIndent(); err != nil {
				return WrapFormatError("failed to write newline and indent", err)
			}
		}
	}

	// Write the boolean value, add space if it's a value after a key
	var boolStr string
	if value {
		boolStr = "true"
	} else {
		boolStr = "false"
	}
	if p.depth > 0 && !p.isInArray() {
		// This is an object value, add space after colon
		if _, err := p.builder.WriteString(" " + boolStr); err != nil {
			return WrapFormatError("failed to write boolean value", err)
		}
	} else {
		if _, err := p.builder.WriteString(boolStr); err != nil {
			return WrapFormatError("failed to write boolean value", err)
		}
	}

	// Mark that we've processed an element
	p.isFirstElement = false

	// If we're in an object, next string will be a key
	if !p.isInArray() {
		p.expectingKey = true
	}

	return nil
}

// handleNull handles null values
func (p *TokenParser) handleNull() error {
	// Validate parser state
	if p.builder == nil {
		return NewFormatError("invalid parser state: builder is nil")
	}
	if p.config == nil {
		return NewFormatError("invalid parser state: config is nil")
	}

	// Validate that we're not expecting a key (null can't be object keys)
	if p.expectingKey {
		return NewFormatError("malformed JSON: unexpected null, expected object key")
	}

	// Only add comma if we're in an array and not the first element
	if !p.isFirstElement && p.isInArray() {
		if _, err := p.builder.WriteString(","); err != nil {
			return WrapFormatError("failed to write comma separator", err)
		}
		if p.shouldFormatCompact() {
			if _, err := p.builder.WriteString(" "); err != nil {
				return WrapFormatError("failed to write space", err)
			}
		} else {
			if err := p.writeNewlineAndIndent(); err != nil {
				return WrapFormatError("failed to write newline and indent", err)
			}
		}
	} else if p.depth > 0 && p.isInArray() {
		if p.shouldFormatCompact() {
			// For compact formatting, don't add newline
		} else {
			if err := p.writeNewlineAndIndent(); err != nil {
				return WrapFormatError("failed to write newline and indent", err)
			}
		}
	}

	// Write null value, add space if it's a value after a key
	if p.depth > 0 && !p.isInArray() {
		// This is an object value, add space after colon
		if _, err := p.builder.WriteString(" null"); err != nil {
			return WrapFormatError("failed to write null value", err)
		}
	} else {
		if _, err := p.builder.WriteString("null"); err != nil {
			return WrapFormatError("failed to write null value", err)
		}
	}

	// Mark that we've processed an element
	p.isFirstElement = false

	// If we're in an object, next string will be a key
	if !p.isInArray() {
		p.expectingKey = true
	}

	return nil
}

// enterArray updates parser state when entering an array
func (p *TokenParser) enterArray() error {
	// Validate state before entering array
	if p.depth < 0 {
		return NewFormatError("invalid parser state: negative depth")
	}
	if p.depth >= 100 {
		return NewFormatError("JSON structure too deeply nested (max depth: 100)")
	}

	p.depth++
	p.inArray = append(p.inArray, true)
	return nil
}

// exitArray updates parser state when exiting an array
func (p *TokenParser) exitArray() error {
	// Validate state before exiting array
	if p.depth <= 0 {
		return NewFormatError("invalid parser state: cannot exit array, depth is zero or negative")
	}
	if len(p.inArray) == 0 {
		return NewFormatError("invalid parser state: cannot exit array, no array context")
	}
	if !p.inArray[len(p.inArray)-1] {
		return NewFormatError("invalid parser state: cannot exit array, currently in object context")
	}

	p.depth--
	p.inArray = p.inArray[:len(p.inArray)-1]
	return nil
}

// enterObject updates parser state when entering an object
func (p *TokenParser) enterObject() error {
	// Validate state before entering object
	if p.depth < 0 {
		return NewFormatError("invalid parser state: negative depth")
	}
	if p.depth >= 100 {
		return NewFormatError("JSON structure too deeply nested (max depth: 100)")
	}

	p.depth++
	p.inArray = append(p.inArray, false)
	return nil
}

// exitObject updates parser state when exiting an object
func (p *TokenParser) exitObject() error {
	// Validate state before exiting object
	if p.depth <= 0 {
		return NewFormatError("invalid parser state: cannot exit object, depth is zero or negative")
	}
	if len(p.inArray) == 0 {
		return NewFormatError("invalid parser state: cannot exit object, no object context")
	}
	if p.inArray[len(p.inArray)-1] {
		return NewFormatError("invalid parser state: cannot exit object, currently in array context")
	}

	p.depth--
	p.inArray = p.inArray[:len(p.inArray)-1]
	return nil
}

// isInArray returns true if currently inside an array
func (p *TokenParser) isInArray() bool {
	if len(p.inArray) == 0 {
		return false
	}
	return p.inArray[len(p.inArray)-1]
}

// shouldFormatCompact determines if elements at current depth should be formatted compactly
func (p *TokenParser) shouldFormatCompact() bool {
	// Format compactly if we're at or beyond the configured compact depth
	return p.config.CompactDepth > 0 && p.depth >= p.config.CompactDepth
}

// writeIndent writes the appropriate indentation based on current depth and config
func (p *TokenParser) writeIndent() error {
	// Validate parser state
	if p.builder == nil {
		return NewFormatError("invalid parser state: builder is nil")
	}
	if p.config == nil {
		return NewFormatError("invalid parser state: config is nil")
	}
	if p.depth < 0 {
		return NewFormatError("invalid parser state: negative depth")
	}

	var indentStr string
	if p.config.UseTab {
		indentStr = strings.Repeat("\t", p.depth)
	} else {
		// Validate indent size to prevent excessive memory usage
		totalSpaces := p.depth * p.config.IndentSize
		if totalSpaces > 10000 { // Limit total indentation to prevent memory issues
			return NewFormatError("indentation too large (exceeds 10000 characters)")
		}
		indentStr = strings.Repeat(" ", totalSpaces)
	}

	if _, err := p.builder.WriteString(indentStr); err != nil {
		return WrapFormatError("failed to write indentation", err)
	}

	return nil
}

// writeNewlineAndIndent writes a newline followed by proper indentation
func (p *TokenParser) writeNewlineAndIndent() error {
	// Validate parser state
	if p.builder == nil {
		return NewFormatError("invalid parser state: builder is nil")
	}

	if _, err := p.builder.WriteString("\n"); err != nil {
		return WrapFormatError("failed to write newline", err)
	}

	if err := p.writeIndent(); err != nil {
		return WrapFormatError("failed to write indentation after newline", err)
	}

	return nil
}

// escapeString properly escapes a string for JSON output
func (p *TokenParser) escapeString(s string) (string, error) {
	// Validate input string
	if len(s) > 1000000 { // 1MB limit for individual strings
		return "", NewFormatError("string too large for escaping (exceeds 1MB limit)")
	}

	// Use json.Marshal to properly escape the string, then remove the surrounding quotes
	escaped, err := json.Marshal(s)
	if err != nil {
		return "", WrapFormatError("failed to escape string for JSON output", err)
	}

	// Remove the surrounding quotes that json.Marshal adds
	escapedStr := string(escaped)
	if len(escapedStr) >= 2 && escapedStr[0] == '"' && escapedStr[len(escapedStr)-1] == '"' {
		return escapedStr[1 : len(escapedStr)-1], nil
	}

	// If the escaped string doesn't have quotes (shouldn't happen), return as-is
	return escapedStr, nil
}

// formatNumber formats a float64 number for JSON output
func (p *TokenParser) formatNumber(value float64) (string, error) {
	// Validate number value
	if value != value { // NaN check
		return "", NewFormatError("cannot format NaN as JSON number")
	}
	if value == value+1 && value == value*2 { // Infinity check
		return "", NewFormatError("cannot format infinite value as JSON number")
	}

	// Use json.Marshal to properly format the number
	formatted, err := json.Marshal(value)
	if err != nil {
		return "", WrapFormatError("failed to format number for JSON output", err)
	}

	return string(formatted), nil
}

// calculatePosition estimates the current position in the input stream
func (p *TokenParser) calculatePosition(reader *strings.Reader) int {
	// Get the current position by checking how much has been read
	currentPos := p.inputLength - reader.Len()
	if currentPos < 0 {
		return 0
	}
	if currentPos > p.inputLength {
		return p.inputLength
	}
	return currentPos
}

// FormatError represents an error that occurred during JSON formatting.
// It provides detailed information about what went wrong, including
// the position in the input where the error occurred and the underlying cause.
type FormatError struct {
	// Msg contains a human-readable description of what went wrong
	Msg string

	// Position indicates the approximate character position in the input
	// where the error occurred (0-based). A value of 0 means position
	// information is not available.
	Position int

	// Original contains the underlying error that caused this formatting error.
	// It may be nil if the error originated within the formatter itself.
	Original error
}

// Error implements the error interface and returns a formatted error message.
// The message includes position information when available and details
// about any underlying error.
func (e *FormatError) Error() string {
	if e.Position > 0 {
		if e.Original != nil {
			return fmt.Sprintf("%s at position %d: %v", e.Msg, e.Position, e.Original)
		}
		return fmt.Sprintf("%s at position %d", e.Msg, e.Position)
	}

	if e.Original != nil {
		return fmt.Sprintf("%s: %v", e.Msg, e.Original)
	}
	return e.Msg
}

// Unwrap returns the underlying error for error unwrapping.
// This allows the use of errors.Is() and errors.As() with FormatError.
// Returns nil if there is no underlying error.
func (e *FormatError) Unwrap() error {
	return e.Original
}

// NewFormatError creates a new FormatError with the given message.
// The position is set to 0 (unknown) and there is no underlying error.
func NewFormatError(msg string) *FormatError {
	return &FormatError{
		Msg:      msg,
		Position: 0,
		Original: nil,
	}
}

// NewFormatErrorWithPosition creates a new FormatError with message and position.
// The position should be a 0-based character index in the input where the error occurred.
func NewFormatErrorWithPosition(msg string, position int) *FormatError {
	return &FormatError{
		Msg:      msg,
		Position: position,
		Original: nil,
	}
}

// WrapFormatError wraps an existing error with formatting context.
// This is useful for adding context to errors from underlying libraries
// while preserving the original error for unwrapping.
func WrapFormatError(msg string, err error) *FormatError {
	return &FormatError{
		Msg:      msg,
		Position: 0,
		Original: err,
	}
}

// WrapFormatErrorWithPosition wraps an existing error with formatting context and position.
// This provides the most detailed error information, including both position
// and the underlying cause of the error.
func WrapFormatErrorWithPosition(msg string, position int, err error) *FormatError {
	return &FormatError{
		Msg:      msg,
		Position: position,
		Original: err,
	}
}
