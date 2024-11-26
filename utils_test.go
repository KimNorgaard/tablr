package tablr

import (
	"testing"
)

func TestPad(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		width    int
		align    Alignment
		expected string
	}{
		{
			name:     "AlignLeft with padding",
			input:    "test",
			width:    10,
			align:    AlignLeft,
			expected: "test      ",
		},
		{
			name:     "AlignLeft without padding",
			input:    "test",
			width:    4,
			align:    AlignLeft,
			expected: "test",
		},
		{
			name:     "AlignCenter with padding",
			input:    "test",
			width:    10,
			align:    AlignCenter,
			expected: "   test   ",
		},
		{
			name:     "AlignCenter without padding",
			input:    "test",
			width:    4,
			align:    AlignCenter,
			expected: "test",
		},
		{
			name:     "AlignRight with padding",
			input:    "test",
			width:    10,
			align:    AlignRight,
			expected: "      test",
		},
		{
			name:     "AlignRight without padding",
			input:    "test",
			width:    4,
			align:    AlignRight,
			expected: "test",
		},
		{
			name:     "AlignDefault with padding",
			input:    "test",
			width:    10,
			align:    AlignDefault,
			expected: "test      ",
		},
		{
			name:     "AlignDefault without padding",
			input:    "test",
			width:    4,
			align:    AlignDefault,
			expected: "test",
		},
		{
			name:     "AlignLeft with exact width",
			input:    "test",
			width:    4,
			align:    AlignLeft,
			expected: "test",
		},
		{
			name:     "AlignCenter with exact width",
			input:    "test",
			width:    4,
			align:    AlignCenter,
			expected: "test",
		},
		{
			name:     "AlignRight with exact width",
			input:    "test",
			width:    4,
			align:    AlignRight,
			expected: "test",
		},
		{
			name:     "AlignDefault with exact width",
			input:    "test",
			width:    4,
			align:    AlignDefault,
			expected: "test",
		},
		{
			name:     "AlignDefault with unknown alignment",
			input:    "test",
			width:    10,
			align:    Alignment(10),
			expected: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pad(tt.input, tt.width, tt.align)
			if result != tt.expected {
				t.Errorf("pad(%q, %d, %v) = %q, want %q", tt.input, tt.width, tt.align, result, tt.expected)
			}
		})
	}
}
