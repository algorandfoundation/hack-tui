package accounts

import (
	"bytes"
	"github.com/algorandfoundation/hack-tui/internal"
	"github.com/algorandfoundation/hack-tui/ui/test"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/exp/golden"
	"github.com/charmbracelet/x/exp/teatest"
	"testing"
	"time"
)

func Test_New(t *testing.T) {
	m := New(&internal.StateModel{})
	acc := m.SelectedAccount()

	if acc != nil {
		t.Errorf("Expected no accounts to exist, got %s", acc.Address)
	}
	m, cmd := m.HandleMessage(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("enter"),
	})

	if cmd != nil {
		t.Errorf("Expected no comand")
	}

	m = New(test.GetState())
	m, _ = m.HandleMessage(tea.WindowSizeMsg{Width: 80, Height: 40})

	if m.Data.Admin {
		t.Errorf("Admin flag should be false, got true")
	}

	// Fetch state after message handling
	acc = m.SelectedAccount()
	if acc == nil {
		t.Errorf("expected true, got false")
	}

	// Update syncing state
	m.Data.Status.State = internal.SyncingState
	m.makeRows()
	if m.Data.Status.State != internal.SyncingState {

	}
}

func Test_Snapshot(t *testing.T) {
	t.Run("Visible", func(t *testing.T) {
		model := New(test.GetState())

		model, _ = model.HandleMessage(tea.WindowSizeMsg{Width: 80, Height: 40})
		got := ansi.Strip(model.View())
		golden.RequireEqual(t, []byte(got))
	})
}

func Test_Messages(t *testing.T) {
	// Create the Model
	m := New(test.GetState())

	tm := teatest.NewTestModel(
		t, m,
		teatest.WithInitialTermSize(80, 40),
	)

	// Wait for prompt to exit
	teatest.WaitFor(
		t, tm.Output(),
		func(bts []byte) bool {
			return bytes.Contains(bts, []byte("accounts | keys"))
		},
		teatest.WithCheckInterval(time.Millisecond*100),
		teatest.WithDuration(time.Second*3),
	)

	tm.Send(*test.GetState())

	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("enter"),
	})

	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("ctrl+c"),
	})

	tm.Send(tea.QuitMsg{})

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}