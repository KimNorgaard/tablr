package tablr

import (
	"fmt"
	"io"
	"strings"
	"sync"
)

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

// Table represents a Markdown table.
type Table struct {
	mu              sync.RWMutex
	writer          io.Writer
	columns         []string
	rows            [][]string
	alignments      []Alignment
	columnMinWidths []int
}

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

// New creates a new Markdown table with the given columns and alignments.
func New(writer io.Writer, columns []string, opts ...TableOption) *Table {
	t := &Table{
		writer:          writer,
		columns:         columns,
		alignments:      make([]Alignment, len(columns)),
		columnMinWidths: make([]int, len(columns)),
		rows:            make([][]string, 0),
	}

	// Initialize column widths with column lengths
	for i, col := range columns {
		t.columns[i] = escapePipes(col)
		t.columnMinWidths[i] = len(col)
	}

	for _, opt := range opts {
		opt(t)
	}

	t.updateAlignments()

	return t
}

// AddRow appends a row to the table.
func (t *Table) AddRow(row []string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.addRow(row)
}

// addRow appends a row to the table without locking.
func (t *Table) addRow(row []string) error {
	if len(row) != len(t.columns) {
		return fmt.Errorf("incorrect number of values in row, should be %d", len(t.columns))
	}

	for i, val := range row {
		row[i] = escapePipes(val)
	}

	t.rows = append(t.rows, row)

	return nil
}

// AddRows appends multiple rows to the table.
func (t *Table) AddRows(rows [][]string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, row := range rows {
		if err := t.addRow(row); err != nil {
			return err
		}
	}

	return nil
}

// GetRow returns the row at the given index.
func (t *Table) GetRow(index int) ([]string, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if index < 0 || index >= len(t.rows) {
		return nil, fmt.Errorf("row index out of range: %d, rows: %d", index, len(t.rows))
	}

	return t.rows[index], nil
}

// GetRows returns the rows in the table.
func (t *Table) GetRows() [][]string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.rows
}

// SetRows sets the rows in the table.
func (t *Table) SetRows(rows [][]string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, row := range rows {
		for i, val := range row {
			row[i] = escapePipes(val)
		}
	}
	t.rows = rows

	return nil
}

// SetRow sets the row at the given index to the new row.
func (t *Table) SetRow(index int, row []string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if index < 0 || index >= len(t.rows) {
		return fmt.Errorf("row index out of range: %d, rows: %d", index, len(t.rows))
	}

	if len(row) != len(t.columns) {
		return fmt.Errorf("incorrect number of values in row, should be %d", len(t.columns))
	}

	for i, val := range row {
		row[i] = escapePipes(val)
	}

	t.rows[index] = row

	return nil
}

// DeleteRow deletes the row at the given index.
func (t *Table) DeleteRow(index int) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if index < 0 || index >= len(t.rows) {
		return fmt.Errorf("row index out of range: %d, rows: %d", index, len(t.rows))
	}

	copy(t.rows[index:], t.rows[index+1:])
	t.rows = t.rows[:len(t.rows)-1]

	return nil
}

// Reset clears all rows from the table.
func (t *Table) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.rows = make([][]string, 0)
}

// AddColumn adds a column to the table.
func (t *Table) AddColumn(column string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.addColumn(column)
}

// addColumn adds a column to the table without locking.
func (t *Table) addColumn(column string) error {
	column = escapePipes(column)
	t.columns = append(t.columns, column)
	t.columnMinWidths = append(t.columnMinWidths, len(column))
	t.alignments = append(t.alignments, AlignDefault)

	return nil
}

// AddColumns appends multiple columns to the table.
func (t *Table) AddColumns(columns []string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, column := range columns {
		if err := t.addColumn(column); err != nil {
			return err
		}
	}

	return nil
}

// DeleteColumn deletes the column at the given index.
func (t *Table) DeleteColumn(index int) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if index < 0 || index >= len(t.columns) {
		return fmt.Errorf("column index out of range: %d, columns: %d", index, len(t.columns))
	}

	copy(t.columns[index:], t.columns[index+1:])
	t.columns = t.columns[:len(t.columns)-1]

	copy(t.columnMinWidths[index:], t.columnMinWidths[index+1:])
	t.columnMinWidths = t.columnMinWidths[:len(t.columnMinWidths)-1]

	copy(t.alignments[index:], t.alignments[index+1:])
	t.alignments = t.alignments[:len(t.alignments)-1]

	return nil
}

// GetColumn returns the column at the given index.
func (t *Table) GetColumn(index int) (string, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if index < 0 || index >= len(t.columns) {
		return "", fmt.Errorf("column index out of range: %d, columns: %d", index, len(t.columns))
	}

	return t.columns[index], nil
}

// GetColumns returns the columns of the table.
func (t *Table) GetColumns() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.columns
}

// SetColumns sets the headers of the table.
func (t *Table) SetColumns(columns []string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.columns = columns
	t.updateAlignments()

	return nil
}

// SetColumn updates the column at the given index.
func (t *Table) SetColumn(index int, column string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if index < 0 || index >= len(t.columns) {
		return fmt.Errorf("column index out of range: %d, columns: %d", index, len(t.columns))
	}

	t.columns[index] = column

	return nil
}

