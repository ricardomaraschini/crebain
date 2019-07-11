package tui

import (
	"strings"
	"sync"

	"github.com/ricardomaraschini/crebain/trunner"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// NewTestDetail returns a new widget capable of rendering a text.
func NewTestDetail(x, width, height int) *TestDetail {
	box := widgets.NewParagraph()
	box.SetRect(0, x, width, height)
	return &TestDetail{
		box:     box,
		maxRows: height - x - 2,
	}
}

// TestDetail the box where we add the text.
type TestDetail struct {
	sync.Mutex
	maxRows         int
	fullContent     []string
	renderedContent []string
	box             *widgets.Paragraph
}

// Event is called everytime the user generates an event, e.g. presses a key.
func (r *TestDetail) Event(event string) {
	switch event {
	case "J":
		r.box.Text = "down"
	case "K":
		r.box.Text = "up"
	case "H":
		r.box.Text = "left"
	case "L":
		r.box.Text = "right"
	}
	ui.Render(r.box)
}

// Push sets the title and text to be rendered.
func (r *TestDetail) Push(res *trunner.TestResult) {
	r.Lock()
	defer r.Unlock()

	r.fullContent = make([]string, len(res.Out))
	r.renderedContent = make([]string, r.maxRows)
	count := 0
	for _, out := range res.Out {
		if out.Output == "" {
			continue
		}
		r.fullContent[count] = out.Output
		if count < r.maxRows {
			r.renderedContent[count] = out.Output
		}
		count++
	}

	//r.box.Title = title
	r.box.Text = strings.Join(r.renderedContent, "")
	//r.box.Text = strings.Join(r.fullContent, "")
	ui.Render(r.box)
}
