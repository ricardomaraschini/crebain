package main

import (
	"fmt"
	"log"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type OutputTable struct {
	rows [][]string
}

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	width, height := ui.TerminalDimensions()
	go func() {
		tbl := NewStatusTable(height / 3)
		i := 0
		for {
			tbl.Push(i%2 == 0)
			i++
			time.Sleep(time.Second)
		}
	}()

	selStyle := ui.NewStyle(ui.ColorBlack, ui.ColorWhite)
	table := widgets.NewTable()
	table.Rows = [][]string{
		[]string{"Last Execution", "Package", "Coverage"},
		[]string{"05:46:26", "github.com/ricardomaraschini/crebain/cmd", "100%"},
		[]string{"05:46:26", "github.com/ricardomaraschini/crebain/cmd", "100%"},
		[]string{"05:46:26", "github.com/ricardomaraschini/crebain/cmd", "100%"},
		[]string{"05:46:26", "github.com/ricardomaraschini/crebain/cmd", "100%"},
		[]string{"05:46:26", "github.com/ricardomaraschini/crebain/cmd", "100%"},
		[]string{"05:46:26", "github.com/ricardomaraschini/crebain/cmd", "100%"},
		[]string{"05:46:26", "github.com/ricardomaraschini/crebain/cmd", "100%"},
		[]string{"05:46:26", "github.com/ricardomaraschini/crebain/cmd", "100%"},
		[]string{"05:46:26", "github.com/ricardomaraschini/crebain/cmd", "100%"},
		[]string{"05:46:26", "github.com/ricardomaraschini/crebain/cmd", "100%"},
		[]string{"05:46:26", "github.com/ricardomaraschini/crebain/cmd", "100%"},
	}
	table.RowStyles[1] = selStyle
	table.FillRow = true
	table.Border = false
	table.RowSeparator = false
	table.TextStyle = ui.NewStyle(ui.ColorWhite)
	table.BorderStyle = ui.NewStyle(ui.ColorWhite)
	table.SetRect(5, 0, width-3, height/3)

	result := NewResultBox(height/3, width, height)
	result.Render("title", "content")
	ui.Render(table)

	curpos := 1
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			if curpos == len(table.Rows)-1 {
				continue
			}
			table.RowStyles[curpos] = ui.NewStyle(ui.ColorClear)
			curpos++
			table.RowStyles[curpos] = selStyle
			ui.Render(table)
			result.Render(fmt.Sprintf("%d", curpos), fmt.Sprintf("%d", curpos))
		case "k", "<Up>":
			if curpos == 1 {
				continue
			}
			table.RowStyles[curpos] = ui.NewStyle(ui.ColorClear)
			curpos--
			table.RowStyles[curpos] = selStyle
			ui.Render(table)
			result.Render(fmt.Sprintf("%d", curpos), fmt.Sprintf("%d", curpos))
		default:
			table.Rows[0][0] = e.ID
			ui.Render(table)
		}
	}
}