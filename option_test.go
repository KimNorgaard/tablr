package tablr_test

import (
	"bytes"
	"testing"

	"github.com/KimNorgaard/tablr"
)

func TestTable_WithAlignment(t *testing.T) {
	tables := []struct {
		name      string
		alignment tablr.Alignment
		want      []tablr.Alignment
	}{
		{
			name:      "WithAlignment success",
			alignment: tablr.AlignRight,
			want:      []tablr.Alignment{tablr.AlignDefault, tablr.AlignRight, tablr.AlignDefault},
		},
		{
			name:      "WithAlignment invalid alignment sets alignment to default",
			alignment: tablr.Alignment(10),
			want:      []tablr.Alignment{tablr.AlignDefault, tablr.AlignDefault, tablr.AlignDefault},
		},
		{
			name:      "WithAlignment left alignment",
			alignment: tablr.AlignLeft,
			want:      []tablr.Alignment{tablr.AlignDefault, tablr.AlignLeft, tablr.AlignDefault},
		},
		{
			name:      "WithAlignment center alignment",
			alignment: tablr.AlignCenter,
			want:      []tablr.Alignment{tablr.AlignDefault, tablr.AlignCenter, tablr.AlignDefault},
		},
	}

	for _, tt := range tables {
		t.Run(tt.name, func(t *testing.T) {
			table := tablr.New(&bytes.Buffer{}, defaultColumns, tablr.WithAlignment(1, tt.alignment))

			got := table.GetAlignments()
			if !equalSlices(got, tt.want) {
				t.Errorf("WithAlignment() got = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestTable_WithAlignments(t *testing.T) {
	tables := []struct {
		name       string
		alignments []tablr.Alignment
		want       []tablr.Alignment
	}{
		{
			name:       "WithAlignments invalid alignment sets alignment to default",
			alignments: []tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.Alignment(10)},
			want:       []tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignDefault},
		},
		{
			name:       "WithAlignments all valid alignments",
			alignments: []tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight},
			want:       []tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight},
		},
		{
			name:       "WithAlignments all default alignments",
			alignments: []tablr.Alignment{tablr.AlignDefault, tablr.AlignDefault, tablr.AlignDefault},
			want:       []tablr.Alignment{tablr.AlignDefault, tablr.AlignDefault, tablr.AlignDefault},
		},
		{
			name:       "WithAlignments mixed valid and invalid alignments",
			alignments: []tablr.Alignment{tablr.AlignLeft, tablr.Alignment(10), tablr.AlignRight},
			want:       []tablr.Alignment{tablr.AlignLeft, tablr.AlignDefault, tablr.AlignRight},
		},
		{
			name:       "WithAlignments more alignments than columns",
			alignments: []tablr.Alignment{tablr.AlignLeft, tablr.AlignDefault, tablr.AlignRight, tablr.AlignCenter},
			want:       []tablr.Alignment{tablr.AlignLeft, tablr.AlignDefault, tablr.AlignRight},
		},
		{
			name:       "WithAlignments fewer alignments than columns",
			alignments: []tablr.Alignment{tablr.AlignLeft, tablr.AlignDefault},
			want:       []tablr.Alignment{tablr.AlignLeft, tablr.AlignDefault, tablr.AlignDefault},
		},
	}

	for _, tt := range tables {
		t.Run(tt.name, func(t *testing.T) {
			alignments := tablr.WithAlignments(tt.alignments)
			table := tablr.New(&bytes.Buffer{}, defaultColumns, alignments)

			got := table.GetAlignments()
			if !equalSlices(got, tt.want) {
				t.Errorf("WithAlignments() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTable_WithMinColumnWidths(t *testing.T) {
	tables := []struct {
		name      string
		minWidths []int
		want      []int
		columns   []string
		rows      [][]string
	}{
		{
			name:      "valid minimum column widths",
			minWidths: []int{10, 5, 15},
			want:      []int{10, 5, 15},
			columns:   defaultColumns,
		},
		{
			name:      "too small minimum column widths should extend with column header lengths",
			minWidths: []int{10, 5},
			want:      []int{10, 5, len(defaultColumns[2])},
			columns:   defaultColumns,
		},
		{
			name:      "too large minimum column widths should truncate using column lengths",
			minWidths: []int{10, 5, 10, 20},
			want:      []int{10, 5, 10},
			columns:   defaultColumns,
		},
		{
			name:      "with rows and minimum width larger than header width",
			minWidths: []int{0, 10, 0},
			want:      []int{len(defaultColumns[0]), 10, 8},
			columns:   defaultColumns,
			rows:      [][]string{{"John", "42", "New York"}},
		},
		{
			name:      "with no rows and header width smaller than minimum column width",
			minWidths: []int{0, 10, 0},
			want:      []int{len(defaultColumns[0]), 10, len(defaultColumns[2])},
			columns:   defaultColumns,
		},
		{
			name:      "with no rows and header width larger than minimum column width",
			minWidths: []int{3},
			want:      []int{8},
			columns:   []string{"Location"},
		},
	}

	for _, tt := range tables {
		t.Run(tt.name, func(t *testing.T) {
			table := tablr.New(&bytes.Buffer{}, tt.columns, tablr.WithMinColumWidths(tt.minWidths))
			if len(tt.rows) > 0 {
				for _, row := range tt.rows {
					table.AddRow(row)
				}
			}
			got := table.GetColumnWidths()
			if !equalSlices(got, tt.want) {
				t.Errorf("WithMinColumnWidths() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTable_WithHeaderAlignment(t *testing.T) {
	tables := []struct {
		name        string
		alignments  []tablr.Alignment
		headerIdx   int
		headerAlign tablr.Alignment
		want        []tablr.Alignment
	}{
		{
			name:        "WithHeaderAlignment success",
			alignments:  []tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight},
			headerIdx:   1,
			headerAlign: tablr.AlignRight,
			want:        []tablr.Alignment{tablr.AlignLeft, tablr.AlignRight, tablr.AlignRight},
		},
		{
			name:        "WithHeaderAlignment with invalid alignment",
			alignments:  []tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight},
			headerIdx:   1,
			headerAlign: tablr.Alignment(10),
			want:        []tablr.Alignment{tablr.AlignLeft, tablr.AlignDefault, tablr.AlignRight},
		},
		{
			name:        "WithHeaderAlignment all default alignments",
			alignments:  []tablr.Alignment{tablr.AlignDefault, tablr.AlignDefault, tablr.AlignDefault},
			headerIdx:   2,
			headerAlign: tablr.AlignLeft,
			want:        []tablr.Alignment{tablr.AlignDefault, tablr.AlignDefault, tablr.AlignLeft},
		},
		{
			name:        "WithHeaderAlignment mixed valid and invalid alignments",
			alignments:  []tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight},
			headerIdx:   0,
			headerAlign: tablr.Alignment(10),
			want:        []tablr.Alignment{tablr.AlignDefault, tablr.AlignCenter, tablr.AlignRight},
		},
	}

	for _, tt := range tables {
		t.Run(tt.name, func(t *testing.T) {
			alignments := tablr.WithHeaderAlignments(tt.alignments)
			table := tablr.New(&bytes.Buffer{}, defaultColumns, alignments, tablr.WithHeaderAlignment(tt.headerIdx, tt.headerAlign))

			got := table.GetHeaderAlignments()
			if !equalSlices(got, tt.want) {
				t.Errorf("WithHeaderAlignment() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTable_WithHeaderAlignments(t *testing.T) {
	tables := []struct {
		name       string
		alignments []tablr.Alignment
		want       []tablr.Alignment
	}{
		{
			name:       "WithHeaderAlignments all valid alignments",
			alignments: []tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight},
			want:       []tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight},
		},
		{
			name:       "WithHeaderAlignments with invalid alignment",
			alignments: []tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.Alignment(10)},
			want:       []tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignDefault},
		},
		{
			name:       "WithHeaderAlignments all default alignments",
			alignments: []tablr.Alignment{tablr.AlignDefault, tablr.AlignDefault, tablr.AlignDefault},
			want:       []tablr.Alignment{tablr.AlignDefault, tablr.AlignDefault, tablr.AlignDefault},
		},
		{
			name:       "WithHeaderAlignments mixed valid and invalid alignments",
			alignments: []tablr.Alignment{tablr.AlignLeft, tablr.Alignment(10), tablr.AlignRight},
			want:       []tablr.Alignment{tablr.AlignLeft, tablr.AlignDefault, tablr.AlignRight},
		},
		{
			name:       "WithHeaderAlignments more alignments than columns",
			alignments: []tablr.Alignment{tablr.AlignLeft, tablr.AlignDefault, tablr.AlignRight, tablr.AlignCenter},
			want:       []tablr.Alignment{tablr.AlignLeft, tablr.AlignDefault, tablr.AlignRight},
		},
		{
			name:       "WithHeaderAlignments fewer alignments than columns",
			alignments: []tablr.Alignment{tablr.AlignLeft, tablr.AlignDefault},
			want:       []tablr.Alignment{tablr.AlignLeft, tablr.AlignDefault, tablr.AlignDefault},
		},
	}

	for _, tt := range tables {
		t.Run(tt.name, func(t *testing.T) {
			alignments := tablr.WithHeaderAlignments(tt.alignments)
			table := tablr.New(&bytes.Buffer{}, defaultColumns, alignments)

			got := table.GetHeaderAlignments()
			if !equalSlices(got, tt.want) {
				t.Errorf("WithHeaderAlignments() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTable_WithColumnAlignment(t *testing.T) {
	tables := []struct {
		name      string
		alignment tablr.Alignment
		want      tablr.Alignment
	}{
		{
			name:      "WithColumnAlignment success",
			alignment: tablr.AlignRight,
			want:      tablr.AlignRight,
		},
		{
			name:      "WithColumnAlignment with invalid alignment",
			alignment: tablr.Alignment(10),
			want:      tablr.AlignDefault,
		},
		{
			name:      "WithColumnAlignment default alignment",
			alignment: tablr.AlignDefault,
			want:      tablr.AlignDefault,
		},
		{
			name:      "WithColumnAlignment center alignment",
			alignment: tablr.AlignCenter,
			want:      tablr.AlignCenter,
		},
	}

	for _, tt := range tables {
		t.Run(tt.name, func(t *testing.T) {
			table := tablr.New(&bytes.Buffer{}, defaultColumns)
			table.AddColumn("Name", tablr.WithColumnAlignment(tt.alignment))

			got, err := table.GetAlignment(3)
			if err != nil {
				t.Errorf("GetAlignment() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("WithColumnAlignment() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTable_WithColumnHeaderAlignment(t *testing.T) {
	tables := []struct {
		name      string
		alignment tablr.Alignment
		want      tablr.Alignment
	}{
		{
			name:      "WithColumnHeaderAlignment success",
			alignment: tablr.AlignLeft,
			want:      tablr.AlignLeft,
		},
		{
			name:      "WithColumnHeaderAlignment with invalid alignment",
			alignment: tablr.Alignment(10),
			want:      tablr.AlignDefault,
		},
		{
			name:      "WithColumnHeaderAlignment default alignment",
			alignment: tablr.AlignDefault,
			want:      tablr.AlignDefault,
		},
		{
			name:      "WithColumnHeaderAlignment center alignment",
			alignment: tablr.AlignCenter,
			want:      tablr.AlignCenter,
		},
	}

	for _, tt := range tables {
		t.Run(tt.name, func(t *testing.T) {
			table := tablr.New(&bytes.Buffer{}, defaultColumns)
			table.AddColumn("Name", tablr.WithColumnHeaderAlignment(tt.alignment))

			got, err := table.GetHeaderAlignment(3)
			if err != nil {
				t.Errorf("GetAlignment() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("WithColumnHeaderAlignment() got = %v, want %v", got, tt.want)
			}
		})
	}
}
