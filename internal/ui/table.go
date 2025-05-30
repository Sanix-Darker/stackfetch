package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
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
	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader(headers)
	t.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	t.SetCenterSeparator("")
	t.AppendBulk(rows)
	t.Render()
}
