package tui

import (
	"testing"
)

func TestViewport_GetContent_StripsAnsi(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "plain text",
			content:  "Hello, World!",
			expected: "Hello, World!",
		},
		{
			name:     "text with ANSI color codes",
			content:  "\x1b[31mRed Text\x1b[0m Normal Text",
			expected: "Red Text Normal Text",
		},
		{
			name:     "text with multiple ANSI codes",
			content:  "\x1b[1;32mBold Green\x1b[0m \x1b[34mBlue\x1b[0m",
			expected: "Bold Green Blue",
		},
		{
			name:     "empty content",
			content:  "",
			expected: "",
		},
		{
			name:     "terraform-like output with colors",
			content:  "\x1b[32m+ resource \"aws_instance\"\x1b[0m will be created",
			expected: "+ resource \"aws_instance\" will be created",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewViewport(ViewportOptions{
				Width:  80,
				Height: 24,
			})
			v.content = []byte(tt.content)

			result := v.GetContent()
			if result != tt.expected {
				t.Errorf("GetContent() = %q, want %q", result, tt.expected)
			}
		})
	}
}
