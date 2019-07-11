package tui

import (
	"github.com/ricardomaraschini/crebain/trunner"

	ui "github.com/gizak/termui/v3"
)

// New returns a new Terminal User Interface reference.
func New() (*TUI, error) {
	if err := ui.Init(); err != nil {
		return nil, err
	}

	width, height := ui.TerminalDimensions()
	t := &TUI{
		statusTable: NewStatusTable(height / 3),
		testDetail:  NewTestDetail(height/3, width, height),
		testsTable:  NewTestsTable(width, height/3),
	}
	t.testsTable.OnSelRowChange = t.testDetail.Push
	return t, nil
}

// TUI is our Terminal User Interface.
type TUI struct {
	statusTable *StatusTable
	testsTable  *TestsTable
	testDetail  *TestDetail
}

// Start initiates the Terminal User Interface.
//
// It only returns when the user requires to close the interface.
func (t *TUI) Start() {
	events := ui.PollEvents()
	for {
		e := <-events
		if e.ID == "q" || e.ID == "<C-c>" {
			break
		}
		t.testsTable.Event(e.ID)
		t.testDetail.Event(e.ID)
	}
	ui.Close()
}

// PushResult pushes a new test result into the interface.
func (t *TUI) PushResult(res *trunner.TestResult) {
	t.statusTable.Push(res)
	t.testsTable.Push(res)
}
