package tablr

import (
	"fmt"
	"io"
	"sync"
)

// Table represents a Markdown table.
type Table struct {
	mu               sync.RWMutex
	writer           io.Writer
	columns          []string
	rows             [][]string
	headerAlignments []Alignment
	alignments       []Alignment
	columnMinWidths  []int
}

// New creates a new Markdown table with the given columns and options.
// Each column string is the header of the column.
func New(writer io.Writer, columns []string, opts ...TableOption) *Table {
	t := &Table{
		writer:           writer,
		columns:          columns,
		headerAlignments: make([]Alignment, len(columns)),
		alignments:       make([]Alignment, len(columns)),
		columnMinWidths:  make([]int, len(columns)),
		rows:             make([][]string, 0),
	}

	// Initialize columns
	for i, col := range columns {
		t.columns[i] = escapePipes(col)
		t.columnMinWidths[i] = len(col)
		t.headerAlignments[i] = AlignDefault
		t.alignments[i] = AlignDefault
	}

	for _, opt := range opts {
		opt(t)
	}

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
func (t *Table) AddColumn(col string, opts ...ColumnOption) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	c := &column{
		alignment:       AlignDefault,
		headerAlignment: AlignDefault,
	}

	for _, opt := range opts {
		opt(c)
	}

	return t.addColumn(col, c)
}

// addColumn adds a column to the table without locking.
func (t *Table) addColumn(column string, c *column) error {
	column = escapePipes(column)

	if !c.alignment.IsValid() {
		return fmt.Errorf("invalid column alignment: %v", c.alignment)
	}
	if !c.headerAlignment.IsValid() {
		return fmt.Errorf("invalid header alignment: %v", c.headerAlignment)
	}

	t.columns = append(t.columns, column)
	t.columnMinWidths = append(t.columnMinWidths, len(column))
	t.headerAlignments = append(t.headerAlignments, c.headerAlignment)
	t.alignments = append(t.alignments, c.alignment)

	return nil
}

// AddColumns appends multiple columns to the table.
func (t *Table) AddColumns(columns []string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	col := &column{
		alignment:       AlignDefault,
		headerAlignment: AlignDefault,
	}

	for _, column := range columns {
		if err := t.addColumn(column, col); err != nil {
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

	copy(t.headerAlignments[index:], t.headerAlignments[index+1:])
	t.headerAlignments = t.headerAlignments[:len(t.headerAlignments)-1]

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
	newHeaderAlignments := make([]Alignment, len(t.headerAlignments))
	newAlignments := make([]Alignment, len(t.alignments))
	newColumnMinWidths := make([]int, len(t.columnMinWidths))
	for i, newIndex := range newOrder {
		newColumns[i] = t.columns[newIndex]
		newHeaderAlignments[i] = t.headerAlignments[newIndex]
		newAlignments[i] = t.alignments[newIndex]
		newColumnMinWidths[i] = t.columnMinWidths[newIndex]
	}
	t.columns = newColumns
	t.headerAlignments = newHeaderAlignments
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

// GetHeaderAlignments returns the header alignments of the columns.
func (t *Table) GetHeaderAlignments() []Alignment {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.updateAlignments()

	return t.headerAlignments
}

// GetHeaderAlignment returns the header alignment of the column at the given index.
func (t *Table) GetHeaderAlignment(index int) (Alignment, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if index < 0 || index >= len(t.headerAlignments) {
		return 0, fmt.Errorf("header alignment index out of range: %d, alignments: %d", index, len(t.headerAlignments))
	}

	t.updateAlignments()

	return t.headerAlignments[index], nil
}

// GetAlignments returns the alignments of the columns.
func (t *Table) GetAlignments() []Alignment {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.updateAlignments()

	return t.alignments
}

// GetAlignment returns the alignment of the column at the given index.
func (t *Table) GetAlignment(index int) (Alignment, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if index < 0 || index >= len(t.alignments) {
		return 0, fmt.Errorf("alignment index out of range: %d, alignments: %d", index, len(t.alignments))
	}

	t.updateAlignments()

	return t.alignments[index], nil
}

// SetHeaderAlignments sets the header alignments of the columns.
func (t *Table) SetHeaderAlignments(alignments []Alignment) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, a := range alignments {
		if !a.IsValid() {
			return fmt.Errorf("invalid alignment: %v", a)
		}
	}

	t.headerAlignments = alignments

	return nil
}

// SetHeaderAlignment sets the header alignment of the column at the given index.
func (t *Table) SetHeaderAlignment(index int, alignment Alignment) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if index < 0 || index >= len(t.headerAlignments) {
		return fmt.Errorf("header alignment index out of range: %d, alignments: %d", index, len(t.headerAlignments))
	}

	if !alignment.IsValid() {
		return fmt.Errorf("invalid alignment: %v", alignment)
	}

	t.headerAlignments[index] = alignment

	return nil
}

// SetAlignments sets the alignments of the columns.
func (t *Table) SetAlignments(alignments []Alignment) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, a := range alignments {
		if !a.IsValid() {
			return fmt.Errorf("invalid alignment: %v", a)
		}
	}

	t.alignments = alignments

	return nil
}

// SetAlignment sets the alignment of the column at the given index.
func (t *Table) SetAlignment(index int, alignment Alignment) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if index < 0 || index >= len(t.alignments) {
		return fmt.Errorf("alignment index out of range: %d, alignments: %d", index, len(t.alignments))
	}

	if !alignment.IsValid() {
		return fmt.Errorf("invalid alignment: %v", alignment)
	}

	t.alignments[index] = alignment

	return nil
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

// updateAlignments updates the alignments and headerAlignments slices to match the number of columns.
// If the alignments slices have more elements than the number of columns, remove the extra elements.
// If the alignments slices have fewer elements, append the missing elements with the default alignment.
func (t *Table) updateAlignments() {
	if len(t.columns) == len(t.alignments) && len(t.alignments) == len(t.headerAlignments) {
		return
	}

	if len(t.alignments) > len(t.columns) {
		t.alignments = t.alignments[:len(t.columns)]
	} else {
		for i := 0; i < len(t.columns)-len(t.alignments)+1; i++ {
			t.alignments = append(t.alignments, AlignDefault)
		}
	}

	if len(t.headerAlignments) > len(t.columns) {
		t.headerAlignments = t.headerAlignments[:len(t.columns)]
	} else {
		for i := 0; i < len(t.columns)-len(t.headerAlignments)+1; i++ {
			t.headerAlignments = append(t.headerAlignments, AlignDefault)
		}
	}
}
