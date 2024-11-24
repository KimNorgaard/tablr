package main

import (
	"os"

	"github.com/KimNorgaard/tablr"
)

func main() {
	// Create a new table with two columns
	table := tablr.New(os.Stdout, []string{"Name", "Age", "City"}, tablr.WithPretty(true))

	// Add some rows
	table.AddRow([]string{"John Doe", "30", "New York"})
	table.AddRow([]string{"Jane Doe", "25", "Los Angeles"})
	table.AddRow([]string{"John Smith", "40", "Chicago"})

	// Render the table
	table.Render()
}
