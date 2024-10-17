package generate

import (
	"fmt"
	"github.com/algorandfoundation/hack-tui/api"
	"github.com/algorandfoundation/hack-tui/ui/utils"
	"github.com/charmbracelet/bubbles/table"
	"sort"
)

func (m ViewModel) View() string {
	return fmt.Sprintf(
		"What is your algorand address?\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}

func makeRows(keys *[]api.ParticipationKey) *[]table.Row {
	rows := make([]table.Row, 0)
	for _, key := range *keys {
		rows = append(rows, table.Row{
			key.Id,
			key.Address,
			"TODO",
			utils.StrOrNA(key.EffectiveLastValid),
			"Unknown",
		})
	}
	sort.SliceStable(rows, func(i, j int) bool {
		return rows[i][0] < rows[j][0]
	})
	return &rows
}
