# tablr

[![CI](https://github.com/KimNorgaard/tablr/actions/workflows/ci.yaml/badge.svg)](https://github.com/KimNorgaard/tablr/actions/workflows/ci.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/KimNorgaard/tablr)](https://pkg.go.dev/github.com/KimNorgaard/tablr)
[![Go Report Card](https://goreportcard.com/badge/github.com/KimNorgaard/tablr)](https://goreportcard.com/report/github.com/KimNorgaard/tablr)
[![License](https://img.shields.io/github/license/KimNorgaard/tablr)](LICENSE)

A Go package for creating and rendering GitHub flavored Markdown tables.

## Features

-   Straightforward API
-   Flexible row manipulation (add, update, delete)
-   Customizable column and header alignment (left, center, right)
-   Customizable column width
-   Column reordering
-   Thread-safe
-   Input validation
-   No dependencies

## Installation

Install using `go get`:

```bash
go get github.com/KimNorgaard/tablr
```

## Usage

Here's a basic example:

```go
package main

import (
	"fmt"
	"os"

	"github.com/KimNorgaard/tablr"
)

func main() {
	// Create a new table with two columns
	table := tablr.New(os.Stdout, []string{"Name", "Age"})

	// Add some rows
	table.AddRow([]string{"John Doe", "30"})
	table.AddRow([]string{"Jane Doe", "25"})
	table.AddRow([]string{"John Smith", "40"})

	// Render the table
	table.Render()
}
```

This will output:

```
| Name       | Age |
|------------|-----|
| John Doe   | 30  |
| Jane Doe   | 25  |
| John Smith | 40  |
```

### Alignment

Customize column alignment:

```go
table := tablr.New(
    os.Stdout, []string{"Name", "Age"},
    tablr.WithAlignments([]tablr.Alignment{tablr.AlignLeft, tablr.AlignRight}),
)
table.AddRow([]string{"John Doe", "30"})
table.Render()
```

This will right-align the "Age" column:

```
| Name     | Age |
|:---------|----:|
| John Doe |  30 |
```

### Column Reordering

Reorder columns in the table:

```go
table := tablr.New(os.Stdout, []string{"Name", "Age", "City"})
table.AddRow([]string{"John Doe", "30", "New York"})
table.AddRow([]string{"Jane Doe", "25", "Los Angeles"})
table.AddRow([]string{"John Smith", "40", "Chicago"})

// Reorder columns to "City", "Name", "Age"
err := table.ReorderColumns([]int{2, 0, 1})
if err != nil {
    fmt.Println("Error:", err)
}
table.Render()
```

This will produce a table with reordered columns:

```
| City        | Name       | Age |
|-------------|------------|-----|
| New York    | John Doe   | 30  |
| Los Angeles | Jane Doe   | 25  |
| Chicago     | John Smith | 40  |
```

### Error Handling

Many functions return errors. Always check for errors after calling methods like
`GetRow`, `GetColumn`, `SetRow`, `SetColumn`, `DeleteRow`,
`DeleteColumn`, `SetAlignment`, and `SetAlignments`.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

[MIT](LICENSE)
