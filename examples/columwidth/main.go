package main

import (
	"os"

	"github.com/KimNorgaard/tablr"
)

func main() {
	// Create a new table with two columns
	// First column has a minimum width of 20 characters
	table := tablr.New(os.Stdout, []string{"Name", "Age"},
		tablr.WithMinColumnWidth(0, 20)) //nolint:revive

	// Add some rows
	table.AddRow([]string{"John Doe", "30"})
	table.AddRow([]string{"Jane Doe", "25"})

	// Render the table
	table.Render()
}
