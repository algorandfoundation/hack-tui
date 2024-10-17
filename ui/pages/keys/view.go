package keys

import "fmt"

func (m ViewModel) View() string {
	//return m.table.View()
	return fmt.Sprintf("%s\n%s", m.table.View(), m.controls.View())
}
