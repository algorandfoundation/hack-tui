package installer

import (
	"github.com/algorandfoundation/algorun-tui/ui/style"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error
type Model struct {
	Width        int
	Height       int
	confirmation chan bool
	spinner      spinner.Model
	quitting     bool
	err          error
}

func New() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return Model{spinner: s, confirmation: make(chan bool)}
}
func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case errMsg:
		m.err = msg
		return m, nil
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m Model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	//str := fmt.Sprintf("\n\n   %s Loading forever...press q to quit\n\n", m.spinner.View())
	//if m.quitting {
	//	return str + "\n"
	//}

	return style.WithTitle(
		"Install Confirmation",
		style.
			ApplyBorder(m.Width-2, m.Height-2, "24").
			Render("Are you sure you want to install this node?"),
	)
}
