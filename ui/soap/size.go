package soap

import (
	tea "github.com/charmbracelet/bubbletea"
)

type WindowSizeMsg struct {
	Height int
	Width  int
}

func WithWindowSizeMsg(width int, height int) tea.Cmd {
	return func() tea.Msg {
		return WindowSizeMsg{
			Height: height,
			Width:  width,
		}
	}
}
