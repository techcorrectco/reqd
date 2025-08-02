package commands

import "testing"

func TestGenerateIDPrefix(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "multi-word with spaces",
			input:    "Product Requirements Document",
			expected: "PRD",
		},
		{
			name:     "multi-word with hyphens",
			input:    "my-awesome-project",
			expected: "MAP",
		},
		{
			name:     "multi-word with underscores",
			input:    "user_management_system",
			expected: "UMS",
		},
		{
			name:     "multi-word with dots",
			input:    "api.gateway.service",
			expected: "AGS",
		},
		{
			name:     "mixed delimiters",
			input:    "web-app_version.2",
			expected: "WAV",
		},
		{
			name:     "single word",
			input:    "reqd",
			expected: "REQD",
		},
		{
			name:     "single word longer",
			input:    "application",
			expected: "APPL",
		},
		{
			name:     "single word short",
			input:    "go",
			expected: "GO",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "whitespace only",
			input:    "   ",
			expected: "",
		},
		{
			name:     "extra whitespace",
			input:    "  web   app  service  ",
			expected: "WAS",
		},
		{
			name:     "numbers and letters",
			input:    "version-2-api",
			expected: "VA",
		},
		{
			name:     "numbers only",
			input:    "123",
			expected: "123",
		},
		{
			name:     "special characters",
			input:    "test@#$",
			expected: "TEST",
		},
		{
			name:     "lowercase input",
			input:    "lower case project",
			expected: "LCP",
		},
		{
			name:     "uppercase input",
			input:    "UPPER CASE PROJECT",
			expected: "UCP",
		},
		{
			name:     "mixed case input",
			input:    "MiXeD CaSe PrOjEcT",
			expected: "MCP",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateIDPrefix(tt.input)
			if result != tt.expected {
				t.Errorf("generateIDPrefix(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}