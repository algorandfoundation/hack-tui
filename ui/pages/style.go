package pages

import (
	"github.com/charmbracelet/lipgloss"
	"strings"
)

var (
	Padding1 = lipgloss.NewStyle().Padding(1).Render
	Border   = func() lipgloss.Style {
		return lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder())
	}()
	PageBorder = func(width int) lipgloss.Style {
		return Border.
			Width(width).
			Padding(0).
			Margin(0).
			BorderForeground(lipgloss.Color("8"))
	}
)

func WithTitle(title string, view string) string {
	r := []rune(view)
	if lipgloss.Width(view) >= len(title)+4 {
		b, _, _, _, _ := Border.GetBorder()
		id := strings.IndexRune(view, []rune(b.Top)[0])
		start := string(r[0:id])
		return start + title + string(r[len(title)+id:])
	}
	return view
}
