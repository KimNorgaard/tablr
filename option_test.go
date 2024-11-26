package tablr_test

import (
	"bytes"
	"testing"

	"github.com/KimNorgaard/tablr"
)

func TestTable_WithAlignment(t *testing.T) {
	t.Run("WithAlignment success", func(t *testing.T) {
		columns := []string{"Name", "Age", "City"}
		alignments := tablr.WithAlignments([]tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight})
		table := tablr.New(&bytes.Buffer{}, columns, alignments, tablr.WithAlignment(1, tablr.AlignRight))

		got := table.GetAlignments()
		want := []tablr.Alignment{tablr.AlignLeft, tablr.AlignRight, tablr.AlignRight}
		if !equalSlices(got, want) {
			t.Errorf("WithAlignment() got = %v, want %v", got, want)
		}
	})

	t.Run("WithAlignment invalid alignment", func(t *testing.T) {
		columns := []string{"Name", "Age", "City"}
		alignments := tablr.WithAlignments([]tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight})
		table := tablr.New(&bytes.Buffer{}, columns, alignments, tablr.WithAlignment(1, tablr.Alignment(10)))

		got := table.GetAlignments()
		want := []tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight}
		if !equalSlices(got, want) {
			t.Errorf("WithAlignment() got = %v, want %v", got, want)
		}
	})

	t.Run("WithAlignments invalid alignment", func(t *testing.T) {
		columns := []string{"Name", "Age", "City"}
		alignments := tablr.WithAlignments([]tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.Alignment(10)})
		table := tablr.New(&bytes.Buffer{}, columns, alignments)

		got := table.GetAlignments()
		want := []tablr.Alignment{tablr.AlignDefault, tablr.AlignDefault, tablr.AlignDefault}
		if !equalSlices(got, want) {
			t.Errorf("WithAlignment() got = %v, want %v", got, want)
		}
	})
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

func TestTable_WithHeaderAlignment(t *testing.T) {
	t.Run("WithHeaderAlignment success", func(t *testing.T) {
		columns := []string{"Name", "Age", "City"}
		alignments := tablr.WithHeaderAlignments([]tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight})
		table := tablr.New(&bytes.Buffer{}, columns, alignments, tablr.WithHeaderAlignment(1, tablr.AlignRight))

		got := table.GetHeaderAlignments()
		want := []tablr.Alignment{tablr.AlignLeft, tablr.AlignRight, tablr.AlignRight}
		if !equalSlices(got, want) {
			t.Errorf("WithHeaderAlignment() got = %v, want %v", got, want)
		}
	})

	t.Run("WithHeaderAlignment with invalid alignment", func(t *testing.T) {
		columns := []string{"Name", "Age", "City"}
		alignments := tablr.WithHeaderAlignments([]tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight})
		table := tablr.New(&bytes.Buffer{}, columns, alignments, tablr.WithHeaderAlignment(1, tablr.Alignment(10)))

		got := table.GetHeaderAlignments()
		want := []tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight}
		if !equalSlices(got, want) {
			t.Errorf("WithHeaderAlignment() got = %v, want %v", got, want)
		}
	})

	t.Run("WithHeaderAlignments with invalid alignment", func(t *testing.T) {
		columns := []string{"Name", "Age", "City"}
		alignments := tablr.WithHeaderAlignments([]tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.Alignment(10)})
		table := tablr.New(&bytes.Buffer{}, columns, alignments)

		got := table.GetHeaderAlignments()
		want := []tablr.Alignment{tablr.AlignDefault, tablr.AlignDefault, tablr.AlignDefault}
		if !equalSlices(got, want) {
			t.Errorf("WithHeaderAlignment() got = %v, want %v", got, want)
		}
	})
}

func TestTable_WithColumnAlignment(t *testing.T) {
	t.Run("WithColumnAlignment success", func(t *testing.T) {
		columns := []string{"Name", "Age", "City"}
		table := tablr.New(&bytes.Buffer{}, columns)
		err := table.AddColumn("test", tablr.WithColumnAlignment(tablr.AlignRight))
		if err != nil {
			t.Errorf("AddColumn() error = %v", err)
		}

		want := tablr.AlignRight
		got, err := table.GetAlignment(3)
		if err != nil {
			t.Errorf("GetAlignment() error = %v", err)
		}
		if got != want {
			t.Errorf("WithColumnAlignment() got = %v, want %v", got, want)
		}
	})

	t.Run("WithColumnAlignment with invalid alignment", func(t *testing.T) {
		columns := []string{"Name", "Age", "City"}
		table := tablr.New(&bytes.Buffer{}, columns)
		opts := tablr.WithColumnAlignment(tablr.Alignment(10))
		err := table.AddColumn("test", opts)
		if err == nil {
			t.Error("AddColumn() expected an error for invalid alignment, but got nil")
		}
	})
}

func TestTable_WithColumnHeaderAlignment(t *testing.T) {
	t.Run("WithColumnHeaderAlignment success", func(t *testing.T) {
		columns := []string{"Name", "Age", "City"}
		table := tablr.New(&bytes.Buffer{}, columns)
		err := table.AddColumn("test", tablr.WithColumnHeaderAlignment(tablr.AlignLeft))
		if err != nil {
			t.Errorf("AddColumn() error = %v", err)
		}

		want := tablr.AlignLeft
		got, err := table.GetHeaderAlignment(3)
		if err != nil {
			t.Errorf("GetAlignment() error = %v", err)
		}
		if got != want {
			t.Errorf("WithColumnHeaderAlignment() got = %v, want %v", got, want)
		}
	})

	t.Run("WithColumnHeaderAlignment with invalid alignment", func(t *testing.T) {
		columns := []string{"Name", "Age", "City"}
		table := tablr.New(&bytes.Buffer{}, columns)
		opts := tablr.WithColumnHeaderAlignment(tablr.Alignment(10))
		err := table.AddColumn("test", opts)
		if err == nil {
			t.Error("AddColumn() expected an error for invalid alignment, but got nil")
		}
	})
}
