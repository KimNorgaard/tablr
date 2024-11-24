package main

import (
	"os"

	"github.com/KimNorgaard/tablr"
)

func main() {
	// Create a new table with two columns
	table := tablr.New(os.Stdout, []string{"Name", "Age"})

	// Add some rows
	table.AddRow([]string{"John Doe", "30"})
	table.AddRow([]string{"Jane Doe", "25"})

	// Render the table
	table.Render()
}
