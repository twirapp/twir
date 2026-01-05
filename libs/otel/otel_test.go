package otel

import (
	"reflect"
	"testing"
)

func TestParseHeaders(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: nil,
		},
		{
			name:  "single header",
			input: "authorization=Bearer token",
			expected: map[string]string{
				"authorization": "Bearer token",
			},
		},
		{
			name:  "multiple headers",
			input: "authorization=Bearer token,x-api-key=secret123",
			expected: map[string]string{
				"authorization": "Bearer token",
				"x-api-key":     "secret123",
			},
		},
		{
			name:  "headers with spaces",
			input: "key1 = value1 , key2 = value2",
			expected: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name:  "header with equals in value",
			input: "authorization=Basic dXNlcjpwYXNz==",
			expected: map[string]string{
				"authorization": "Basic dXNlcjpwYXNz==",
			},
		},
		{
			name:  "uptrace dsn header",
			input: "uptrace-dsn=https://project@uptrace.dev",
			expected: map[string]string{
				"uptrace-dsn": "https://project@uptrace.dev",
			},
		},
		{
			name:  "honeycomb headers",
			input: "x-honeycomb-team=abcd1234,x-honeycomb-dataset=my-service",
			expected: map[string]string{
				"x-honeycomb-team":    "abcd1234",
				"x-honeycomb-dataset": "my-service",
			},
		},
		{
			name:     "invalid format - no equals",
			input:    "invalid-header",
			expected: map[string]string{},
		},
		{
			name:  "mixed valid and invalid",
			input: "valid=value,invalid,another=good",
			expected: map[string]string{
				"valid":   "value",
				"another": "good",
			},
		},
		{
			name:     "only commas",
			input:    ",,,",
			expected: map[string]string{},
		},
		{
			name:  "trailing comma",
			input: "key=value,",
			expected: map[string]string{
				"key": "value",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseHeaders(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("parseHeaders(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCleanEndpoint(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "plain endpoint",
			input:    "localhost:4317",
			expected: "localhost:4317",
		},
		{
			name:     "http prefix",
			input:    "http://localhost:4317",
			expected: "localhost:4317",
		},
		{
			name:     "https prefix",
			input:    "https://localhost:4317",
			expected: "localhost:4317",
		},
		{
			name:     "grpc prefix",
			input:    "grpc://localhost:4317",
			expected: "localhost:4317",
		},
		{
			name:     "with domain",
			input:    "https://otlp.uptrace.dev:4317",
			expected: "otlp.uptrace.dev:4317",
		},
		{
			name:     "without port",
			input:    "https://api.honeycomb.io",
			expected: "api.honeycomb.io",
		},
		{
			name:     "with path (should keep)",
			input:    "http://localhost:4317/v1/traces",
			expected: "localhost:4317/v1/traces",
		},
		{
			name:     "with spaces",
			input:    "  http://localhost:4317  ",
			expected: "localhost:4317",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only spaces",
			input:    "   ",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanEndpoint(tt.input)
			if result != tt.expected {
				t.Errorf("cleanEndpoint(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
