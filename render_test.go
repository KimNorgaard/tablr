package tablr_test

import (
	"bytes"
	"testing"

	"github.com/KimNorgaard/tablr"
)

func TestTable_Render(t *testing.T) {
	tests := []struct {
		name    string
		columns []string
		rows    [][]string
		options []tablr.TableOption
		want    string
	}{
		{
			name:    "Simple table",
			columns: []string{"Name", "Age", "City"},
			options: []tablr.TableOption{},
			rows: [][]string{
				{"John Doe", "30", "New York"},
				{"Jane Smith", "25", "Los Angeles"},
			},
			want: `| Name       | Age | City        |
|------------|-----|-------------|
| John Doe   | 30  | New York    |
| Jane Smith | 25  | Los Angeles |
`,
		},
		{
			name:    "Simple table with alignment",
			columns: []string{"Name", "Age", "City"},
			options: []tablr.TableOption{
				tablr.WithAlignments([]tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight}),
			},
			rows: [][]string{
				{"John Doe", "30", "New York"},
				{"Jane Smith", "25", "Los Angeles"},
			},
			want: `| Name       | Age |        City |
|:-----------|:---:|------------:|
| John Doe   | 30  |    New York |
| Jane Smith | 25  | Los Angeles |
`,
		},
		{
			name:    "Center aligned with odd length",
			columns: []string{"Name", "Age", "City"},
			options: []tablr.TableOption{
				tablr.WithAlignments([]tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight}),
			},
			rows: [][]string{
				{"John Doe", "100", "New York"},
				{"Jane Smith", "255", "Los Angeles"},
			},
			want: `| Name       | Age |        City |
|:-----------|:---:|------------:|
| John Doe   | 100 |    New York |
| Jane Smith | 255 | Los Angeles |
`,
		},
		{
			name:    "No rows",
			columns: []string{"Name", "Age", "City"},
			options: []tablr.TableOption{
				tablr.WithAlignments([]tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight}),
			},
			rows: [][]string{},
			want: `| Name | Age | City |
|:-----|:---:|-----:|
`,
		},
		{
			name:    "With pipes",
			columns: []string{"Name | Lastname", "Age", "City"},
			options: []tablr.TableOption{
				tablr.WithAlignments([]tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight}),
			},
			rows: [][]string{
				{"John | Doe", "30", "New | York"},
				{"Jane | Smith", "25", "Los | Angeles"},
			},
			want: `| Name \| Lastname | Age |           City |
|:----------------|:---:|---------------:|
| John \| Doe     | 30  |    New \| York |
| Jane \| Smith   | 25  | Los \| Angeles |
`,
		},
		{
			name:    "WithAlignment option",
			columns: []string{"Name", "Age", "City"},
			options: []tablr.TableOption{
				tablr.WithAlignments([]tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight}),
				tablr.WithAlignment(1, tablr.AlignRight),
			},
			rows: [][]string{
				{"John Doe", "30", "New York"},
				{"Jane Smith", "25", "Los Angeles"},
			},
			want: `| Name       | Age |        City |
|:-----------|----:|------------:|
| John Doe   |  30 |    New York |
| Jane Smith |  25 | Los Angeles |
`,
		},
		{
			name:    "WithAlignment option with wrong index",
			columns: []string{"Name", "Age", "City"},
			options: []tablr.TableOption{
				tablr.WithAlignments([]tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight}),
				tablr.WithAlignment(5, tablr.AlignRight),
			},
			rows: [][]string{
				{"John Doe", "30", "New York"},
				{"Jane Smith", "25", "Los Angeles"},
			},
			want: `| Name       | Age |        City |
|:-----------|:---:|------------:|
| John Doe   | 30  |    New York |
| Jane Smith | 25  | Los Angeles |
`,
		},
		{
			name:    "WithMinColumnWidths option",
			columns: []string{"Name", "Age", "City"},
			options: []tablr.TableOption{
				tablr.WithAlignments([]tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight}),
				tablr.WithMinColumWidths([]int{10, 5, 15}),
			},
			rows: [][]string{
				{"John Doe", "30", "New York"},
				{"Jane Smith", "25", "Los Angeles"},
			},
			want: `| Name       |  Age  |            City |
|:-----------|:-----:|----------------:|
| John Doe   |  30   |        New York |
| Jane Smith |  25   |     Los Angeles |
`,
		},
		{
			name:    "WithMinColumnWidth option",
			columns: []string{"Name", "Age", "City"},
			options: []tablr.TableOption{
				tablr.WithAlignments([]tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight}),
				tablr.WithMinColumnWidth(1, 10),
			},
			rows: [][]string{
				{"John Doe", "30", "New York"},
				{"Jane Smith", "25", "Los Angeles"},
			},
			want: `| Name       |    Age     |        City |
|:-----------|:----------:|------------:|
| John Doe   |     30     |    New York |
| Jane Smith |     25     | Los Angeles |
`,
		},
		{
			name:    "WithMinColumnWidth option with wrong index",
			columns: []string{"Name", "Age", "City"},
			options: []tablr.TableOption{
				tablr.WithAlignments([]tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight}),
				tablr.WithMinColumnWidth(5, 10),
			},
			rows: [][]string{
				{"John Doe", "30", "New York"},
				{"Jane Smith", "25", "Los Angeles"},
			},
			want: `| Name       | Age |        City |
|:-----------|:---:|------------:|
| John Doe   | 30  |    New York |
| Jane Smith | 25  | Los Angeles |
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			table := tablr.New(w, tt.columns, tt.options...)
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

func TestTable_String(t *testing.T) {
	tests := []struct {
		name       string
		columns    []string
		rows       [][]string
		alignments func(t *tablr.Table)
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
			want: `| Name | Age | City |
|:-----|:---:|-----:|
`,
		},
		{
			name:       "With pipes",
			columns:    []string{"Name | Lastname", "Age", "City"},
			alignments: tablr.WithAlignments([]tablr.Alignment{tablr.AlignLeft, tablr.AlignCenter, tablr.AlignRight}),
			rows: [][]string{
				{"John | Doe", "30", "New | York"},
				{"Jane | Smith", "25", "Los | Angeles"},
			},
			want: `| Name \| Lastname | Age |           City |
|:----------------|:---:|---------------:|
| John \| Doe     | 30  |    New \| York |
| Jane \| Smith   | 25  | Los \| Angeles |
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := tablr.New(nil, tt.columns, tt.alignments)
			err := table.AddRows(tt.rows)
			if err != nil {
				t.Errorf("Table.AddRows() error = %v", err)
			}
			got := table.String()
			if got != tt.want {
				t.Errorf("Table.String() got = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}
