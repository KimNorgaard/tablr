// Package tablr provides a simple way to create GitHub flavored Markdown
// tables.
//
// Example usage:
//
//	import (
//		"fmt"
//		"io"
//		"os"
//
//		"github.com/KimNorgaard/tablr"
//	)
//
//	func main() {
//		// Create a new table with two columns
//		table := tablr.New(os.Stdout, []string{"Name", "Age"})
//
//		// Add some rows
//		table.AddRow([]string{"John Doe", "30"})
//		table.AddRow([]string{"Jane Doe", "25"})
//
//		// Render the table
//		table.Render()
//	}
//
// This will output:
//
// | Name     | Age |
// |:---------|:----|
// | John Doe | 30  |
// | Jane Doe | 25  |
//
// You can also specify the alignment of each column:
//
//	table := tablr.New(os.Stdout, []string{"Name", "Age"}, tablr.WithAlignments([]tablr.Alignment{tablr.AlignLeft, tablr.AlignRight}))
//
// This will output:
//
// | Name       | Age | City        |
// |:-----------|:---:|:------------|
// | John Doe   | 30  | New York    |
// | Jane Doe   | 25  | Los Angeles |
// | John Smith | 40  | Chicago     |
//
// Column Reordering:
//
//	table := tablr.New(os.Stdout, []string{"Name", "Age", "City"})
//	table.AddRow([]string{"John Doe", "30", "New York"})
//	table.AddRow([]string{"Jane Doe", "25", "Los Angeles"})
//	table.AddRow([]string{"John Smith", "40", "Chicago"})
//
//	// Reorder columns to "City", "Name", "Age"
//	err := table.ReorderColumns([]int{2, 0, 1})
//	if err != nil {
//		fmt.Println("Error:", err)
//	}
//	table.Render()
//
// This will output a table with reordered columns:
//
// | City        | Name       | Age |
// |:------------|:-----------|:---:|
// | New York    | John Doe   | 30  |
// | Los Angeles | Jane Doe   | 25  |
//
// Column Reordering:
//
//	table := tablr.New(os.Stdout, []string{"Name", "Age", "City"})
//	table.AddRow([]string{"John Doe", "30", "New York"})
//	table.AddRow([]string{"Jane Doe", "25", "Los Angeles"})
//	table.AddRow([]string{"John Smith", "40", "Chicago"})
//
//	// Reorder columns to "City", "Name", "Age"
//	err := table.ReorderColumns([]int{2, 0, 1})
//	if err != nil {
//		fmt.Println("Error:", err)
//	}
//	table.Render()
//
// This will output a table with reordered columns:
//
// | City        | Name       | Age |
// |:------------|:-----------|:---:|
// | New York    | John Doe   | 30  |
// | Los Angeles | Jane Doe   | 25  |
//
// Column Reordering:
//
//	table := tablr.New(os.Stdout, []string{"Name", "Age", "City"})
//	table.AddRow([]string{"John Doe", "30", "New York"})
//	table.AddRow([]string{"Jane Doe", "25", "Los Angeles"})
//	table.AddRow([]string{"John Smith", "40", "Chicago"})
//
//	// Reorder columns to "City", "Name", "Age"
//	err := table.ReorderColumns([]int{2, 0, 1})
//	if err != nil {
//		fmt.Println("Error:", err)
//	}
//	table.Render()
//
// This will output a table with reordered columns:
//
// | City        | Name       | Age |
// |:------------|:-----------|:---:|
// | New York    | John Doe   | 30  |
// | Los Angeles | Jane Doe   | 25  |
//
// Column Reordering:
//
//	table := tablr.New(os.Stdout, []string{"Name", "Age", "City"})
//	table.AddRow([]string{"John Doe", "30", "New York"})
//	table.AddRow([]string{"Jane Doe", "25", "Los Angeles"})
//	table.AddRow([]string{"John Smith", "40", "Chicago"})
//
//	// Reorder columns to "City", "Name", "Age"
//	err := table.ReorderColumns([]int{2, 0, 1})
//	if err != nil {
//		fmt.Println("Error:", err)
//	}
//	table.Render()
//
// This will output a table with reordered columns:
//
// | City        | Name       | Age |
// |:------------|:-----------|:---:|
// | New York    | John Doe   | 30  |
// | Los Angeles | Jane Doe   | 25  |
// | Chicago     | John Smith | 40  |
//
// Error Handling:
//
// Many methods in this package return an error.  Always check for errors.  For example, when accessing a row:
//
//		row, err := table.GetRow(0)
//		if err != nil {
//			fmt.Println("Error:", err)
//		}
//
//	 Similar error handling should be used with `GetColumn`, `SetRow`, `SetColumn`, `DeleteRow`, `DeleteColumn`, `SetAlignment`, and `SetAlignments`.
package tablr
