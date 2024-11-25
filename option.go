package tablr

// TableOption represents an option for configuring a Table.
type TableOption func(*Table)

// WithAlignments sets the alignments for each column.
func WithAlignments(alignments []Alignment) TableOption {
	return func(t *Table) {
		t.alignments = alignments
	}
}

// WithAlignment sets the alignment for a column.
func WithAlignment(index int, alignment Alignment) TableOption {
	return func(t *Table) {
		if index < 0 || index >= len(t.columns) {
			return
		}
		t.alignments[index] = alignment
	}
}

// WithMinColumnWidths sets the minimum widths for multiple columns.
func WithMinColumWidths(minColumnWidths []int) TableOption {
	return func(t *Table) {
		t.columnMinWidths = minColumnWidths
	}
}

// WithMinColumnWidth sets the minimum width for a column.
func WithMinColumnWidth(index, minWidth int) TableOption {
	return func(t *Table) {
		if index < 0 || index >= len(t.columns) {
			return
		}
		t.columnMinWidths[index] = minWidth
	}
}
