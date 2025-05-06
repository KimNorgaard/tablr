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
	columnAlignments []Alignment
	columnMinWidths  []int
}

// New creates a new Markdown table with the given columns and options.
// Each column string is the header of the column.
func New(writer io.Writer, columns []string, opts ...TableOption) *Table {
	t := &Table{
		writer:           writer,
		columns:          columns,
		headerAlignments: make([]Alignment, len(columns)),
		columnAlignments: make([]Alignment, len(columns)),
		columnMinWidths:  make([]int, len(columns)),
		rows:             make([][]string, 0),
	}

	// Initialize columns
	for i, col := range columns {
		t.columns[i] = escapePipes(col)
		t.columnMinWidths[i] = len(col)
		t.headerAlignments[i] = AlignDefault
		t.columnAlignments[i] = AlignDefault
	}

	for _, opt := range opts {
		opt(t)
	}

	return t
}

// AddRow appends a row to the table.
// If the number of columns in the row is less than the number of columns in the
// table, the row will be padded with empty strings.
// If the number of columns in the row is greater than the number of columns in
// the table, the row will be truncated.
func (t *Table) AddRow(row []string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.addRowInternal(row)
	t.adjustColumnWidths()
}

// addRowInternal appends a row to the table without locking.
func (t *Table) addRowInternal(row []string) {
	row = t.adjustRowLength(row)

	for i, val := range row {
		row[i] = escapePipes(val)
	}

	t.rows = append(t.rows, row)
}

// AddRows appends multiple rows to the table.
// If the number of columns in a row is less than the number of columns in the
// table, the row will be padded with empty strings.
// If the number of columns in a row is greater than the number of columns in
// the table, the row will be truncated.
func (t *Table) AddRows(rows [][]string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, row := range rows {
		t.addRowInternal(row)
	}
	t.adjustColumnWidths()
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
// If the number of columns in a row is less than the number of columns in the
// table, the row will be padded with empty strings.
// If the number of columns in a row is greater than the number of columns in
// the table, the row will be truncated.
func (t *Table) SetRows(rows [][]string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Initialize newRows with the correct length
	newRows := make([][]string, len(rows))

	// Adjust the length of each row
	for rowIndex, row := range rows {
		adjustedRow := t.adjustRowLength(row)
		for colIndex, val := range adjustedRow {
			escapedVal := escapePipes(val)
			adjustedRow[colIndex] = escapedVal
			// Update columnMinWidths with the minimum width. This makes sure a
			// previously set larger columnMinWidths[colIndex] it not used
			t.columnMinWidths[colIndex] = max(t.columnMinWidths[colIndex], len(escapedVal), len(t.columns[colIndex]))
		}
		newRows[rowIndex] = adjustedRow
	}

	t.rows = newRows
	t.adjustColumnWidths()
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

	row = t.adjustRowLength(row)

	t.rows[index] = row
	t.adjustColumnWidths()

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

	t.adjustColumnWidths()

	return nil
}

// Reset clears all columns and rows from the table.
func (t *Table) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.columns = make([]string, 0)
	t.rows = make([][]string, 0)

	t.adjustColumnWidths()
}

// AddColumn adds a column to the table.
func (t *Table) AddColumn(header string, opts ...ColumnOption) {
	t.mu.Lock()
	defer t.mu.Unlock()

	c := &column{
		alignment:       AlignDefault,
		headerAlignment: AlignDefault,
	}

	for _, opt := range opts {
		opt(c)
	}

	t.addColumnInternal(header, c)
}

// addColumnInternal adds a column to the table without locking.
func (t *Table) addColumnInternal(header string, c *column) {
	header = escapePipes(header)

	if !c.alignment.IsValid() {
		c.alignment = AlignDefault
	}
	if !c.headerAlignment.IsValid() {
		c.headerAlignment = AlignDefault
	}

	t.columns = append(t.columns, header)
	t.columnMinWidths = append(t.columnMinWidths, len(header))
	t.headerAlignments = append(t.headerAlignments, c.headerAlignment)
	t.columnAlignments = append(t.columnAlignments, c.alignment)

	for i, row := range t.rows {
		t.rows[i] = append(row, "")
	}
}

// AddColumns appends multiple columns to the table.
func (t *Table) AddColumns(headers []string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	col := &column{
		alignment:       AlignDefault,
		headerAlignment: AlignDefault,
	}

	for _, header := range headers {
		t.addColumnInternal(header, col)
	}
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

	copy(t.columnAlignments[index:], t.columnAlignments[index+1:])
	t.columnAlignments = t.columnAlignments[:len(t.columnAlignments)-1]

	for i, row := range t.rows {
		copy(row[index:], row[index+1:])
		t.rows[i] = row[:len(row)-1]
	}
	t.adjustRowLenghts()

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
// If the new column count is less than the old one, the rows are truncated.
// If the new column count is greater than the old one, the rows are extended.
// Be careful when using this method, as it can lead to inconsistent data.
func (t *Table) SetColumns(columns []string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.columns = columns

	t.adjustRowLenghts()
	t.adjustAlignments()
	t.adjustColumnWidths()
}

// SetColumn updates the column at the given index.
func (t *Table) SetColumn(index int, column string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if index < 0 || index >= len(t.columns) {
		return fmt.Errorf("column index out of range: %d, columns: %d", index, len(t.columns))
	}

	t.columns[index] = column

	t.adjustColumnWidths()

	return nil
}

// GetColumnWidth returns the width of the column at the given index.
// Returns 0 if column index is out of range.
func (t *Table) GetColumnWidth(index int) int {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.columnMinWidth(index)
}

// GetColumnWidths returns the width of all columns.
func (t *Table) GetColumnWidths() []int {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.columnMinWidths
}

// SetColumnWidths sets the widths of the given column.
func (t *Table) SetColumnMinWidth(index, width int) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if index < 0 || index >= len(t.columns) {
		return fmt.Errorf("column index out of range: %d, columns: %d", index, len(t.columns))
	}

	if width < 0 {
		width = t.columnMinWidths[index]
	}

	t.columnMinWidths[index] = width
	t.adjustColumnWidths()

	return nil
}

