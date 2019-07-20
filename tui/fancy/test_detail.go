package fancy

import (
	"strings"
	"sync"

	"github.com/ricardomaraschini/crebain/tui"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// NewTestDetail returns a new widget capable of rendering a test output.
func NewTestDetail(y0, x1, y1 int) *TestDetail {
	box := widgets.NewParagraph()
	box.SetRect(0, y0, x1, y1)
	return &TestDetail{
		box:     box,
		maxRows: y1 - y0 - 2,
		maxCols: x1 - 2,
	}
}

// TestDetail the box where we add the text.
type TestDetail struct {
	sync.Mutex
	firstRenderedRow int
	firstRenderedCol int
	maxRows          int
	maxCols          int
	fullContent      []string
	renderedContent  []string
	box              *widgets.Paragraph
}

// Event is called everytime the user generates an event, e.g. presses a key.
func (r *TestDetail) Event(event string) {
	r.Lock()
	defer r.Unlock()

	switch event {
	case "J":
		r.down()
	case "K":
		r.up()
	case "H":
		r.left()
	case "L":
		r.right()
	}
}

// up moves box content one row up.
func (r *TestDetail) up() {
	if r.firstRenderedRow == 0 {
		return
	}
	r.firstRenderedRow--
	r.writeContent()
}

// left moves the box content one char to the left.
func (r *TestDetail) left() {
	if r.firstRenderedCol == 0 {
		return
	}
	r.firstRenderedCol--
	r.writeContent()
}

// right moves the box content one char to the right.
func (r *TestDetail) right() {
	r.firstRenderedCol++
	r.writeContent()
}

// down moves the box content one row down.
func (r *TestDetail) down() {
	r.firstRenderedRow++
	r.writeContent()
}

// Update sets the title and text to be rendered.
func (r *TestDetail) Update(res tui.Drawable) {
	r.Lock()
	defer r.Unlock()

	r.fullContent = make([]string, 0)
	for _, line := range res.Content() {
		line = strings.Replace(line, "\t", " ", -1)
		r.fullContent = append(r.fullContent, line)
	}

	r.box.Title = res.Title()
	r.firstRenderedRow = 0
	r.firstRenderedCol = 0
	r.writeContent()
}

// writeContent writes content on the screen.
//
// Takes care to render only the portion of the output we can present on the interface.
func (r *TestDetail) writeContent() {
	lastLine := r.firstRenderedRow + r.maxRows
	r.renderedContent = make([]string, 0)
	for i := r.firstRenderedRow; i < lastLine; i++ {
		// nothing more to write.
		if i >= len(r.fullContent) {
			break
		}
		r.renderedContent = append(
			r.renderedContent,
			r.cutLine(r.fullContent[i]),
		)
	}
	r.box.Text = strings.Join(r.renderedContent, "\n")
	ui.Render(r.box)
}

// cutLine returns the part of the line from firstRenderedCol with maxCols columns.
func (r *TestDetail) cutLine(line string) string {
	// we are rendering beyond the end of this line.
	if r.firstRenderedCol >= len(line) {
		return ""
	}

	lastRenderedCol := r.firstRenderedCol + r.maxCols
	if lastRenderedCol >= len(line) {
		return line[r.firstRenderedCol:]
	}

	return line[r.firstRenderedCol:lastRenderedCol]
}
