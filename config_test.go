package jsonformat

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config == nil {
		t.Fatal("DefaultConfig() returned nil")
	}

	if config.IndentSize != 2 {
		t.Errorf("Expected IndentSize to be 2, got %d", config.IndentSize)
	}

	if config.UseTab != false {
		t.Errorf("Expected UseTab to be false, got %t", config.UseTab)
	}

	if config.CompactDepth != 3 {
		t.Errorf("Expected CompactDepth to be 3, got %d", config.CompactDepth)
	}
}

func TestDefaultConfigReturnsNewInstance(t *testing.T) {
	config1 := DefaultConfig()
	config2 := DefaultConfig()

	if config1 == config2 {
		t.Error("DefaultConfig() should return new instances, not the same pointer")
	}

	// Modify one config to ensure they're independent
	config1.IndentSize = 4
	if config2.IndentSize != 2 {
		t.Error("Modifying one config instance affected another")
	}
}

func TestNewConfigWithOptions(t *testing.T) {
	tests := []struct {
		name     string
		options  []ConfigOption
		expected *Config
	}{
		{
			name:    "no options",
			options: []ConfigOption{},
			expected: &Config{
				IndentSize:   2,
				UseTab:       false,
				CompactDepth: 3,
			},
		},
		{
			name:    "with indent size",
			options: []ConfigOption{WithIndentSize(4)},
			expected: &Config{
				IndentSize:   4,
				UseTab:       false,
				CompactDepth: 3,
			},
		},
		{
			name:    "with tabs",
			options: []ConfigOption{WithTabs()},
			expected: &Config{
				IndentSize:   2,
				UseTab:       true,
				CompactDepth: 3,
			},
		},
		{
			name:    "with spaces",
			options: []ConfigOption{WithSpaces()},
			expected: &Config{
				IndentSize:   2,
				UseTab:       false,
				CompactDepth: 3,
			},
		},
		{
			name:    "set compact depth to 0",
			options: []ConfigOption{WithCompactDepth(0)},
			expected: &Config{
				IndentSize:   2,
				UseTab:       false,
				CompactDepth: 0,
			},
		},
		{
			name: "multiple options",
			options: []ConfigOption{
				WithIndentSize(8),
				WithTabs(),
				WithCompactDepth(0),
			},
			expected: &Config{
				IndentSize:   8,
				UseTab:       true,
				CompactDepth: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := NewConfig(tt.options...)

			if config.IndentSize != tt.expected.IndentSize {
				t.Errorf("Expected IndentSize %d, got %d", tt.expected.IndentSize, config.IndentSize)
			}

			if config.UseTab != tt.expected.UseTab {
				t.Errorf("Expected UseTab %t, got %t", tt.expected.UseTab, config.UseTab)
			}

			if config.CompactDepth != tt.expected.CompactDepth {
				t.Errorf("Expected CompactDepth %d, got %d", tt.expected.CompactDepth, config.CompactDepth)
			}
		})
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		options     []ConfigOption
		expectValid bool
	}{
		{
			name:        "valid config with default values",
			options:     []ConfigOption{},
			expectValid: true,
		},
		{
			name:        "valid config with custom indent size",
			options:     []ConfigOption{WithIndentSize(4)},
			expectValid: true,
		},
		{
			name:        "valid config with zero indent size",
			options:     []ConfigOption{WithIndentSize(0)},
			expectValid: true,
		},
		{
			name:        "valid config with max indent size",
			options:     []ConfigOption{WithIndentSize(20)},
			expectValid: true,
		},
		{
			name:        "invalid config with negative indent size",
			options:     []ConfigOption{WithIndentSize(-1)},
			expectValid: true, // WithIndentSize should ignore invalid values
		},
		{
			name:        "invalid config with too large indent size",
			options:     []ConfigOption{WithIndentSize(25)},
			expectValid: true, // WithIndentSize should ignore invalid values
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := NewConfig(tt.options...)

			// All configs should be valid because NewConfig handles validation
			// and returns default config if validation fails
			if config == nil {
				t.Error("NewConfig returned nil")
			}

			// Verify that invalid indent sizes are ignored
			if config.IndentSize < 0 || config.IndentSize > 20 {
				t.Errorf("Config has invalid IndentSize: %d", config.IndentSize)
			}
		})
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
	}{
		{
			name:        "nil config",
			config:      nil,
			expectError: true,
		},
		{
			name: "valid config",
			config: &Config{
				IndentSize:   2,
				UseTab:       false,
				CompactDepth: 3,
			},
			expectError: false,
		},
		{
			name: "negative indent size",
			config: &Config{
				IndentSize:   -1,
				UseTab:       false,
				CompactDepth: 3,
			},
			expectError: true,
		},
		{
			name: "too large indent size",
			config: &Config{
				IndentSize:   25,
				UseTab:       false,
				CompactDepth: 3,
			},
			expectError: true,
		},
		{
			name: "zero indent size",
			config: &Config{
				IndentSize:   0,
				UseTab:       false,
				CompactDepth: 3,
			},
			expectError: false,
		},
		{
			name: "max indent size",
			config: &Config{
				IndentSize:   20,
				UseTab:       false,
				CompactDepth: 3,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.config)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestConfigOptions(t *testing.T) {
	t.Run("WithIndentSize", func(t *testing.T) {
		config := &Config{IndentSize: 2}

		// Valid size
		WithIndentSize(4)(config)
		if config.IndentSize != 4 {
			t.Errorf("Expected IndentSize 4, got %d", config.IndentSize)
		}

		// Invalid negative size - should not change
		WithIndentSize(-1)(config)
		if config.IndentSize != 4 {
			t.Errorf("Expected IndentSize to remain 4, got %d", config.IndentSize)
		}

		// Invalid too large size - should not change
		WithIndentSize(25)(config)
		if config.IndentSize != 4 {
			t.Errorf("Expected IndentSize to remain 4, got %d", config.IndentSize)
		}

		// Edge cases
		WithIndentSize(0)(config)
		if config.IndentSize != 0 {
			t.Errorf("Expected IndentSize 0, got %d", config.IndentSize)
		}

		WithIndentSize(20)(config)
		if config.IndentSize != 20 {
			t.Errorf("Expected IndentSize 20, got %d", config.IndentSize)
		}
	})

	t.Run("WithTabs", func(t *testing.T) {
		config := &Config{UseTab: false}
		WithTabs()(config)
		if !config.UseTab {
			t.Error("Expected UseTab to be true")
		}
	})

	t.Run("WithSpaces", func(t *testing.T) {
		config := &Config{UseTab: true}
		WithSpaces()(config)
		if config.UseTab {
			t.Error("Expected UseTab to be false")
		}
	})

	t.Run("WithCompactDepth", func(t *testing.T) {
		config := &Config{CompactDepth: 3}
		WithCompactDepth(0)(config)
		if config.CompactDepth != 0 {
			t.Error("Expected CompactDepth to be 0")
		}

		WithCompactDepth(2)(config)
		if config.CompactDepth != 2 {
			t.Error("Expected CompactDepth to be 2")
		}
	})
}

func TestConfigOptionsCombinations(t *testing.T) {
	tests := []struct {
		name     string
		options  []ConfigOption
		expected *Config
	}{
		{
			name: "all custom options",
			options: []ConfigOption{
				WithIndentSize(8),
				WithTabs(),
				WithCompactDepth(0),
			},
			expected: &Config{
				IndentSize:   8,
				UseTab:       true,
				CompactDepth: 0,
			},
		},
		{
			name: "conflicting tab options - last wins",
			options: []ConfigOption{
				WithTabs(),
				WithSpaces(),
			},
			expected: &Config{
				IndentSize:   2,
				UseTab:       false,
				CompactDepth: 3,
			},
		},
		{
			name: "multiple indent size options - last wins",
			options: []ConfigOption{
				WithIndentSize(4),
				WithIndentSize(6),
				WithIndentSize(8),
			},
			expected: &Config{
				IndentSize:   8,
				UseTab:       false,
				CompactDepth: 3,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := NewConfig(tt.options...)

			if config.IndentSize != tt.expected.IndentSize {
				t.Errorf("Expected IndentSize %d, got %d", tt.expected.IndentSize, config.IndentSize)
			}

			if config.UseTab != tt.expected.UseTab {
				t.Errorf("Expected UseTab %t, got %t", tt.expected.UseTab, config.UseTab)
			}

			if config.CompactDepth != tt.expected.CompactDepth {
				t.Errorf("Expected CompactDepth %d, got %d", tt.expected.CompactDepth, config.CompactDepth)
			}
		})
	}
}

func TestNewConfigWithInvalidOptions(t *testing.T) {
	// Test that NewConfig handles invalid options gracefully
	config := NewConfig(
		WithIndentSize(-5),  // Invalid
		WithIndentSize(100), // Invalid
		WithIndentSize(4),   // Valid - should be used
		WithTabs(),
	)

	// Should use the valid indent size
	if config.IndentSize != 4 {
		t.Errorf("Expected IndentSize 4, got %d", config.IndentSize)
	}

	if !config.UseTab {
		t.Error("Expected UseTab to be true")
	}
}

func TestConfigValidationEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
		errorMsg    string
	}{
		{
			name:        "nil config",
			config:      nil,
			expectError: true,
			errorMsg:    "config cannot be nil",
		},
		{
			name: "boundary valid - zero indent",
			config: &Config{
				IndentSize:   0,
				UseTab:       false,
				CompactDepth: 3,
			},
			expectError: false,
		},
		{
			name: "boundary valid - max indent",
			config: &Config{
				IndentSize:   20,
				UseTab:       false,
				CompactDepth: 3,
			},
			expectError: false,
		},
		{
			name: "boundary invalid - negative indent",
			config: &Config{
				IndentSize:   -1,
				UseTab:       false,
				CompactDepth: 3,
			},
			expectError: true,
			errorMsg:    "IndentSize must be non-negative",
		},
		{
			name: "boundary invalid - too large indent",
			config: &Config{
				IndentSize:   21,
				UseTab:       false,
				CompactDepth: 3,
			},
			expectError: true,
			errorMsg:    "IndentSize must not exceed 20 spaces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.config)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestNewConfigValidationFallback(t *testing.T) {
	// Test that NewConfig falls back to default when validation fails
	// We can't directly test this since validateConfig is called internally,
	// but we can test the behavior indirectly

	// Create a config that would fail validation if passed directly
	config := NewConfig(WithIndentSize(-1)) // Invalid option should be ignored

	// Should get default config since invalid options are ignored by the option functions
	expected := DefaultConfig()
	if config.IndentSize != expected.IndentSize {
		t.Errorf("Expected fallback to default IndentSize %d, got %d", expected.IndentSize, config.IndentSize)
	}
	if config.UseTab != expected.UseTab {
		t.Errorf("Expected fallback to default UseTab %t, got %t", expected.UseTab, config.UseTab)
	}
	if config.CompactDepth != expected.CompactDepth {
		t.Errorf("Expected fallback to default CompactDepth %d, got %d", expected.CompactDepth, config.CompactDepth)
	}
}
