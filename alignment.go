package tablr

// Alignment represents the alignment of a column in a Markdown table.
type Alignment uint8

const (
	// AlignDefault aligns the column content according to the default
	// alignment.
	AlignDefault Alignment = iota
	// AlignLeft aligns the column content to the left.
	AlignLeft
	// AlignCenter aligns the column content to the center.
	AlignCenter
	// AlignRight aligns the column content to the right.
	AlignRight
)
