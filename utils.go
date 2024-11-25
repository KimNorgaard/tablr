package tablr

import "strings"

// pad pads a string to the given width with spaces, aligning it as specified.
func pad(s string, width int, align Alignment) string {
	if len(s) >= width {
		return s
	}

	padding := width - len(s)
	switch align {
	case AlignLeft, AlignDefault:
		return s + strings.Repeat(" ", padding)
	case AlignCenter:
		left := padding / 2
		right := padding - left
		return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
	case AlignRight:
		return strings.Repeat(" ", padding) + s
	default:
		return s
	}
}

// escapePipes escapes pipe characters in a string with a backslash.
func escapePipes(s string) string {
	return strings.ReplaceAll(s, "|", "\\|")
}
