# Go JSON Formatter

A Go library for formatting JSON with specialized formatting rules. The library provides custom indentation for JSON structures, with the unique feature of formatting objects within arrays on a single line for improved readability of array-heavy JSON data.

## Features

- **Custom Indentation**: Configure spaces or tabs with customizable indent size
- **Compact Depth Formatting**: Elements at specified depth or deeper are formatted on a single line
- **Flexible Configuration**: Use functional options pattern for easy customization
- **Robust Error Handling**: Comprehensive error reporting with position information
- **Memory Efficient**: Streaming token-based parsing for large JSON files
- **Go Conventions**: Follows standard Go practices and idioms

## Installation

```bash
go get go-json-formatter
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    formatter "go-json-formatter"
)

func main() {
    // JSON with nested structures
    jsonStr := `{"users":[{"id":1,"name":"Alice"},{"id":2,"name":"Bob"}],"meta":{"count":2}}`
    
    // Create formatter with default configuration
    config := formatter.DefaultConfig()
    f := formatter.NewFormatter(config)
    
    // Format the JSON
    formatted, err := f.Format(jsonStr)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(formatted)
}
```

**Output:**
```json
{
  "users": [
    {"id": 1, "name": "Alice"},
    {"id": 2, "name": "Bob"}
  ],
  "meta": {
    "count": 2
  }
}
```

## Configuration Options

### Default Configuration

The default configuration provides sensible defaults:
- **IndentSize**: 2 spaces
- **UseTab**: false (uses spaces)
- **CompactDepth**: 3 (elements at depth 3+ formatted compactly)

```go
config := formatter.DefaultConfig()
```

### Custom Configuration

Use functional options to customize the formatter:

```go
config := formatter.NewConfig(
    formatter.WithIndentSize(4),                    // Use 4 spaces
    formatter.WithTabs(),                           // Use tabs instead of spaces
    formatter.WithCompactDepth(2),                  // Compact at depth 2 and deeper
)
```

### Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `WithIndentSize(n)` | Set number of spaces for indentation (0-20) | 2 |
| `WithTabs()` | Use tabs instead of spaces | false |
| `WithSpaces()` | Use spaces instead of tabs | true |
| `WithCompactDepth(n)` | Set depth for compact formatting (0 disables) | 3 |

## Usage Examples

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    
    formatter "go-json-formatter"
)

func main() {
    jsonStr := `{"name":"Alice","items":[{"id":1,"type":"book"},{"id":2,"type":"pen"}]}`
    
    config := formatter.DefaultConfig()
    f := formatter.NewFormatter(config)
    
    formatted, err := f.Format(jsonStr)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(formatted)
}
```

### Custom Indentation

```go
// 4-space indentation
config := formatter.NewConfig(formatter.WithIndentSize(4))
f := formatter.NewFormatter(config)

// Tab indentation
config := formatter.NewConfig(formatter.WithTabs())
f := formatter.NewFormatter(config)
```

### Compact Depth Configuration

```go
// Compact at depth 2 and deeper
config := formatter.NewConfig(
    formatter.WithCompactDepth(2),
)
f := formatter.NewFormatter(config)

// Disable compact formatting entirely
config := formatter.NewConfig(
    formatter.WithCompactDepth(0),
)
f := formatter.NewFormatter(config)
```

### Working with Bytes

```go
jsonBytes := []byte(`{"users":[{"id":1,"name":"Alice"}]}`)

config := formatter.DefaultConfig()
f := formatter.NewFormatter(config)

formattedBytes, err := f.FormatBytes(jsonBytes)
if err != nil {
    log.Fatal(err)
}

fmt.Println(string(formattedBytes))
```

## Error Handling

The library provides detailed error information through the `FormatError` type:

```go
formatted, err := f.Format(invalidJSON)
if err != nil {
    var formatErr *formatter.FormatError
    if errors.As(err, &formatErr) {
        fmt.Printf("Error: %s\n", formatErr.Error())
        fmt.Printf("Position: %d\n", formatErr.Position)
        
        if formatErr.Unwrap() != nil {
            fmt.Printf("Underlying error: %v\n", formatErr.Unwrap())
        }
    }
}
```

### Common Error Scenarios

- **Empty Input**: Returns error for empty JSON strings
- **Invalid JSON**: Provides position information for syntax errors
- **Malformed Structure**: Detects unclosed objects/arrays
- **Deep Nesting**: Prevents stack overflow with depth limits (max: 100)
- **Large Strings**: Handles memory efficiently with size limits

## Formatting Behavior

### Normal Objects

Objects at the top level or nested within other objects use standard multi-line formatting:

```json
{
  "user": {
    "id": 1,
    "name": "Alice",
    "email": "alice@example.com"
  }
}
```

### Objects in Arrays

Objects that are direct children of arrays are formatted on a single line (when enabled):

```json
{
  "users": [
    {"id": 1, "name": "Alice", "email": "alice@example.com"},
    {"id": 2, "name": "Bob", "email": "bob@example.com"}
  ]
}
```

### Nested Structures

The formatter handles complex nested structures correctly:

```json
{
  "data": [
    {
      "type": "user",
      "attributes": {"name": "Alice", "age": 30},
      "relationships": {
        "posts": [
          {"id": 1, "title": "Hello World"},
          {"id": 2, "title": "Go Programming"}
        ]
      }
    }
  ]
}
```

## Performance

The library uses a streaming token-based approach for efficient memory usage:

- **Memory Efficient**: Processes JSON tokens sequentially without loading entire structure
- **Large File Support**: Handles large JSON files without excessive memory usage
- **Depth Limits**: Prevents stack overflow with configurable depth limits
- **String Limits**: Protects against memory exhaustion with string size limits

## API Reference

### Types

#### `Config`
Configuration struct for the formatter.

#### `ConfigOption`
Functional option type for configuring the formatter.

#### `Formatter`
Main formatter struct that handles JSON formatting.

#### `FormatError`
Error type that provides detailed formatting error information.

### Functions

#### `DefaultConfig() *Config`
Returns a configuration with default values.

#### `NewConfig(options ...ConfigOption) *Config`
Creates a new configuration with the provided options.

#### `NewFormatter(config *Config) *Formatter`
Creates a new formatter with the given configuration.

### Methods

#### `(f *Formatter) Format(jsonStr string) (string, error)`
Formats a JSON string according to the configured rules.

#### `(f *Formatter) FormatBytes(jsonBytes []byte) ([]byte, error)`
Formats JSON bytes according to the configured rules.

#### `(e *FormatError) Error() string`
Returns a formatted error message.

#### `(e *FormatError) Unwrap() error`
Returns the underlying error for error unwrapping.

### Configuration Options

#### `WithIndentSize(size int) ConfigOption`
Sets the number of spaces for indentation (0-20).

#### `WithTabs() ConfigOption`
Enables tab indentation.

#### `WithSpaces() ConfigOption`
Enables space indentation (default).

#### `WithCompactDepth(depth int) ConfigOption`
Sets the depth at which elements should be formatted compactly on a single line.

## Examples

See the `examples/` directory for complete working examples:

- `examples/basic_usage/` - Basic formatting example
- `examples/custom_config/` - Configuration options demonstration  
- `examples/error_handling/` - Error handling patterns

To run an example:
```bash
cd examples/basic_usage
go run main.go
```

## Requirements

- Go 1.18 or later
- No external dependencies (uses only standard library)

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## Changelog

### v1.0.0
- Initial release
- Basic JSON formatting with custom indentation
- Single-line array object formatting
- Comprehensive error handling
- Full test coverage