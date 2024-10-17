package ui

import (
	"context"
	"fmt"
	"github.com/algorandfoundation/hack-tui/api"
	"github.com/algorandfoundation/hack-tui/ui/pages/accounts"
	"github.com/algorandfoundation/hack-tui/ui/pages/generate"
	"github.com/algorandfoundation/hack-tui/ui/pages/keys"
	"github.com/algorandfoundation/hack-tui/ui/pages/transaction"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ViewportPage string

const (
	AccountsPage    ViewportPage = "accounts"
	KeysPage        ViewportPage = "keys"
	GeneratePage    ViewportPage = "generate"
	TransactionPage ViewportPage = "transaction"
)

type ViewportViewModel struct {
	PageWidth, PageHeight         int
	TerminalWidth, TerminalHeight int
	// Header Components
	status   StatusViewModel
	protocol ProtocolViewModel

	// Pages
	accountsPage    accounts.ViewModel
	keysPage        keys.ViewModel
	generatePage    generate.ViewModel
	transactionPage transaction.ViewModel

	// Viewport Statue
	ready bool
	//viewport              viewport.Model
	viewportPage          ViewportPage
	viewportStatusChannel chan ViewportPage
}

// Init is a no-op
func (m ViewportViewModel) Init() tea.Cmd {
	return nil
}

// Update Handle the viewport lifecycle
func (m ViewportViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	// Handle Header Updates
	m.protocol, cmd = m.protocol.HandleMessage(msg)
	cmds = append(cmds, cmd)
	m.status, cmd = m.status.HandleMessage(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// TODO: Disable these handlers
		switch msg.String() {
		case "h":
			m.PageHeight = max(0, m.TerminalHeight-lipgloss.Height(m.headerView()))
		case "a":
			m.viewportPage = AccountsPage
		case "k":
			m.keysPage.Address = m.accountsPage.SelectedAccount()
			m.viewportPage = KeysPage
		case "t":
			m.viewportPage = TransactionPage
		//case "g":
		//	m.viewportPage = GeneratePage
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.TerminalWidth = msg.Width
		m.TerminalHeight = msg.Height

		m.PageWidth = msg.Width
		m.PageHeight = max(0, msg.Height-lipgloss.Height(m.headerView()))

		m.accountsPage.ViewHeight = m.PageHeight
		m.accountsPage.ViewWidth = m.PageWidth
		m.keysPage.ViewHeight = m.PageHeight
		m.keysPage.ViewWidth = m.PageWidth
		m.transactionPage.ViewHeight = m.PageHeight
		m.transactionPage.ViewWidth = m.PageWidth
	}
	// Get Page Updates
	switch m.viewportPage {
	case AccountsPage:
		m.accountsPage, cmd = m.accountsPage.HandleMessage(msg)
	case KeysPage:
		m.keysPage, cmd = m.keysPage.HandleMessage(msg)
	case GeneratePage:
		m.generatePage, cmd = m.generatePage.HandleMessage(msg)
	case TransactionPage:
		m.transactionPage, cmd = m.transactionPage.HandleMessage(msg)
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

// View renders the viewport.Model
func (m ViewportViewModel) View() string {
	// Handle Page render
	var page tea.Model
	switch m.viewportPage {
	case AccountsPage:
		page = m.accountsPage
	case GeneratePage:
		page = m.generatePage
	case KeysPage:
		page = m.keysPage
	case TransactionPage:
		page = m.transactionPage
	}
	return fmt.Sprintf("%s\n%s", m.headerView(), page.View())
	if m.headerView() != "" {
		//return m.headerView()
		return fmt.Sprintf("%s\n%s", m.headerView(), page.View())
	}
	return page.View()
}

// headerView generates the top elements
func (m ViewportViewModel) headerView() string {
	if m.TerminalHeight < 15 {
		return ""
	}

	if m.TerminalWidth < 90 {
		if m.protocol.View() == "" {
			return lipgloss.JoinVertical(lipgloss.Center, m.status.View())
		}
		return lipgloss.JoinVertical(lipgloss.Center, m.status.View(), m.protocol.View())
	}

	return lipgloss.JoinHorizontal(lipgloss.Center, m.status.View(), m.protocol.View())
}

// MakeViewportViewModel handles the construction of the TUI viewport
func MakeViewportViewModel(ctx context.Context, client *api.ClientWithResponses) (*ViewportViewModel, error) {

	status, err := MakeStatusViewModel(ctx, client)
	if err != nil {
		return nil, err
	}
	protocol := MakeProtocolViewModel(status.Status)

	ap, err := accounts.New(ctx, client)
	if err != nil {
		return nil, err
	}

	kp, err := keys.New(ctx, client)
	if err != nil {
		return nil, err
	}
	kp.Address = ap.SelectedAccount()
	m := ViewportViewModel{
		status:          status,
		protocol:        protocol,
		accountsPage:    ap,
		keysPage:        kp,
		generatePage:    generate.New(ctx, client),
		transactionPage: transaction.New(ctx, client),
		viewportPage:    AccountsPage,
	}

	return &m, nil
}
