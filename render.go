package tablr

import (
	"fmt"
	"strings"
)

// Render renders the table to the writer.
func (t *Table) Render() {
	t.updateAlignments()
	t.calculateColumnWidths()

	// Write header column
	for i, col := range t.columns {
		fmt.Fprint(t.writer, "| ")
		fmt.Fprint(t.writer, pad(col, t.columnMinWidths[i], t.alignments[i]))
		fmt.Fprint(t.writer, " ")
	}
	fmt.Fprintln(t.writer, "|")

	// Write alignment row
	for i, align := range t.alignments {
		fmt.Fprint(t.writer, "|")
		switch align {
		case AlignDefault:
			fmt.Fprint(t.writer, "-", strings.Repeat("-", t.columnMinWidths[i]), "-")
		case AlignLeft:
			fmt.Fprint(t.writer, ":", strings.Repeat("-", t.columnMinWidths[i]), "-")
		case AlignCenter:
			fmt.Fprint(t.writer, ":", strings.Repeat("-", t.columnMinWidths[i]), ":")
		case AlignRight:
			fmt.Fprint(t.writer, "-", strings.Repeat("-", t.columnMinWidths[i]), ":")
		}
	}
	fmt.Fprintln(t.writer, "|")

	// Write data rows
	for _, row := range t.rows {
		for i, cell := range row {
			fmt.Fprint(t.writer, "| ")
			fmt.Fprint(t.writer, pad(cell, t.columnMinWidths[i], t.alignments[i]))
			fmt.Fprint(t.writer, " ")
		}
		fmt.Fprintln(t.writer, "|")
	}
}

// String returns the table as a string.
func (t *Table) String() string {
	var sb strings.Builder
	t.writer = &sb
	t.Render()
	return sb.String()
}
