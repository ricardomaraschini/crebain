package tui

import (
	"fmt"
	"sync"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// New returns a new text user interface manager.
func New() (*TUI, error) {
	if err := ui.Init(); err != nil {
		return nil, err
	}

	width, height := ui.TerminalDimensions()

	tabPane := widgets.NewTabPane()
	tabPane.SetRect(0, 0, width, 3)
	tabPane.Border = true

	ti := &TUI{
		tabs:     tabPane,
		txtBoxes: make([]*widgets.Paragraph, 0),
		width:    width,
		height:   height,
	}

	go ti.eventLoop()
	return ti, nil
}

// TUI controls our text user interface.
type TUI struct {
	sync.Mutex
	tabs     *widgets.TabPane
	txtBoxes []*widgets.Paragraph
	width    int
	height   int
}

// eventLoop awaits for keys to be pressed.
func (t *TUI) eventLoop() {
	uiEvents := ui.PollEvents()
	for {
		e, ok := <-uiEvents
		if !ok {
			return
		}

		switch e.ID {
		case "q", "<C-c>":
			t.Stop()
			return
		case "h":
			t.tabs.FocusLeft()
		case "l":
			t.tabs.FocusRight()
		}

		ui.Clear()
		ui.Render(t.tabs)
		t.renderTXTBox()
	}

}

func (t *TUI) renderTXTBox() {
	t.Lock()
	defer t.Unlock()

	active := t.tabs.ActiveTabIndex
	last := len(t.txtBoxes) - 1
	if active > last {
		panic("rendering an invalid tab")
	}

	ui.Render(t.txtBoxes[active])
}

// Stop ends the text user interface.
func (t *TUI) Stop() {
	ui.Close()
}

// AddBox adds a new box to the text interface.
func (t *TUI) AddBox(content fmt.Stringer) {
	t.Lock()
	defer t.Unlock()

	box := widgets.NewParagraph()
	box.Text = content.String()
	box.Title = "title"
	box.SetRect(0, 3, t.width, t.height)
	box.BorderStyle.Fg = ui.ColorGreen

	t.tabs.TabNames = append(t.tabs.TabNames, "tabname")
	t.txtBoxes = append(t.txtBoxes, box)
}
