package tui

// Drawable represents items we are able to draw on the interface.
type Drawable interface {
	Success() bool
	Title() string
	Content() []string
}

// UI interface is implemented by a text based user interface or any other
// implementation that renders test results.
type UI interface {
	PushResult(Drawable)
	Start()
	Close()
}
