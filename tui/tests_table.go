package tui

import (
	"sync"
	"time"

	"github.com/ricardomaraschini/crebain/trunner"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var (
	normalRowStyle = ui.Style{
		Fg: ui.ColorClear,
		Bg: ui.ColorClear,
	}
	selRowStyle = ui.Style{
		Fg: ui.ColorWhite,
		Bg: ui.Color(17),
	}
)

// NewTestsTable returns a terminal ui component capable of rendering a
// table were we present the test outputs, one per row.
func NewTestsTable(width, height int) *TestsTable {
	table := widgets.NewTable()
	table.FillRow = true
	table.Border = false
	table.RowSeparator = false
	table.SetRect(5, 0, width, height)

	return &TestsTable{
		OnSelRowChange: func(*trunner.TestResult) {},
		maxRows:        height - 2,
		table:          table,
	}
}

// TestsTable is an ui component for rendering go test results.
//
// Allow user to navigate through test results.
type TestsTable struct {
	sync.Mutex
	testResults    []*trunner.TestResult
	OnSelRowChange func(*trunner.TestResult)
	selRow         int
	maxRows        int
	table          *widgets.Table
}

// SelectedRow returns the currently selected row index.
func (o *TestsTable) SelectedRow() int {
	o.Lock()
	defer o.Unlock()
	return o.selRow
}

// Event is called everytime the user takes an action, e.g. presses a key.
func (o *TestsTable) Event(event string) {
	o.Lock()
	defer o.Unlock()

	switch event {
	case "j", "<Down>":
		if o.selRow == len(o.table.Rows)-1 {
			return
		}
		o.table.RowStyles[o.selRow] = normalRowStyle
		o.selRow++
		o.OnSelRowChange(o.testResults[o.selRow])
		o.table.RowStyles[o.selRow] = selRowStyle
	case "k", "<Up>":
		if o.selRow == 0 {
			return
		}
		o.table.RowStyles[o.selRow] = normalRowStyle
		o.selRow--
		o.OnSelRowChange(o.testResults[o.selRow])
		o.table.RowStyles[o.selRow] = selRowStyle
	default:
		return
	}

	ui.Render(o.table)
}

// Push adds a new row to the begining of the table.
func (o *TestsTable) Push(res *trunner.TestResult) {
	o.Lock()
	defer o.Unlock()

	pkg := "undefined"
	if len(res.Out) > 0 {
		pkg = res.Out[0].Package
	}

	rows := [][]string{
		[]string{time.Now().Format("Mon Jan 2 15:04:05"), pkg},
	}
	results := []*trunner.TestResult{res}

	o.table.Rows = append(rows, o.table.Rows...)
	o.testResults = append(results, o.testResults...)
	if len(o.table.Rows) > o.maxRows {
		o.table.Rows = o.table.Rows[:o.maxRows]
		o.testResults = o.testResults[:o.maxRows]
	}

	// if the selected row is not the first neither the last row we
	// make the selected row move down one row in order to maintain
	// the same row selected.
	if o.selRow != 0 && o.selRow != len(o.table.Rows)-1 {
		o.table.RowStyles[o.selRow] = normalRowStyle
		o.selRow++
	}

	o.OnSelRowChange(o.testResults[o.selRow])
	o.table.RowStyles[o.selRow] = selRowStyle
	ui.Render(o.table)
}
