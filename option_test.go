package tablr_test

import (
	"bytes"
	"testing"

	"github.com/KimNorgaard/tablr"
)

func TestTable_WithAlignment(t *testing.T) {
	columns := []string{"Name", "Age", "City"}
	alignments := tablr.WithAlignments([]tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight})
	table := tablr.New(&bytes.Buffer{}, columns, alignments, tablr.WithAlignment(1, tablr.AlignRight))

	got := table.GetAlignments()
	want := []tablr.Alignment{tablr.AlignLeft, tablr.AlignRight, tablr.AlignRight}
	if !equalSlices(got, want) {
		t.Errorf("WithAlignment() got = %v, want %v", got, want)
	}
}

func TestTable_WithMinColumnWidths(t *testing.T) {
	columns := []string{"Name", "Age", "City"}
	minWidths := []int{10, 5, 15}
	table := tablr.New(&bytes.Buffer{}, columns, tablr.WithMinColumWidths(minWidths))

	got := table.GetColumnWidths()
	if !equalSlices(got, minWidths) {
		t.Errorf("WithMinColumnWidths() got = %v, want %v", got, minWidths)
	}
}

func TestTable_WithMinColumnWidth(t *testing.T) {
	columns := []string{"Name", "Age", "City"}
	table := tablr.New(&bytes.Buffer{}, columns, tablr.WithMinColumnWidth(1, 10))

	got, err := table.GetColumnWidth(1)
	if err != nil {
		t.Errorf("GetColumnWidth() error = %v", err)
	}
	want := 10
	if got != want {
		t.Errorf("WithMinColumnWidth() got = %v, want %v", got, want)
	}
}
