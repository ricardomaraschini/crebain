package tui

import (
	"sync"

	"github.com/ricardomaraschini/crebain/trunner"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var (
	// kinda green
	successStyle = ui.Style{
		Fg: ui.ColorClear,
		Bg: ui.Color(10),
	}

	// kinda red
	failureStyle = ui.Style{
		Fg: ui.ColorClear,
		Bg: ui.Color(124),
	}
)

// NewStatusTable returns a table that renders test statuses.
func NewStatusTable(height int) *StatusTable {
	table := widgets.NewTable()
	table.SetRect(0, 0, 5, height)
	table.Border = false
	table.BorderRight = false
	table.FillRow = true
	table.RowSeparator = false

	maxRows := height - 2
	table.Rows = make([][]string, maxRows)
	table.RowStyles = make(map[int]ui.Style, maxRows)
	for i := 0; i < maxRows; i++ {
		table.Rows[i] = []string{"   "}
		table.RowStyles[i] = ui.NewStyle(ui.ColorClear)
	}

	return &StatusTable{
		maxRows: maxRows,
		table:   table,
	}
}

// StatusTable holds a table where every line represents a test status.
type StatusTable struct {
	sync.Mutex
	maxRows int
	table   *widgets.Table
}

// Event receives user events, this table has no action on events.
func (s *StatusTable) Event(event string) {}

// Push pushes a new status row to the first line of the table.
func (s *StatusTable) Push(res *trunner.TestResult) {
	s.Lock()
	defer s.Unlock()

	style := successStyle
	if res.Code != 0 {
		style = failureStyle
	}

	newStatuses := map[int]ui.Style{0: style}
	for i := 0; i < s.maxRows; i++ {
		newStatuses[i+1] = s.table.RowStyles[i]
	}
	s.table.RowStyles = newStatuses
	ui.Render(s.table)
}
