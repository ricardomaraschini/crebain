package tui

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// NewResultBox returns a new widget capable of rendering a text.
func NewResultBox(x, width, height int) *ResultBox {
	box := widgets.NewParagraph()
	box.SetRect(0, x, width, height)
	return &ResultBox{
		box: box,
	}
}

// ResultBox the box where we add the text.
type ResultBox struct {
	box *widgets.Paragraph
}

// Event is called everytime the user generates an event, e.g. presses a key.
func (r *ResultBox) Event(event string) {}

// Set sets the title and text to be rendered.
func (r *ResultBox) Set(title, content string) {
	r.box.Title = title
	r.box.Text = content
	ui.Render(r.box)
}
