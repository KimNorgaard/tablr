package tablr_test

import (
	"bytes"
	"fmt"
	"sync"
	"testing"

	"github.com/KimNorgaard/tablr"
)

func TestAlignmentConstants(t *testing.T) {
	if tablr.AlignDefault != 0 {
		t.Error("AlignDefault should be 0")
	}
	if tablr.AlignLeft != 1 {
		t.Error("AlignLeft should be 1")
	}
	if tablr.AlignCenter != 2 {
		t.Error("AlignCenter should be 2")
	}
	if tablr.AlignRight != 3 {
		t.Error("AlignRight should be 3")
	}
}

func TestTable_Render(t *testing.T) {
	tests := []struct {
		name       string
		columns    []string
		rows       [][]string
		alignments func(t *tablr.Table)
		pretty     func(t *tablr.Table)
		want       string
	}{
		{
			name:       "Simple table",
			columns:    []string{"Name", "Age", "City"},
			alignments: tablr.WithAlignments([]tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight}),
			rows: [][]string{
				{"John Doe", "30", "New York"},
				{"Jane Smith", "25", "Los Angeles"},
			},
			pretty: tablr.WithPretty(true),
			want: `| Name       | Age |        City |
|:-----------|:---:|------------:|
| John Doe   | 30  |    New York |
| Jane Smith | 25  | Los Angeles |
`,
		},
		{
			name:       "Center aligned with odd length",
			columns:    []string{"Name", "Age", "City"},
			alignments: tablr.WithAlignments([]tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight}),
			rows: [][]string{
				{"John Doe", "100", "New York"},
				{"Jane Smith", "255", "Los Angeles"},
			},
			pretty: tablr.WithPretty(true),
			want: `| Name       | Age |        City |
|:-----------|:---:|------------:|
| John Doe   | 100 |    New York |
| Jane Smith | 255 | Los Angeles |
`,
		},
		{
			name:       "No rows",
			columns:    []string{"Name", "Age", "City"},
			alignments: tablr.WithAlignments([]tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight}),
			rows:       [][]string{},
			pretty:     tablr.WithPretty(true),
			want: `| Name | Age | City |
|:-----|:---:|-----:|
`,
		},
		{
			name:       "Not pretty",
			columns:    []string{"Name", "Age", "City"},
			alignments: tablr.WithAlignments([]tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight}),
			rows: [][]string{
				{"John Doe", "30", "New York"},
				{"Jane Smith", "25", "Los Angeles"},
			},
			pretty: tablr.WithPretty(false),
			want: `| Name | Age | City |
|:-----|:---:|-----:|
| John Doe | 30  | New York |
| Jane Smith | 25  | Los Angeles |
`,
		},
		{
			name:       "Pretty with pipes",
			columns:    []string{"Name | Lastname", "Age", "City"},
			alignments: tablr.WithAlignments([]tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight}),
			rows: [][]string{
				{"John | Doe", "30", "New | York"},
				{"Jane | Smith", "25", "Los | Angeles"},
			},
			pretty: tablr.WithPretty(true),
			want: `| Name \| Lastname | Age |           City |
|:----------------|:---:|---------------:|
| John \| Doe     | 30  |    New \| York |
| Jane \| Smith   | 25  | Los \| Angeles |
`,
		},
		{
			name:       "Not pretty with pipes",
			columns:    []string{"Name | Lastname", "Age", "City"},
			alignments: tablr.WithAlignments([]tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight}),
			rows: [][]string{
				{"John | Doe", "30", "New | York"},
				{"Jane | Smith", "25", "Los | Angeles"},
			},
			pretty: tablr.WithPretty(false),
			want: `| Name \| Lastname | Age | City |
|:----------------|:---:|-----:|
| John \| Doe     | 30  | New \| York |
| Jane \| Smith   | 25  | Los \| Angeles |
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			table := tablr.New(w, tt.columns, tt.alignments, tt.pretty)
			err := table.AddRows(tt.rows)
			if err != nil {
				t.Errorf("Table.AddRows() error = %v", err)
			}
			table.Render()
			got := w.String()
			if got != tt.want {
				t.Errorf("Table.Render() got = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func TestTable_Methods(t *testing.T) {
	// Initialize a table for testing
	w := &bytes.Buffer{}
	table := tablr.New(w, []string{"Name", "Age", "City"}, tablr.WithAlignments([]tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight}), tablr.WithPretty(true))
	err := table.AddRows([][]string{
		{"John Doe", "30", "New York"},
		{"Jane Smith", "25", "Los Angeles"},
	})
	if err != nil {
		t.Errorf("Table.AddRows() error = %v", err)
	}
	err = table.AddRows([][]string{
		{"John Doe", "30", "New York", "USA"},
	})
	if err == nil {
		t.Error("Table.AddRows() expected an error for invalid row length, but got nil")
	}

	t.Run("GetRows", func(t *testing.T) {
		got := table.GetRows()
		want := [][]string{
			{"John Doe", "30", "New York"},
			{"Jane Smith", "25", "Los Angeles"},
		}
		if !equalRows(got, want) {
			t.Errorf("GetRows() got = %v, want %v", got, want)
		}
	})

	t.Run("SetRows", func(t *testing.T) {
		newRows := [][]string{
			{"Alice", "28", "Chicago"},
			{"Bob", "35", "Seattle"},
		}
		err := table.SetRows(newRows)
		if err != nil {
			t.Errorf("SetRows() error = %v", err)
		}
		got := table.GetRows()
		if !equalRows(got, newRows) {
			t.Errorf("SetRows() got = %v, want %v", got, newRows)
		}
	})

	t.Run("GetRow", func(t *testing.T) {
		got, err := table.GetRow(0)
		if err != nil {
			t.Errorf("GetRow() error = %v", err)
		}
		want := []string{"Alice", "28", "Chicago"}
		if !equalSlices(got, want) {
			t.Errorf("GetRow() got = %v, want %v", got, want)
		}

		_, err = table.GetRow(10) // Invalid index
		if err == nil {
			t.Error("GetRow() expected an error for invalid index, but got nil")
		}
	})

	t.Run("SetRow", func(t *testing.T) {
		newRow := []string{"Carol", "40", "Miami"}
		err := table.SetRow(1, newRow)
		if err != nil {
			t.Errorf("SetRow() error = %v", err)
		}
		got, _ := table.GetRow(1)
		if !equalSlices(got, newRow) {
			t.Errorf("SetRow() got = %v, want %v", got, newRow)
		}

		err = table.SetRow(10, newRow) // Invalid index
		if err == nil {
			t.Error("SetRow() expected an error for invalid index, but got nil")
		}

		err = table.SetRow(0, []string{"Invalid"}) // Invalid row length
		if err == nil {
			t.Error("SetRow() expected an error for invalid row length, but got nil")
		}
	})

	t.Run("DeleteRow", func(t *testing.T) {
		err := table.DeleteRow(0)
		if err != nil {
			t.Errorf("DeleteRow() error = %v", err)
		}
		got := table.GetRows()
		want := [][]string{
			{"Carol", "40", "Miami"},
		}
		if !equalRows(got, want) {
			t.Errorf("DeleteRow() got = %v, want %v", got, want)
		}

		err = table.DeleteRow(10) // Invalid index
		if err == nil {
			t.Error("DeleteRow() expected an error for invalid index, but got nil")
		}
	})

	t.Run("GetColumns", func(t *testing.T) {
		got := table.GetColumns()
		want := []string{"Name", "Age", "City"}
		if !equalSlices(got, want) {
			t.Errorf("GetColumns() got = %v, want %v", got, want)
		}
	})

	t.Run("GetColumn", func(t *testing.T) {
		got, err := table.GetColumn(0)
		if err != nil {
			t.Errorf("GetColumn() error = %v", err)
		}
		want := "Name"
		if got != want {
			t.Errorf("GetColumn() got = %v, want %v", got, want)
		}

		_, err = table.GetColumn(10) // Invalid index
		if err == nil {
			t.Error("GetColumn() expected an error for invalid index, but got nil")
		}
	})

	t.Run("SetColumns", func(t *testing.T) {
		newColumns := []string{"First Name", "Age", "Location"}
		err := table.SetColumns(newColumns)
		if err != nil {
			t.Errorf("SetColumns() error = %v", err)
		}
		if !equalSlices(table.GetColumns(), newColumns) {
			t.Errorf("SetColumns() got = %v, want %v", table.GetColumns(), newColumns)
		}
	})

	t.Run("SetColumn", func(t *testing.T) {
		err := table.SetColumn(1, "Years")
		if err != nil {
			t.Errorf("SetColumn() error = %v", err)
		}
		want := []string{"First Name", "Years", "Location"}
		if !equalSlices(table.GetColumns(), want) {
			t.Errorf("SetColumn() got = %v, want %v", table.GetColumns(), want)
		}

		err = table.SetColumn(10, "Invalid") // Invalid index
		if err == nil {
			t.Error("SetColumn() expected an error for invalid index, but got nil")
		}
	})

	t.Run("GetAlignment", func(t *testing.T) {
		got, err := table.GetAlignment(1)
		if err != nil {
			t.Errorf("GetAlignment() error = %v", err)
		}
		want := tablr.AlignCenter
		if got != want {
			t.Errorf("GetAlignment() got = %v, want %v", got, want)
		}
		_, err = table.GetAlignment(10) // Invalid index
		if err == nil {
			t.Error("GetAlignment() expected an error for invalid index, but got nil")
		}
	})

	t.Run("GetAlignments", func(t *testing.T) {
		got := table.GetAlignments()
		want := []tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight}
		if !equalSlices(got, want) {
			t.Errorf("GetAlignments() got = %v, want %v", got, want)
		}
	})

	t.Run("SetAlignment", func(t *testing.T) {
		err := table.SetAlignment(1, tablr.AlignRight)
		if err != nil {
			t.Errorf("SetAlignment() error = %v", err)
		}
		want := []tablr.Alignment{tablr.AlignLeft, tablr.AlignRight, tablr.AlignRight}
		got := table.GetAlignments()
		if !equalSlices(got, want) {
			t.Errorf("SetAlignment() got = %v, want %v", got, want)
		}
		err = table.SetAlignment(10, tablr.AlignRight) // Invalid index
		if err == nil {
			t.Error("SetAlignment() expected an error for invalid index, but got nil")
		}
	})

	t.Run("SetAlignments", func(t *testing.T) {
		err := table.SetAlignments([]tablr.Alignment{tablr.AlignCenter, tablr.AlignLeft, tablr.AlignCenter})
		if err != nil {
			t.Errorf("SetAlignments() error = %v", err)
		}
		want := []tablr.Alignment{tablr.AlignCenter, tablr.AlignLeft, tablr.AlignCenter}
		got := table.GetAlignments()
		if !equalSlices(got, want) {
			t.Errorf("SetAlignments() got = %v, want %v", got, want)
		}
	})

	t.Run("GetColumnWidth", func(t *testing.T) {
		got, err := table.GetColumnWidth(1)
		if err != nil {
			t.Errorf("GetColumnWidth() error = %v", err)
		}
		want := 3
		if got != want {
			t.Errorf("GetColumnWidth() got = %v, want %v", got, want)
		}
		_, err = table.GetColumnWidth(10) // Invalid index
		if err == nil {
			t.Error("GetColumnWidth() expected an error for invalid index, but got nil")
		}
	})

	t.Run("GetColumnWidths", func(t *testing.T) {
		got := table.GetColumnWidths()
		want := []int{5, 3, 5}
		t.Logf("rows: %v", table.GetRows())
		if !equalSlices(got, want) {
			t.Errorf("GetColumnWidths() got = %v, want %v", got, want)
		}
	})

	t.Run("DeleteColumn", func(t *testing.T) {
		err := table.DeleteColumn(1)
		if err != nil {
			t.Errorf("DeleteColumn() error = %v", err)
		}
		got := table.GetColumns()
		want := []string{"First Name", "Location"}
		if !equalSlices(got, want) {
			t.Errorf("DeleteColumn() got = %v, want %v", got, want)
		}
		err = table.DeleteColumn(10) // Invalid index
		if err == nil {
			t.Error("DeleteColumn() expected an error for invalid index, but got nil")
		}
	})

	t.Run("AddColumn", func(t *testing.T) {
		err := table.AddColumn("Country")
		if err != nil {
			t.Errorf("AddColumn() error = %v", err)
		}
		want := []string{"First Name", "Location", "Country"}
		got := table.GetColumns()
		if !equalSlices(got, want) {
			t.Errorf("AddColumn() got = %v, want %v", got, want)
		}
	})

	t.Run("AddColumns", func(t *testing.T) {
		err := table.AddColumns([]string{"Zip", "Phone"})
		if err != nil {
			t.Errorf("AddColumns() error = %v", err)
		}
		want := []string{"First Name", "Location", "Country", "Zip", "Phone"}
		got := table.GetColumns()
		if !equalSlices(got, want) {
			t.Errorf("AddColumns() got = %v, want %v", got, want)
		}
	})

	t.Run("Reset", func(t *testing.T) {
		table.Reset()
		got := table.GetRows()
		if len(got) != 0 {
			t.Errorf("Reset() got = %v, want empty slice", got)
		}
	})
}

func TestTable_ReorderColumns(t *testing.T) {
	w := &bytes.Buffer{}
	table := tablr.New(w, []string{"Name", "Age", "City"}, tablr.WithAlignments([]tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight}), tablr.WithPretty(true))
	err := table.AddRows([][]string{
		{"John Doe", "30", "New York"},
		{"Jane Smith", "25", "Los Angeles"},
	})
	if err != nil {
		t.Errorf("Table.AddRows() error = %v", err)
	}

	t.Run("ReorderColumns", func(t *testing.T) {
		err := table.ReorderColumns([]int{2, 0, 1})
		if err != nil {
			t.Errorf("ReorderColumns() error = %v", err)
		}

		gotColumns := table.GetColumns()
		wantColumns := []string{"City", "Name", "Age"}
		if !equalSlices(gotColumns, wantColumns) {
			t.Errorf("ReorderColumns() got columns = %v, want %v", gotColumns, wantColumns)
		}

		gotRows := table.GetRows()
		wantRows := [][]string{
			{"New York", "John Doe", "30"},
			{"Los Angeles", "Jane Smith", "25"},
		}
		if !equalRows(gotRows, wantRows) {
			t.Errorf("ReorderColumns() got rows = %v, want %v", gotRows, wantRows)
		}
	})

	t.Run("InvalidReorderColumns", func(t *testing.T) {
		err := table.ReorderColumns([]int{2, 0})
		if err == nil {
			t.Error("ReorderColumns() expected an error for invalid new order length, but got nil")
		}

		err = table.ReorderColumns([]int{2, 0, 3})
		if err == nil {
			t.Error("ReorderColumns() expected an error for invalid column index, but got nil")
		}

		err = table.ReorderColumns([]int{2, 0, 2})
		if err == nil {
			t.Error("ReorderColumns() expected an error for duplicate column index, but got nil")
		}
	})
}

func TestTable_Concurrency(t *testing.T) {
	table := tablr.New(&bytes.Buffer{}, []string{"Name", "Age"}, tablr.WithAlignments([]tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter}), tablr.WithPretty(true))

	t.Log("adding 100 rows")
	// Start multiple goroutines to append rows concurrently
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_ = table.AddRow([]string{fmt.Sprintf("Goroutine %d", i), fmt.Sprintf("%d", i)})
		}(i)
	}
	wg.Wait()

	t.Logf("number of rows: %d", len(table.GetRows()))
	// Check if the number of rows is correct
	if got, want := len(table.GetRows()), 100; got != want {
		t.Errorf("Expected %d rows, but got %d", want, got)
	}

	// Start multiple goroutines to access and modify the table concurrently
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			_ = table.SetRow(i, []string{fmt.Sprintf("Updated %d", i), fmt.Sprintf("%d", i*2)})
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
		want := []string{fmt.Sprintf("Updated %d", i), fmt.Sprintf("%d", i*2)}
		if !equalSlices(row, want) {
			t.Errorf("Row %d: got = %v, want %v", i, row, want)
		}
	}
}

func BenchmarkTable_Render(b *testing.B) {
	sizes := []int{10, 100, 1000, 100000}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("%d rows", size), func(b *testing.B) {
			columns := []string{"Name", "Age", "City"}
			alignments := tablr.WithAlignments([]tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight})
			pretty := tablr.WithPretty(true)

			// Prepare the table with the specified number of rows
			w := &bytes.Buffer{}
			table := tablr.New(w, columns, alignments, pretty)
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