// GetColumnWidth returns the width of the column at the given index.
func (t *Table) GetColumnWidth(index int) (int, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if index < 0 || index >= len(t.columnMinWidths) {
		return 0, fmt.Errorf("column index out of range: %d, columns: %d", index, len(t.columnMinWidths))
	}

	t.calculateColumnWidths()

	return t.columnMinWidths[index], nil
}

// GetColumnWidths returns the width of all columns.
func (t *Table) GetColumnWidths() []int {
	t.mu.RLock()
	defer t.mu.RUnlock()

	t.calculateColumnWidths()

	return t.columnMinWidths
}

// SetColumnWidths sets the widths of the given column.
// A negative value will set the width to 0.
func (t *Table) SetColumnMinWidth(index, width int) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if index < 0 || index >= len(t.columns) {
		return fmt.Errorf("column index out of range: %d, columns: %d", index, len(t.columns))
	}

	if width < 0 {
		width = 0
	}

	t.columnMinWidths[index] = width
	return nil
}

// SetColumnMinWidths sets the minimum widths of the columns.
// A negative value will set the width to 0.
func (t *Table) SetColumnMinWidths(widths []int) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(widths) != len(t.columns) {
		return fmt.Errorf("number of widths must match the number of columns: %d", len(t.columns))
	}

	for i, width := range widths {
		if width < 0 {
			width = 0
		}
		t.columnMinWidths[i] = width
	}

	return nil
}

// ReorderColumns reorders the columns according to the specified new order.
func (t *Table) ReorderColumns(newOrder []int) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(newOrder) != len(t.columns) {
		return fmt.Errorf("new order length must match the number of columns: %d", len(t.columns))
	}

	// Validate the new order
	seen := make(map[int]bool)
	for _, index := range newOrder {
		if index < 0 || index >= len(t.columns) {
			return fmt.Errorf("invalid column index in new order: %d, columns: %d", index, len(t.columns))
		}
		if seen[index] {
			return fmt.Errorf("duplicate column index in new order: %d", index)
		}
		seen[index] = true
	}

	// Reorder columns, alignments, and column widths
	newColumns := make([]string, len(t.columns))
	newAlignments := make([]Alignment, len(t.alignments))
	newColumnMinWidths := make([]int, len(t.columnMinWidths))
	for i, newIndex := range newOrder {
		newColumns[i] = t.columns[newIndex]
		newAlignments[i] = t.alignments[newIndex]
		newColumnMinWidths[i] = t.columnMinWidths[newIndex]
	}
	t.columns = newColumns
	t.alignments = newAlignments
	t.columnMinWidths = newColumnMinWidths

	// Reorder rows
	for i, row := range t.rows {
		newRow := make([]string, len(row))
		for j, newIndex := range newOrder {
			newRow[j] = row[newIndex]
		}
		t.rows[i] = newRow
	}

	return nil
}

// GetAlignment returns the alignment of the column at the given index.
func (t *Table) GetAlignment(index int) (Alignment, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if index < 0 || index >= len(t.alignments) {
		return 0, fmt.Errorf("alignment index out of range: %d, alignments: %d", index, len(t.alignments))
	}

	return t.alignments[index], nil
}

// GetAlignments returns the alignments of the columns.
func (t *Table) GetAlignments() []Alignment {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.alignments
}

// SetAlignment sets the alignment of the column at the given index.
func (t *Table) SetAlignment(index int, alignment Alignment) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if index < 0 || index >= len(t.alignments) {
		return fmt.Errorf("alignment index out of range: %d, alignments: %d", index, len(t.alignments))
	}

	t.alignments[index] = alignment

	return nil
}

// SetAlignments sets the alignments of the columns.
func (t *Table) SetAlignments(alignments []Alignment) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.alignments = alignments
	t.updateAlignments()

	return nil
}

// Render renders the table to the writer.
func (t *Table) Render() {
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

// calculateColumnWidths calculates the width of each column based on the data.
func (t *Table) calculateColumnWidths() {
	for _, row := range t.rows {
		for col, cell := range row {
			width := t.columnMinWidth(col)
			cellLen := len(cell)
			if cellLen > width {
				t.columnMinWidths[col] = cellLen
			}
		}
	}
}

// columnMinWidth returns the minimum width of the column at the given index.
func (t *Table) columnMinWidth(index int) int {
	if index < 0 || index >= len(t.columnMinWidths) {
		return 0
	}

	return t.columnMinWidths[index]
}

// updateAlignments updates the alignments slice to match the number of columns.
// If the alignments slice is longer than the number of columns, the extra
// elements are removed. If the alignments slice is shorter, the missing elements
// are appended with the default alignment.
func (t *Table) updateAlignments() {
	if len(t.columns) == len(t.alignments) {
		return
	}
	if len(t.alignments) > len(t.columns) {
		t.alignments = t.alignments[:len(t.columns)]
		return
	}
	for i := 0; i < len(t.columns)-len(t.alignments); i++ {
		t.alignments = append(t.alignments, AlignDefault)
	}
}

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
