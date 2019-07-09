package main

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var (
	successColor = ui.Color(10)
	failureColor = ui.Color(124)
)

// NewStatusTable returns a table that renders test statuses.
func NewStatusTable(lines int) *StatusTable {
	table := widgets.NewTable()
	table.SetRect(0, 0, 5, lines)
	table.Border = false
	table.BorderRight = false
	table.FillRow = true
	table.RowSeparator = false

	texts := make([][]string, lines)
	statuses := make(map[int]ui.Style, lines)
	for i := 0; i < lines; i++ {
		texts[i] = []string{"   "}
		statuses[i] = ui.NewStyle(ui.ColorBlack)
	}

	return &StatusTable{
		statuses: statuses,
		lines:    lines,
		table:    table,
		texts:    texts,
	}
}

// StatusTable holds a table where every line represents a test status.
type StatusTable struct {
	statuses map[int]ui.Style
	texts    [][]string
	lines    int
	table    *widgets.Table
}

// Render renders the StatusTable on the screen.
func (s *StatusTable) Render() {
	s.table.Rows = s.texts
	s.table.RowStyles = s.statuses
	ui.Render(s.table)
}

// Push pushes a new status to the first line of the table.
func (s *StatusTable) Push(success bool) {
	color := ui.NewStyle(ui.ColorBlack, successColor)
	if !success {
		color = ui.NewStyle(ui.ColorBlack, failureColor)
	}

	newStatuses := map[int]ui.Style{0: color}
	for i := 0; i < s.lines; i++ {
		newStatuses[i+1] = s.statuses[i]
	}
	s.statuses = newStatuses
	s.Render()
}
