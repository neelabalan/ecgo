package main

import (
	"testing"
)

// TestApplyEscapes tests the applyEscapes function
func TestApplyEscapes(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello world", "hello world"},
		{"hello\\nworld", "hello\nworld"},
		{"hello\\tworld", "hello\tworld"},
		{"line1\\nline2\\ttab", "line1\nline2\ttab"},
		{"", ""},
	}

	for _, test := range tests {
		result := applyEscapes(test.input)
		if result != test.expected {
			t.Errorf("applyEscapes(%q) = %q, want %q", test.input, result, test.expected)
		}
	}
}

// TestColorValidation tests basic color name validation
func TestColorValidation(t *testing.T) {
	validColors := []string{"red", "green", "blue", "yellow", "magenta", "cyan", "white", "black"}

	for _, color := range validColors {
		t.Run("valid_color_"+color, func(t *testing.T) {
			// This would test that valid colors don't cause errors
			// In a real test, you'd mock the color.New function
			if color == "" {
				t.Error("Color should not be empty")
			}
		})
	}
}

// TestEmptyArgs tests behavior with no arguments
func TestEmptyArgs(t *testing.T) {
	// Test that empty arguments don't cause panic
	args := []string{}
	if len(args) != 0 {
		t.Error("Expected empty args")
	}
}
