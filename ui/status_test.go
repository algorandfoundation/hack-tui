package ui

import (
	"bytes"
	"context"
	"github.com/algorandfoundation/hack-tui/api"
	"github.com/algorandfoundation/hack-tui/internal"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"strings"
	"testing"
	"time"
)

func Test_InvalidStatusViewModel(t *testing.T) {
	client, err := api.NewClientWithResponses("http://255.255.255.255:4001")
	if err != nil {
		t.Fatal(err)
	}

	// Test Invalid Node
	_, err = MakeStatusViewModel(context.Background(), client)
	if !strings.Contains(err.Error(), "dial tcp 255.255.255.255:4001") {
		t.Fatal(err)
	}
}
func Test_StatusViewRender(t *testing.T) {
	status := internal.StatusModel{
		HeartBeat:   make(chan uint64),
		LastRound:   0,
		NeedsUpdate: true,
		State:       "SYNCING",
	}
	// Create the Model
	m := StatusViewModel{
		Status: &status,
		Metrics: &internal.MetricsModel{
			RoundTime: 0,
			TX:        0,
			RX:        0,
			TPS:       0,
		},
		ViewWidth: 80,
		IsVisible: true,
	}

	tm := teatest.NewTestModel(
		t, m,
		teatest.WithInitialTermSize(80, 40),
	)

	// Wait for prompt to exit
	teatest.WaitFor(
		t, tm.Output(),
		func(bts []byte) bool {
			return bytes.Contains(bts, []byte("Latest Round: 0"))
		},
		teatest.WithCheckInterval(time.Millisecond*100),
		teatest.WithDuration(time.Second*3),
	)

	// Send hide key
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("h"),
	})

	// Send quit key
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("q"),
	})

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}
