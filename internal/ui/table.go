package ui

import (
	"fmt"
	"strings"

)

// Heading prints a Markdown-style heading of the given level (1-6).
func Heading(title string, level int) {
	if level < 1 {
		level = 1
	}
	if level > 6 {
		level = 6
	}
	fmt.Printf("%s %s\n\n", strings.Repeat("#", level), title)
}

// PrintTable renders a simple ASCII table to stdout.
// headers → column names; rows → slice of string slices.
func PrintTable(headers []string, rows [][]string) {
	colWidths := calculateColumnWidths(headers, rows)

	printRow(headers, colWidths)
	printSeparator(colWidths)
	for _, row := range rows {
		printRow(row, colWidths)
	}
}

func calculateColumnWidths(headers []string, rows [][]string) []int {
	numCols := len(headers)
	colWidths := make([]int, numCols)

	// Determine max width for headers and data in each column
	for i, header := range headers {
		colWidths[i] = len(header) // Initialize with header length
		for _, row := range rows {
			if len(row[i]) > colWidths[i] {
				colWidths[i] = len(row[i])
			}
		}
	}

	return colWidths
}

func printRow(row []string, colWidths []int) {
	for i, item := range row {
		fmt.Print("|", padString(item, colWidths[i]))
	}
	fmt.Println("|")
}

func padString(s string, width int) string {
	return fmt.Sprintf(" %s%s ", s, strings.Repeat(" ", width-len(s)))
}

func printSeparator(colWidths []int) {
	for _, width := range colWidths {
		fmt.Print("+", strings.Repeat("-", width+2))
	}
	fmt.Println("+")
}
