package accounts

import (
	"fmt"
)

func (m ViewModel) View() string {
	return fmt.Sprintf("%s\n%s", m.table.View(), m.controls.View())
}
