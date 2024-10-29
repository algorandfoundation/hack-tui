package generate

import (
	"fmt"
	"github.com/algorandfoundation/hack-tui/ui/pages"
	"strings"
)

func (m ViewModel) View() string {
	var b strings.Builder

	for i := range m.Inputs {
		b.WriteString(m.Inputs[i].View())
		if i < len(m.Inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.Inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	return pages.WithTitle("Generate", pages.PageBorder(m.Width-3).Render(b.String()))
}
