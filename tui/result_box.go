package main

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// NewResultBox returns a new widget capable of rendering a text.
func NewResultBox(x, sWidth, sHeight int) *ResultBox {
	box := widgets.NewParagraph()
	box.SetRect(0, x, sWidth, sHeight)
	return &ResultBox{
		box: box,
	}
}

// ResultBox renders a square that extends from x to the total of the
// available screen.
type ResultBox struct {
	box *widgets.Paragraph
}

// Render renders the result box on the screen.
func (r *ResultBox) Render(title, content string) {
	r.box.Title = title
	r.box.Text = content
	ui.Render(r.box)
}
