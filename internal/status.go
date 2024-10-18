package internal

import (
	"context"
	"fmt"
	"github.com/algorandfoundation/hack-tui/api"
	"time"
)

// StatusModel represents a status response from algod.Status
type StatusModel struct {
	Metrics     MetricsModel
	State       string
	Version     string
	Network     string
	Voting      bool
	NeedsUpdate bool
	LastRound   uint64 // Last recorded round
}

// String prints the last round value
func (m *StatusModel) String() string {
	return fmt.Sprintf("\nLastRound: %d\nRoundTime: %f \nTPS: %f", m.LastRound, m.Metrics.RoundTime.Seconds(), m.Metrics.TPS)
}

// Fetch handles algod.Status
func (m *StatusModel) Fetch(ctx context.Context, client *api.ClientWithResponses) error {
	if m.Version == "" || m.Version == "NA" {
		v, err := client.GetVersionWithResponse(ctx)
		if err != nil {
			return err
		}
		if v.StatusCode() != 200 {
			return fmt.Errorf("Status code %d: %s", v.StatusCode(), v.Status())
		}
		m.Network = v.JSON200.GenesisId
		m.Version = fmt.Sprintf("v%d.%d.%d-%s", v.JSON200.Build.Major, v.JSON200.Build.Minor, v.JSON200.Build.BuildNumber, v.JSON200.Build.Channel)

	}

	s, err := client.GetStatusWithResponse(ctx)
	if err != nil {
		return err
	}

	if s.StatusCode() != 200 {
		return fmt.Errorf("Status code %d: %s", s.StatusCode(), s.Status())
	}
	m.LastRound = uint64(s.JSON200.LastRound)

	if s.JSON200.UpgradeNodeVote != nil {
		m.Voting = *s.JSON200.UpgradeNodeVote
	}
	return nil
}

func getAverage(data []float64) float64 {
	sum := 0.0
	for _, element := range data {
		sum += element
	}
	return sum / (float64(len(data)))
}
func getAverageDuration(timings []time.Duration) time.Duration {
	sum := 0.0
	for _, element := range timings {
		sum += element.Seconds()
	}
	avg := sum / (float64(len(timings)))
	return time.Duration(avg * float64(time.Second))
}

// Watch uses WaitForBlockWithResponse to wait for changes and emits to the HeartBeat channel
func (m *StatusModel) Watch(cb func(model *StatusModel, err error), ctx context.Context, client *api.ClientWithResponses) {
	lastRound := m.LastRound
	timings := make([]time.Duration, 0)
	txns := make([]float64, 0)
	for {
		startTime := time.Now()
		status, err := client.WaitForBlockWithResponse(ctx, int(lastRound))
		endTime := time.Now()
		if err != nil {
			cb(nil, err)
		}
		var format api.GetBlockParamsFormat = "json"
		block, err := client.GetBlockWithResponse(ctx, int(lastRound), &api.GetBlockParams{
			Format: &format,
		})
		if err != nil {
			cb(nil, err)
		}
		m.LastRound = uint64(status.JSON200.LastRound)

		dur := endTime.Sub(startTime)
		timings = append(timings, dur)
		if block.JSON200.Block["txns"] != nil {
			txns = append(txns, float64(len(block.JSON200.Block["txns"].([]any)))/getAverageDuration(timings).Seconds())
		} else {
			txns = append(txns, 0)
		}

		m.Metrics.RoundTime = getAverageDuration(timings)
		m.Metrics.Window = len(timings)
		m.Metrics.TPS = getAverage(txns)

		// Trim data
		if len(timings) >= 100 {
			timings = timings[1:]
			txns = txns[1:]
		}

		lastRound = m.LastRound
		cb(m, nil)
	}
}
