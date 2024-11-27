package tablr

// TableOption represents an option for configuring a Table.
type TableOption func(*Table)

// WithHeaderAlignments sets the alignment for each header.
func WithHeaderAlignments(alignments []Alignment) TableOption {
	return func(t *Table) {
		copy(t.headerAlignments, alignments)
		t.adjustAlignments()
		for i, a := range t.headerAlignments {
			if !a.IsValid() {
				t.headerAlignments[i] = AlignDefault
				continue
			}
			t.headerAlignments[i] = a
		}
	}
}

// WithHeaderAlignment sets the alignment for a header.
func WithHeaderAlignment(index int, alignment Alignment) TableOption {
	return func(t *Table) {
		if index < 0 || index >= len(t.columns) {
			return
		}
		if !alignment.IsValid() {
			alignment = AlignDefault
		}
		t.headerAlignments[index] = alignment
		t.adjustAlignments()
	}
}

// WithAlignments sets the alignments for each non-header column.
func WithAlignments(alignments []Alignment) TableOption {
	return func(t *Table) {
		copy(t.columnAlignments, alignments)
		t.adjustAlignments()
		// Update the alignments for headers to make sure the column alignments
		// take precedence over the default alignments.
		for i, a := range t.columnAlignments {
			if !a.IsValid() {
				t.columnAlignments[i] = AlignDefault
				t.headerAlignments[i] = AlignDefault
				continue
			}
			t.columnAlignments[i] = a

			if t.headerAlignments[i] == AlignDefault {
				t.headerAlignments[i] = a
			}
		}
	}
}

// WithAlignment sets the alignment for a non-header column.
func WithAlignment(index int, alignment Alignment) TableOption {
	return func(t *Table) {
		if index < 0 || index >= len(t.columns) {
			return
		}
		if !alignment.IsValid() {
			alignment = AlignDefault
		}
		t.columnAlignments[index] = alignment
		if t.headerAlignments[index] == AlignDefault {
			t.headerAlignments[index] = alignment
		}
		t.adjustAlignments()
	}
}

// WithMinColumnWidths sets the minimum widths for multiple columns.
func WithMinColumWidths(minColumnWidths []int) TableOption {
	return func(t *Table) {
		t.columnMinWidths = minColumnWidths
		t.adjustColumnWidths()
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

type column struct {
	headerAlignment Alignment
	alignment       Alignment
}

// ColumnOption represents an option for configuring a Column.
type ColumnOption func(*column)

func WithColumnAlignment(alignment Alignment) ColumnOption {
	return func(c *column) {
		c.alignment = alignment
	}
}

func WithColumnHeaderAlignment(alignment Alignment) ColumnOption {
	return func(c *column) {
		c.headerAlignment = alignment
	}
}
