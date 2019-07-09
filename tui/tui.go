package tui

import (
	"github.com/ricardomaraschini/crebain/trunner"

	ui "github.com/gizak/termui/v3"
)

// Event processor is an entity capable of processing user events.
type EventProcessor interface {
	Event(string)
}

var (
	statusTable *StatusTable
	testsTable  *OutputTable
	testDetail  *ResultBox
)

func NewTestResult(res *trunner.TestResult) {
	statusTable.Push(res.Code == 0)
	pkg := "undefined"
	if len(res.Out) > 0 {
		pkg = res.Out[0].Package
	}
	testsTable.Push(pkg)

	content := ""
	for _, out := range res.Out {
		content += out.Output
	}
	testDetail.Set(pkg, content)
}

// StartTUI initiates our text user interface.
func StartTUI() error {
	if err := ui.Init(); err != nil {
		return err
	}

	width, height := ui.TerminalDimensions()

	statusTable = NewStatusTable(height / 3)
	testDetail = NewResultBox(height/3, width, height)
	testsTable = NewOutputTable(width, height/3)
	testsTable.OnSelRowChange = func(idx int) {
		/*
			testDetail.Set(
				fmt.Sprintf("%d", idx),
				fmt.Sprintf("%d", idx),
			)
		*/
	}

	/*
		go func() {
			time.Sleep(5 * time.Second)
			for i := 0; i < 1000; i++ {
				testsTable.Push(
					fmt.Sprintf("%d", i),
					fmt.Sprintf("%d", i),
					fmt.Sprintf("%d", i),
					fmt.Sprintf("%d", i),
					fmt.Sprintf("%d", i),
				)
				statuses.Push(i%2 == 0)
				time.Sleep(time.Second)
			}
		}()
	*/

	events := ui.PollEvents()
	for {
		e := <-events
		testsTable.Event(e.ID)
		if e.ID == "q" || e.ID == "<C-c>" {
			break
		}
	}
	ui.Close()
	return nil
}