// SetColumnMinWidths sets the minimum widths of the columns.
func (t *Table) SetColumnMinWidths(widths []int) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(widths) != len(t.columns) {
		return fmt.Errorf("number of widths must match the number of columns: %d", len(t.columns))
	}

	for i, width := range widths {
		if width < 0 {
			widths[i] = t.columnMinWidths[i]
		}
	}

	t.columnMinWidths = widths
	t.adjustColumnWidths()

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
	newAlignments := make([]Alignment, len(t.columnAlignments))
	newColumnMinWidths := make([]int, len(t.columnMinWidths))
	for i, newIndex := range newOrder {
		newColumns[i] = t.columns[newIndex]
		newHeaderAlignments[i] = t.headerAlignments[newIndex]
		newAlignments[i] = t.columnAlignments[newIndex]
		newColumnMinWidths[i] = t.columnMinWidths[newIndex]
	}
	t.columns = newColumns
	t.headerAlignments = newHeaderAlignments
	t.columnAlignments = newAlignments
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

	return t.headerAlignments
}

// GetHeaderAlignment returns the header alignment of the column at the given index.
func (t *Table) GetHeaderAlignment(index int) (Alignment, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if index < 0 || index >= len(t.headerAlignments) {
		return 0, fmt.Errorf("header alignment index out of range: %d, alignments: %d", index, len(t.headerAlignments))
	}

	return t.headerAlignments[index], nil
}

// GetAlignments returns the alignments of the columns.
func (t *Table) GetAlignments() []Alignment {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.columnAlignments
}

// GetAlignment returns the alignment of the column at the given index.
func (t *Table) GetAlignment(index int) (Alignment, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if index < 0 || index >= len(t.columnAlignments) {
		return 0, fmt.Errorf("alignment index out of range: %d, alignments: %d", index, len(t.columnAlignments))
	}

	return t.columnAlignments[index], nil
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
	t.adjustAlignments()

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
	t.adjustAlignments()

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

	t.columnAlignments = alignments
	t.adjustAlignments()

	return nil
}

// SetAlignment sets the alignment of the column at the given index.
func (t *Table) SetAlignment(index int, alignment Alignment) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if index < 0 || index >= len(t.columnAlignments) {
		return fmt.Errorf("alignment index out of range: %d, alignments: %d", index, len(t.columnAlignments))
	}

	if !alignment.IsValid() {
		return fmt.Errorf("invalid alignment: %v", alignment)
	}

	t.columnAlignments[index] = alignment
	t.adjustAlignments()

	return nil
}

// adjustColumnWidths calculates the width of each column based on the data.
func (t *Table) adjustColumnWidths() {
	colLen := len(t.columns)
	colWidthsLen := len(t.columnMinWidths)
	switch {
	case colWidthsLen > colLen:
		t.columnMinWidths = t.columnMinWidths[:len(t.columns)]
	case colWidthsLen < colLen:
		t.columnMinWidths = append(t.columnMinWidths, make([]int, colLen-colWidthsLen)...)
	}

	// Default to column header lengths
	for i, col := range t.columns {
		if len(col) > t.columnMinWidths[i] {
			t.columnMinWidths[i] = len(col)
		}
	}

	// Override using longest cell value in each column
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

// adjustAlignments updates the alignments and headerAlignments slices to match the number of columns.
// If the alignments slices have more elements than the number of columns, remove the extra elements.
// If the alignments slices have fewer elements, append the missing elements with the default alignment.
func (t *Table) adjustAlignments() {
	columnLen := len(t.columns)
	columnAlignmentsLen := len(t.columnAlignments)
	headerAlignmentsLen := len(t.headerAlignments)

	if columnLen == columnAlignmentsLen && columnAlignmentsLen == headerAlignmentsLen {
		return
	}

	truncateOrAppend := func(alignments *[]Alignment, length int) {
		if len(*alignments) > length {
			*alignments = (*alignments)[:length]
		} else {
			for i := 0; i < length-len(*alignments)+1; i++ {
				*alignments = append(*alignments, AlignDefault)
			}
		}
	}

	truncateOrAppend(&t.columnAlignments, columnLen)
	truncateOrAppend(&t.headerAlignments, columnLen)
}

// adjustRowLength adjusts the length of the row to match the number of columns.
func (t *Table) adjustRowLength(row []string) []string {
	colLen := len(t.columns)
	rowLen := len(row)

	switch {
	case rowLen > colLen:
		row = row[:colLen]
	case rowLen < colLen:
		row = append(row, make([]string, colLen-rowLen)...)
	}

	return row
}

// adjustRowLenghts adjusts the length of each row to match the number of columns.
func (t *Table) adjustRowLenghts() {
	for i, row := range t.rows {
		t.rows[i] = t.adjustRowLength(row)
	}
}
