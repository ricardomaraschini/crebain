package basic

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ricardomaraschini/crebain/tui"
)

// New return a new instance of the basic user interface.
func New() *TUI {
	return &TUI{
		end: make(chan bool),
	}
}

// TUI is a basic text interfaces that just dumps test contents on the
// screen.
type TUI struct {
	end chan bool
}

// PushResult prints the test result on the screen.
func (t *TUI) PushResult(res tui.Drawable) {
	t.clearScreen()
	fmt.Println(res.Title())
	for _, line := range res.Content() {
		fmt.Println(line)
	}
}

// Start initiates the user interface.
func (t *TUI) Start() {
	t.clearScreen()
	<-t.end
}

// Close ends and closes the interface.
func (t *TUI) Close() {
	t.end <- true
}

func (t *TUI) clearScreen() {
	cmd := exec.Command("clear") //Linux example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()
}
