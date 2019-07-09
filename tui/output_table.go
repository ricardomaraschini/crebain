package main

import (
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

// NewOutputTable ...
func NewOutputTable(width, height int) *OutputTable {
	table := widgets.NewTable()
	table.FillRow = true
	table.Border = false
	table.RowSeparator = false
	table.SetRect(5, 0, width, height)

	return &OutputTable{
		maxLines: height - 2,
		table:    table,
	}
}

// OutputTable ...
type OutputTable struct {
	selRow   int
	maxLines int
	table    *widgets.Table
}

// Event ...
func (o *OutputTable) Event(id string) {
	switch id {
	case "j", "<Down>":
		if o.selRow == len(o.table.Rows)-1 {
			return
		}
		o.table.RowStyles[o.selRow] = normalRowStyle
		o.selRow++
		o.table.RowStyles[o.selRow] = selRowStyle
	case "k", "<Up>":
		if o.selRow == 0 {
			return
		}
		o.table.RowStyles[o.selRow] = normalRowStyle
		o.selRow--
		o.table.RowStyles[o.selRow] = selRowStyle
	}

	o.Render()
}

// Push ...
func (o *OutputTable) Push(content ...string) {
	rows := [][]string{content}
	o.table.Rows = append(rows, o.table.Rows...)
	if len(o.table.Rows) > o.maxLines {
		o.table.Rows = o.table.Rows[:o.maxLines]
	}
	if o.selRow > 0 && o.selRow < len(o.table.Rows)-1 {
		o.table.RowStyles[o.selRow] = normalRowStyle
		o.selRow++
	}

	o.table.RowStyles[o.selRow] = selRowStyle
	o.Render()
}

// Render ...
func (o *OutputTable) Render() {
	ui.Render(o.table)
}
