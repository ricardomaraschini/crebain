package main

// SysError represents a system wide error.
//
// Complies with tui.Drawable interface so we can represent it on the interface.
type SysError struct {
	error
}

// Title returns a default title for system errors.
func (s *SysError) Title() string {
	return "system error"
}

// Content return the embed error content.
func (s *SysError) Content() []string {
	return []string{s.Error()}
}

// Success in case of errors, is always false.
func (s *SysError) Success() bool {
	return false
}
