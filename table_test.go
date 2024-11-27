package tablr_test

import (
	"bytes"
	"fmt"
	"sync"
	"testing"

	"github.com/KimNorgaard/tablr"
)

var (
	defaultColumns    = []string{"Name", "Age", "City"}
	defaultAlignments = []tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight}
)

func newTable() *tablr.Table {
	var (
		columns    []string
		alignments []tablr.Alignment
	)
	columns = append(columns, defaultColumns...)
	alignments = append(alignments, defaultAlignments...)
	return tablr.New(&bytes.Buffer{}, columns, tablr.WithAlignments(alignments))
}

func TestTable_Methods(t *testing.T) {
	t.Parallel()

	t.Run("AddRow", func(t *testing.T) {
		tests := []struct {
			name string
			row  []string
			want [][]string
		}{
			{
				name: "AddRow",
				row:  []string{"John Doe", "30", "New York"},
				want: [][]string{
					{"John Doe", "30", "New York"},
				},
			},
			{
				name: "Row index larger than column length",
				row:  []string{"John Doe", "30", "New York", "USA"},
				want: [][]string{
					{"John Doe", "30", "New York"},
				},
			},
			{
				name: "Row index smaller than column length",
				row:  []string{"Jane Doe", "30"},
				want: [][]string{
					{"Jane Doe", "30", ""},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				table := newTable()
				table.AddRow(tt.row)
				got := table.GetRows()
				if !equalRows(tt.want, got) {
					t.Errorf("AddRow() got = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("AddRows", func(t *testing.T) {
		tests := []struct {
			name string
			rows [][]string
			want [][]string
		}{
			{
				name: "AddRows",
				rows: [][]string{
					{"John Doe", "30", "New York"},
					{"Jane Smith", "25", "Los Angeles"},
				},
				want: [][]string{
					{"John Doe", "30", "New York"},
					{"Jane Smith", "25", "Los Angeles"},
				},
			},
			{
				name: "Row index larger than column length",
				rows: [][]string{
					{"John Doe", "30", "New York", "USA"},
				},
				want: [][]string{
					{"John Doe", "30", "New York"},
				},
			},
			{
				name: "Row index smaller than column length",
				rows: [][]string{
					{"Jane Doe", "30"},
				},
				want: [][]string{
					{"Jane Doe", "30", ""},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				table := newTable()
				table.AddRows(tt.rows)
				got := table.GetRows()
				if !equalRows(tt.want, got) {
					t.Errorf("AddRows() got = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("GetRow", func(t *testing.T) {
		tests := []struct {
			name    string
			index   int
			want    []string
			wantErr bool
		}{
			{
				name:    "Valid index",
				index:   0,
				want:    []string{"Alice", "28", "Chicago"},
				wantErr: false,
			},
			{
				name:    "Invalid index",
				index:   10,
				want:    nil,
				wantErr: true,
			},
		}

		table := newTable()
		rows := [][]string{
			{"Alice", "28", "Chicago"},
			{"Bob", "35", "Seattle"},
		}
		table.AddRows(rows)

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := table.GetRow(tt.index)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetRow() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !equalSlices(got, tt.want) {
					t.Errorf("GetRow() got = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("GetRows", func(t *testing.T) {
		tests := []struct {
			name string
			cols []string
			rows [][]string
			want [][]string
		}{
			{
				name: "GetRows with valid data",
				rows: [][]string{
					{"John Doe", "30", "New York"},
					{"Jane Smith", "25", "Los Angeles"},
				},
				want: [][]string{
					{"John Doe", "30", "New York"},
					{"Jane Smith", "25", "Los Angeles"},
				},
			},
			{
				name: "GetRows with extra columns and fewer cells per row",
				cols: []string{"Age"},
				rows: [][]string{
					{"John Doe", "30", "New York"},
					{"Jane Smith", "25", "Los Angeles"},
				},
				want: [][]string{
					{"John Doe", "30", "New York", ""},
					{"Jane Smith", "25", "Los Angeles", ""},
				},
			},
			{
				name: "GetRows with extra columns and extra cells per row",
				cols: []string{"Age"},
				rows: [][]string{
					{"John Doe", "30", "New York", "42", "USA"},
					{"Jane Smith", "25", "Los Angeles", "42", "USA"},
				},
				want: [][]string{
					{"John Doe", "30", "New York", "42"},
					{"Jane Smith", "25", "Los Angeles", "42"},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				table := newTable()
				if len(tt.cols) > 0 {
					table.AddColumns(tt.cols)
				}
				table.AddRows(tt.rows)
				got := table.GetRows()
				if !equalRows(got, tt.want) {
					t.Errorf("GetRows() got = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("SetRow", func(t *testing.T) {
		tests := []struct {
			name    string
			index   int
			newRow  []string
			want    []string
			wantErr bool
		}{
			{
				name:    "Valid index",
				index:   1,
				newRow:  []string{"Carol", "40", "Miami"},
				want:    []string{"Carol", "40", "Miami"},
				wantErr: false,
			},
			{
				name:    "Invalid index",
				index:   10,
				newRow:  []string{"Carol", "40", "Miami"},
				want:    nil,
				wantErr: true,
			},
			{
				name:    "Invalid row length",
				index:   0,
				newRow:  []string{"Invalid"},
				want:    nil,
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				table := newTable()
				rows := [][]string{
					{"Alice", "28", "Chicago"},
					{"Bob", "35", "Seattle"},
				}
				table.AddRows(rows)
				err := table.SetRow(tt.index, tt.newRow)
				if (err != nil) != tt.wantErr {
					t.Errorf("SetRow() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !tt.wantErr {
					got, _ := table.GetRow(tt.index)
					if !equalSlices(got, tt.want) {
						t.Errorf("SetRow() got = %v, want %v", got, tt.want)
					}
				}
			})
		}
	})

	t.Run("SetRows", func(t *testing.T) {
		tests := []struct {
			name    string
			newRows [][]string
			want    [][]string
		}{
			{
				name: "SetRows with valid data",
				newRows: [][]string{
					{"Alice", "28", "Chicago"},
					{"Bob", "35", "Seattle"},
				},
				want: [][]string{
					{"Alice", "28", "Chicago"},
					{"Bob", "35", "Seattle"},
				},
			},
			{
				name: "SetRows with rows of different lengths",
				newRows: [][]string{
					{"Alice", "28", "Chicago", "Long row"},
					{"Bob", "35"},
				},
				want: [][]string{
					{"Alice", "28", "Chicago"},
					{"Bob", "35", ""},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				table := newTable()
				table.SetRows(tt.newRows)
				got := table.GetRows()
				if !equalRows(got, tt.want) {
					t.Errorf("SetRows() got = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("DeleteRow", func(t *testing.T) {
		tests := []struct {
			name    string
			index   int
			want    [][]string
			wantErr bool
		}{
			{
				name:    "Valid index",
				index:   0,
				want:    [][]string{{"Bob", "35", "Seattle"}},
				wantErr: false,
			},
			{
				name:    "Invalid index",
				index:   10,
				want:    [][]string{{"Alice", "28", "Chicago"}, {"Bob", "35", "Seattle"}},
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				table := newTable()
				rows := [][]string{
					{"Alice", "28", "Chicago"},
					{"Bob", "35", "Seattle"},
				}
				table.AddRows(rows)
				err := table.DeleteRow(tt.index)
				if (err != nil) != tt.wantErr {
					t.Errorf("DeleteRow() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				got := table.GetRows()
				if !equalRows(got, tt.want) {
					t.Errorf("DeleteRow() got = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("GetColumn", func(t *testing.T) {
		tests := []struct {
			name    string
			index   int
			want    string
			wantErr bool
		}{
			{
				name:    "Valid index",
				index:   0,
				want:    "Name",
				wantErr: false,
			},
			{
				name:    "Invalid index",
				index:   10,
				want:    "",
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				table := newTable()
				got, err := table.GetColumn(tt.index)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetColumn() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("GetColumn() got = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("GetColumns", func(t *testing.T) {
		table := newTable()
		got := table.GetColumns()
		want := []string{"Name", "Age", "City"}
		if !equalSlices(got, want) {
			t.Errorf("GetColumns() got = %v, want %v", got, want)
		}
	})

	t.Run("SetColumn", func(t *testing.T) {
		tests := []struct {
			name      string
			index     int
			newColumn string
			want      []string
			wantErr   bool
		}{
			{
				name:      "Valid index",
				index:     1,
				newColumn: "Years",
				want:      []string{"Name", "Years", "City"},
				wantErr:   false,
			},
			{
				name:      "Invalid index",
				index:     10,
				newColumn: "Invalid",
				want:      []string{"Name", "Age", "City"},
				wantErr:   true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				table := newTable()
				err := table.SetColumn(tt.index, tt.newColumn)
				if (err != nil) != tt.wantErr {
					t.Errorf("SetColumn() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !tt.wantErr {
					if !equalSlices(table.GetColumns(), tt.want) {
						t.Errorf("SetColumn() got = %v, want %v", table.GetColumns(), tt.want)
					}
				}
			})
		}
	})

	t.Run("SetColumns", func(t *testing.T) {
		tests := []struct {
			name        string
			newCols     []string
			initialRows [][]string
			wantRows    [][]string
		}{
			{
				name:    "SetColumns with same length",
				newCols: []string{"First Name", "Age", "Location"},
				initialRows: [][]string{
					{"John", "30", "New York"},
					{"Jane", "42", "San Francisco"},
				},
				wantRows: [][]string{
					{"John", "30", "New York"},
					{"Jane", "42", "San Francisco"},
				},
			},
			{
				name:    "SetColumns with additional column",
				newCols: []string{"First Name", "Age", "Location", "Country"},
				initialRows: [][]string{
					{"John", "30", "New York"},
					{"Jane", "42", "San Francisco"},
				},
				wantRows: [][]string{
					{"John", "30", "New York", ""},
					{"Jane", "42", "San Francisco", ""},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				table := newTable()
				table.SetColumns(tt.newCols)
				if !equalSlices(table.GetColumns(), tt.newCols) {
					t.Errorf("SetColumns() got = %v, want %v", table.GetColumns(), tt.newCols)
				}

				table.SetRows(tt.initialRows)
				gotRows := table.GetRows()

				if !equalRows(gotRows, tt.wantRows) {
					t.Errorf("SetColumns() got = %v, want %v", gotRows, tt.wantRows)
				}
			})
		}
	})

	t.Run("GetAlignment", func(t *testing.T) {
		tests := []struct {
			name    string
			index   int
			want    tablr.Alignment
			wantErr bool
		}{
			{
				name:    "Valid index",
				index:   1,
				want:    tablr.AlignCenter,
				wantErr: false,
			},
			{
				name:    "Invalid index",
				index:   10,
				want:    tablr.Alignment(0),
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				table := newTable()
				got, err := table.GetAlignment(tt.index)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetAlignment() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("GetAlignment() got = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("GetAlignments", func(t *testing.T) {
		table := newTable()
		got := table.GetAlignments()
		want := []tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight}
		if !equalSlices(got, want) {
			t.Errorf("GetAlignments() got = %v, want %v", got, want)
		}
	})

	t.Run("SetAlignment", func(t *testing.T) {
		tests := []struct {
			name      string
			index     int
			alignment tablr.Alignment
			want      []tablr.Alignment
			wantErr   bool
		}{
			{
				name:      "Valid index",
				index:     1,
				alignment: tablr.AlignRight,
				want:      []tablr.Alignment{tablr.AlignLeft, tablr.AlignRight, tablr.AlignRight},
				wantErr:   false,
			},
			{
				name:      "Invalid index",
				index:     10,
				alignment: tablr.AlignRight,
				want:      []tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight},
				wantErr:   true,
			},
			{
				name:      "Invalid alignment",
				index:     1,
				alignment: tablr.Alignment(10),
				want:      []tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight},
				wantErr:   true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				table := newTable()
				err := table.SetAlignment(tt.index, tt.alignment)
				if (err != nil) != tt.wantErr {
					t.Errorf("SetAlignment() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				got := table.GetAlignments()
				if !equalSlices(got, tt.want) {
					t.Errorf("SetAlignment() got = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("SetAlignments", func(t *testing.T) {
		tests := []struct {
			name       string
			alignments []tablr.Alignment
			want       []tablr.Alignment
			wantErr    bool
		}{
			{
				name:       "Valid alignments",
				alignments: []tablr.Alignment{tablr.AlignCenter, tablr.AlignLeft, tablr.AlignCenter},
				want:       []tablr.Alignment{tablr.AlignCenter, tablr.AlignLeft, tablr.AlignCenter},
				wantErr:    false,
			},
			{
				name:       "Shorter alignments",
				alignments: []tablr.Alignment{tablr.AlignCenter, tablr.AlignLeft},
				want:       []tablr.Alignment{tablr.AlignCenter, tablr.AlignLeft, tablr.AlignDefault},
				wantErr:    false,
			},
			{
				name:       "Longer alignments",
				alignments: []tablr.Alignment{tablr.AlignCenter, tablr.AlignLeft, tablr.AlignCenter, tablr.AlignDefault, tablr.AlignRight},
				want:       []tablr.Alignment{tablr.AlignCenter, tablr.AlignLeft, tablr.AlignCenter},
				wantErr:    false,
			},
			{
				name:       "Invalid alignment",
				alignments: []tablr.Alignment{tablr.AlignCenter, tablr.AlignLeft, tablr.Alignment(10)},
				want:       nil,
				wantErr:    true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				table := newTable()
				err := table.SetAlignments(tt.alignments)
				if (err != nil) != tt.wantErr {
					t.Errorf("SetAlignments() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !tt.wantErr {
					got := table.GetAlignments()
					if !equalSlices(got, tt.want) {
						t.Errorf("SetAlignments() got = %v, want %v", got, tt.want)
					}
				}
			})
		}
	})

	t.Run("GetColumnWidth", func(t *testing.T) {
		tests := []struct {
			name      string
			index     int
			minWidth  int
			wantWidth int
			wantErr   bool
		}{
			{
				name:      "Valid index",
				index:     1,
				minWidth:  10,
				wantWidth: 10,
				wantErr:   false,
			},
			{
				name:      "Invalid index",
				index:     10,
				minWidth:  10,
				wantWidth: 0,
				wantErr:   true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				table := newTable()
				err := table.SetColumnMinWidth(tt.index, tt.minWidth)
				if (err != nil) != tt.wantErr {
					t.Errorf("SetColumnMinWidth() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				got := table.GetColumnWidth(tt.index)
				if got != tt.wantWidth {
					t.Errorf("GetColumnWidth() got = %v, want %v", got, tt.wantWidth)
				}
			})
		}
	})

	t.Run("GetColumnWidths", func(t *testing.T) {
		tests := []struct {
			name        string
			initialRows [][]string
			want        []int
		}{
			{
				name:        "Initial column widths",
				initialRows: nil,
				want:        []int{len(defaultColumns[0]), len(defaultColumns[1]), len(defaultColumns[2])},
			},
			{
				name: "Updated column widths after adding rows",
				initialRows: [][]string{
					{"John Longname", "10000", "San Francisco"},
				},
				want: []int{13, 5, 13},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				table := newTable()
				if tt.initialRows != nil {
					table.SetRows(tt.initialRows)
				}
				got := table.GetColumnWidths()
				if !equalSlices(got, tt.want) {
					t.Errorf("GetColumnWidths() got = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("SetColumnMinWidth", func(t *testing.T) {
		tests := []struct {
			name      string
			index     int
			minWidth  int
			wantWidth int
			wantErr   bool
		}{
			{
				name:      "Valid index",
				index:     1,
				minWidth:  10,
				wantWidth: 10,
				wantErr:   false,
			},
			{
				name:      "Invalid index",
				index:     10,
				minWidth:  10,
				wantWidth: 0,
				wantErr:   true,
			},
			{
				name:      "Negative width",
				index:     1,
				minWidth:  -1,
				wantWidth: len(defaultColumns[1]), // the initial column width
				wantErr:   false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				table := newTable()
				err := table.SetColumnMinWidth(tt.index, tt.minWidth)
				if (err != nil) != tt.wantErr {
					t.Errorf("SetColumnMinWidth() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !tt.wantErr {
					got := table.GetColumnWidth(tt.index)
					if got != tt.wantWidth {
						t.Errorf("SetColumnMinWidth() got = %v, want %v", got, tt.wantWidth)
					}
				}
			})
		}
	})

	t.Run("SetColumnMinWidths", func(t *testing.T) {
		tests := []struct {
			name       string
			addColumn  string
			minWidths  []int
			wantWidths []int
			wantErr    bool
		}{
			{
				name:       "Valid widths",
				minWidths:  []int{10, 15, 20},
				wantWidths: []int{10, 15, 20},
				wantErr:    false,
			},
			{
				name:       "Invalid length",
				minWidths:  []int{10, 15},
				wantWidths: []int{len(defaultColumns[0]), len(defaultColumns[1]), len(defaultColumns[2])}, // the initial column widths
				wantErr:    true,
			},
			{
				name:       "Negative width",
				minWidths:  []int{10, 15, 20, -1},
				wantWidths: []int{len(defaultColumns[0]), len(defaultColumns[1]), len(defaultColumns[2])}, // the initial column widths
				wantErr:    true,
			},
			{
				name:       "Negative width after adding column",
				addColumn:  "Age",
				minWidths:  []int{10, 15, 20, -1},
				wantWidths: []int{10, 15, 20, 3},
				wantErr:    false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				table := newTable()
				if tt.addColumn != "" {
					table.AddColumn(tt.addColumn)
				}
				err := table.SetColumnMinWidths(tt.minWidths)
				if (err != nil) != tt.wantErr {
					t.Errorf("SetColumnMinWidths() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				got := table.GetColumnWidths()
				if !equalSlices(got, tt.wantWidths) {
					t.Errorf("SetColumnMinWidths() got = %v, want %v", got, tt.wantWidths)
				}
			})
		}
	})

	t.Run("DeleteColumn", func(t *testing.T) {
		tests := []struct {
			name     string
			index    int
			wantCols []string
			wantRows [][]string
			wantErr  bool
		}{
			{
				name:     "Valid index",
				index:    1,
				wantCols: []string{"Name", "City"},
				wantRows: [][]string{
					{"John", "Los Angeles"},
					{"Jane", "New York"},
				},
				wantErr: false,
			},
			{
				name:     "Invalid index",
				index:    10,
				wantCols: []string{"Name", "Age", "City"},
				wantRows: [][]string{
					{"John", "30", "Los Angeles"},
					{"Jane", "42", "New York"},
				},
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				table := newTable()
				oldRows := [][]string{
					{"John", "30", "Los Angeles"},
					{"Jane", "42", "New York"},
				}
				table.SetRows(oldRows)
				err := table.DeleteColumn(tt.index)
				if (err != nil) != tt.wantErr {
					t.Errorf("DeleteColumn() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				got := table.GetColumns()
				if !equalSlices(got, tt.wantCols) {
					t.Errorf("DeleteColumn() got = %v, want %v", got, tt.wantCols)
				}
				gotRows := table.GetRows()
				if !equalRows(gotRows, tt.wantRows) {
					t.Errorf("DeleteColumn() got = %v, want %v", gotRows, tt.wantRows)
				}
			})
		}
	})

	t.Run("AddColumn", func(t *testing.T) {
		tests := []struct {
			name     string
			initial  []string
			newCol   string
			wantCols []string
			newRows  [][]string
		}{
			{
				name:     "Add single column",
				initial:  []string{"Name"},
				newCol:   "Location",
				wantCols: []string{"Name", "Location"},
			},
			{
				name:     "Add column to empty table",
				initial:  []string{},
				newCol:   "Age",
				wantCols: []string{"Age"},
			},
			{
				name:     "Add column to multiple existing columns",
				initial:  []string{"Name", "Age"},
				newCol:   "City",
				wantCols: []string{"Name", "Age", "City"},
			},
			{
				name:     "Add column after adding rows",
				initial:  []string{"Name"},
				newCol:   "Location",
				wantCols: []string{"Name", "Location"},
				newRows: [][]string{
					{"John", "New York"},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				table := newTable()
				table.SetColumns(tt.initial)
				if len(tt.newRows) > 0 {
					table.AddRows(tt.newRows)
				}
				table.AddColumn(tt.newCol)
				got := table.GetColumns()
				if !equalSlices(got, tt.wantCols) {
					t.Errorf("AddColumn() got = %v, want %v", got, tt.wantCols)
				}
			})
		}
	})

	t.Run("AddColumn with alignments", func(t *testing.T) {
		tests := []struct {
			name            string
			initialColumns  []string
			newColumn       string
			columnAlignment tablr.Alignment
			headerAlignment tablr.Alignment
			wantColumnAlign tablr.Alignment
			wantHeaderAlign tablr.Alignment
			wantErr         bool
		}{
			{
				name:            "Add column with alignments",
				initialColumns:  []string{"Name"},
				newColumn:       "Age",
				columnAlignment: tablr.AlignCenter,
				headerAlignment: tablr.AlignRight,
				wantColumnAlign: tablr.AlignCenter,
				wantHeaderAlign: tablr.AlignRight,
				wantErr:         false,
			},
			{
				name:            "Add column with default alignments",
				initialColumns:  []string{"Name"},
				newColumn:       "City",
				columnAlignment: tablr.AlignDefault,
				headerAlignment: tablr.AlignDefault,
				wantColumnAlign: tablr.AlignDefault,
				wantHeaderAlign: tablr.AlignDefault,
				wantErr:         false,
			},
			{
				name:            "Add column with invalid alignment",
				initialColumns:  []string{"Name"},
				newColumn:       "Country",
				columnAlignment: tablr.Alignment(10),
				headerAlignment: tablr.AlignRight,
				wantColumnAlign: tablr.Alignment(0),
				wantHeaderAlign: tablr.Alignment(0),
				wantErr:         true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				table := newTable()
				table.SetColumns(tt.initialColumns)
				table.AddColumn(tt.newColumn, tablr.WithColumnAlignment(tt.columnAlignment), tablr.WithColumnHeaderAlignment(tt.headerAlignment))
				if !tt.wantErr {
					got, err := table.GetAlignment(len(table.GetColumns()) - 1)
					if err != nil {
						t.Errorf("GetAlignment() error = %v", err)
					}
					if tt.wantColumnAlign != got {
						t.Errorf("AddColumn() got = %v, want %v", got, tt.wantColumnAlign)
					}

					got, err = table.GetHeaderAlignment(len(table.GetColumns()) - 1)
					if err != nil {
						t.Errorf("GetHeaderAlignment() error = %v", err)
					}
					if tt.wantHeaderAlign != got {
						t.Errorf("AddColumn() got = %v, want %v", got, tt.wantHeaderAlign)
					}
				}
			})
		}
	})

	t.Run("AddColumns", func(t *testing.T) {
		tests := []struct {
			name       string
			initial    []string
			newColumns []string
			want       []string
		}{
			{
				name:       "Add multiple columns",
				initial:    []string{"Name"},
				newColumns: []string{"Zip", "Phone"},
				want:       []string{"Name", "Zip", "Phone"},
			},
			{
				name:       "Add single column",
				initial:    []string{"Name", "Age"},
				newColumns: []string{"City"},
				want:       []string{"Name", "Age", "City"},
			},
			{
				name:       "Add columns to empty table",
				initial:    []string{},
				newColumns: []string{"Country", "Continent"},
				want:       []string{"Country", "Continent"},
			},
			{
				name:       "Add no columns",
				initial:    []string{"Name", "Age"},
				newColumns: []string{},
				want:       []string{"Name", "Age"},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				table := newTable()
				table.SetColumns(tt.initial)
				table.AddColumns(tt.newColumns)
				got := table.GetColumns()
				if !equalSlices(got, tt.want) {
					t.Errorf("AddColumns() got = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("GetHeaderAlignment", func(t *testing.T) {
		tests := []struct {
			name    string
			index   int
			want    tablr.Alignment
			wantErr bool
		}{
			{
				name:    "Valid index",
				index:   0,
				want:    tablr.AlignLeft,
				wantErr: false,
			},
			{
				name:    "Negative index",
				index:   -1,
				want:    tablr.Alignment(0),
				wantErr: true,
			},
			{
				name:    "Invalid index",
				index:   10,
				want:    tablr.Alignment(0),
				wantErr: true,
			},
			{
				name:    "Valid index with center alignment",
				index:   1,
				want:    tablr.AlignCenter,
				wantErr: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				table := newTable()
				got, err := table.GetHeaderAlignment(tt.index)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetHeaderAlignment() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("GetHeaderAlignment() got = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("SetHeaderAlignments", func(t *testing.T) {
		tests := []struct {
			name       string
			alignments []tablr.Alignment
			wantErr    bool
		}{
			{
				name:       "Valid alignments",
				alignments: []tablr.Alignment{tablr.AlignCenter, tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight, tablr.AlignCenter, tablr.AlignLeft},
				wantErr:    false,
			},
			{
				name:       "Invalid alignment",
				alignments: []tablr.Alignment{tablr.AlignCenter, tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight, tablr.AlignCenter, tablr.Alignment(10)},
				wantErr:    true,
			},
			{
				name:       "Shorter alignments",
				alignments: []tablr.Alignment{tablr.AlignCenter},
				wantErr:    false,
			},
			{
				name:       "Longer alignments",
				alignments: []tablr.Alignment{tablr.AlignCenter, tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight, tablr.AlignCenter, tablr.AlignLeft, tablr.AlignRight},
				wantErr:    false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				table := newTable()
				err := table.SetHeaderAlignments(tt.alignments)
				if (err != nil) != tt.wantErr {
					t.Errorf("SetHeaderAlignments() error = %v, wantErr %v", err, tt.wantErr)
				}
			})
		}
	})

	t.Run("SetHeaderAlignment", func(t *testing.T) {
		tests := []struct {
			name      string
			index     int
			alignment tablr.Alignment
			wantErr   bool
		}{
			{
				name:      "Valid index",
				index:     1,
				alignment: tablr.AlignRight,
				wantErr:   false,
			},
			{
				name:      "Invalid index",
				index:     10,
				alignment: tablr.AlignRight,
				wantErr:   true,
			},
			{
				name:      "Invalid alignment",
				index:     1,
				alignment: tablr.Alignment(10),
				wantErr:   true,
			},
			{
				name:      "Valid index with default alignment",
				index:     0,
				alignment: tablr.AlignDefault,
				wantErr:   false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				table := newTable()
				err := table.SetHeaderAlignment(tt.index, tt.alignment)
				if (err != nil) != tt.wantErr {
					t.Errorf("SetHeaderAlignment() error = %v, wantErr %v", err, tt.wantErr)
				}
			})
		}
	})

	t.Run("Reset", func(t *testing.T) {
		table := newTable()
		table.AddRow([]string{"John", "25", "New York"})
		table.Reset()
		want := 0
		got := table.GetRows()
		if want != len(got) {
			t.Errorf("Reset() got = %v, want empty slice", got)
		}
		gotWidths := table.GetColumnWidths()
		if want != len(gotWidths) {
			t.Errorf("Reset() got = %v, want empty slice", gotWidths)
		}
	})
}

func TestTable_updateAlignments(t *testing.T) {
	t.Parallel()

	table := newTable()

	tests := []struct {
		name       string
		alignments []tablr.Alignment
		want       []tablr.Alignment
		wantErr    bool
	}{
		{
			name:       "Valid alignments with extra values",
			alignments: []tablr.Alignment{tablr.AlignCenter, tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight, tablr.AlignCenter, tablr.AlignLeft},
			want:       []tablr.Alignment{tablr.AlignCenter, tablr.AlignLeft, tablr.AlignCenter},
			wantErr:    false,
		},
		{
			name:       "Valid alignments with fewer values",
			alignments: []tablr.Alignment{tablr.AlignCenter},
			want:       []tablr.Alignment{tablr.AlignCenter, tablr.AlignDefault, tablr.AlignDefault},
			wantErr:    false,
		},
		{
			name:       "Valid alignments with exact values",
			alignments: []tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight},
			want:       []tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight},
			wantErr:    false,
		},
		{
			name:       "Invalid alignment value",
			alignments: []tablr.Alignment{tablr.AlignCenter, tablr.Alignment(10), tablr.AlignLeft},
			want:       nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := table.SetHeaderAlignments(tt.alignments)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetHeaderAlignments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				got := table.GetHeaderAlignments()
				if !equalSlices(got, tt.want) {
					t.Errorf("SetHeaderAlignments() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestTable_ReorderColumns(t *testing.T) {
	t.Parallel()

	table := newTable()
	table.AddRows([][]string{
		{"John Doe", "30", "New York"},
		{"Jane Smith", "25", "Los Angeles"},
	})

	t.Run("ReorderColumns", func(t *testing.T) {
		tests := []struct {
			name     string
			newOrder []int
			wantCols []string
			wantRows [][]string
			wantErr  bool
		}{
			{
				name:     "Valid reorder",
				newOrder: []int{2, 0, 1},
				wantCols: []string{"City", "Name", "Age"},
				wantRows: [][]string{
					{"New York", "John Doe", "30"},
					{"Los Angeles", "Jane Smith", "25"},
				},
				wantErr: false,
			},
			{
				name:     "Invalid reorder length",
				newOrder: []int{2, 0},
				wantCols: []string{"Name", "Age", "City"},
				wantRows: [][]string{
					{"John Doe", "30", "New York"},
					{"Jane Smith", "25", "Los Angeles"},
				},
				wantErr: true,
			},
			{
				name:     "Invalid column index",
				newOrder: []int{2, 0, 3},
				wantCols: []string{"Name", "Age", "City"},
				wantRows: [][]string{
					{"John Doe", "30", "New York"},
					{"Jane Smith", "25", "Los Angeles"},
				},
				wantErr: true,
			},
			{
				name:     "Duplicate column index",
				newOrder: []int{2, 0, 2},
				wantCols: []string{"Name", "Age", "City"},
				wantRows: [][]string{
					{"John Doe", "30", "New York"},
					{"Jane Smith", "25", "Los Angeles"},
				},
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := table.ReorderColumns(tt.newOrder)
				if (err != nil) != tt.wantErr {
					t.Errorf("ReorderColumns() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if !tt.wantErr {
					gotColumns := table.GetColumns()
					if !equalSlices(gotColumns, tt.wantCols) {
						t.Errorf("ReorderColumns() got columns = %v, want %v", gotColumns, tt.wantCols)
					}

					gotRows := table.GetRows()
					if !equalRows(gotRows, tt.wantRows) {
						t.Errorf("ReorderColumns() got rows = %v, want %v", gotRows, tt.wantRows)
					}
				}
			})
		}
	})

	t.Run("InvalidReorderColumns", func(t *testing.T) {
		tests := []struct {
			name     string
			newOrder []int
			wantErr  bool
		}{
			{
				name:     "Invalid new order length",
				newOrder: []int{2, 0},
				wantErr:  true,
			},
			{
				name:     "Invalid column index",
				newOrder: []int{2, 0, 3},
				wantErr:  true,
			},
			{
				name:     "Duplicate column index",
				newOrder: []int{2, 0, 2},
				wantErr:  true,
			},
			{
				name:     "Empty new order",
				newOrder: []int{},
				wantErr:  true,
			},
			{
				name:     "Negative column index",
				newOrder: []int{2, 0, -1},
				wantErr:  true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := table.ReorderColumns(tt.newOrder)
				if (err != nil) != tt.wantErr {
					t.Errorf("ReorderColumns() error = %v, wantErr %v", err, tt.wantErr)
				}
			})
		}
	})
}

func TestTable_Concurrency(t *testing.T) {
	t.Parallel()

	table := newTable()

	// Start multiple goroutines to append rows concurrently
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			table.AddRow([]string{fmt.Sprintf("Goroutine %d", i), fmt.Sprintf("%d", i), "Test"})
		}(i)
	}
	wg.Wait()

	// Check if the number of rows is correct
	if got, want := len(table.GetRows()), 100; got != want {
		t.Errorf("Expected %d rows, but got %d", want, got)
	}

	// Start multiple goroutines to access and modify the table concurrently
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			_ = table.SetRow(i, []string{fmt.Sprintf("Updated %d", i), fmt.Sprintf("%d", i*2), "Test"})
		}
	}()

	go func() {
		defer wg.Done()
		for i := 50; i < 100; i++ {
			_, _ = table.GetRow(i)
		}
	}()
	wg.Wait()

	// Perform assertions on the final state of the table
	for i := 0; i < 50; i++ {
		row, err := table.GetRow(i)
		if err != nil {
			t.Errorf("GetRow(%d) error = %v", i, err)
		}
		want := []string{fmt.Sprintf("Updated %d", i), fmt.Sprintf("%d", i*2), "Test"}
		if !equalSlices(row, want) {
			t.Errorf("Row %d: got = %v, want %v", i, row, want)
		}
	}
}

func BenchmarkTable_Render(b *testing.B) {
	sizes := []int{10, 100, 1000, 100000}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("%d rows", size), func(b *testing.B) {
			// Prepare the table with the specified number of rows
			w := &bytes.Buffer{}
			table := tablr.New(w, defaultColumns, tablr.WithAlignments(defaultAlignments))
			for i := 0; i < size; i++ {
				table.AddRow([]string{fmt.Sprintf("Name %d", i), fmt.Sprintf("%d", i), fmt.Sprintf("City %d", i)})
			}

			// Reset the timer to exclude setup time
			b.ResetTimer()

			// Run the benchmark
			for i := 0; i < b.N; i++ {
				w.Reset()
				table.Render()
			}
		})
	}
}

// Helper functions for comparing slices and rows
func equalSlices[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func equalRows(a, b [][]string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, row := range a {
		if !equalSlices(row, b[i]) {
			return false
		}
	}
	return true
}
