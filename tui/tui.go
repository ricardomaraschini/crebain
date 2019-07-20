package tui

import "github.com/ricardomaraschini/crebain/trunner"

// UI interface is implemented by a text based user interface or any other
// implementation that renders test results.
type UI interface {
	PushResult(res *trunner.TestResult)
	Start()
	Close()
}
