package transaction

import (
	"github.com/charmbracelet/lipgloss"
	"strings"
)

func (m ViewModel) View() string {
	if m.asciiQR == "" || m.urlTxn == "" {
		return lipgloss.JoinVertical(
			lipgloss.Center,
			"No QR Code or TxnURL available",
			"\n",
			m.controls.View())
	}

	// Build QR Parts
	hint := yellow.Render(m.hint)
	qrCode := qrStyle.Render(m.asciiQR)
	url := urlStyle.Render(strings.Replace(m.urlTxn, m.Data.Address, m.FormatedAddress(), 1))

	controls := m.controls.View()

	qrFullRender := lipgloss.JoinVertical(
		lipgloss.Center,
		hint,
		"",
		qrCode,
		url,
	)

	remainingHeight := max(0, m.Height-lipgloss.Height(controls))
	isLargeScreen := lipgloss.Height(qrFullRender) <= remainingHeight && lipgloss.Width(qrFullRender) <= m.Width
	isSmallScreen := lipgloss.Height(qrCode) <= remainingHeight && lipgloss.Width(qrCode) < m.Width

	if isLargeScreen {
		qrRenderPadHeight := max(0, remainingHeight-2-lipgloss.Height(qrFullRender))
		qrPad := strings.Repeat("\n", max(0, qrRenderPadHeight/2))
		if qrRenderPadHeight > 2 {
			return lipgloss.JoinVertical(
				lipgloss.Center,
				qrPad,
				qrFullRender,
				qrPad,
				m.controls.View(),
			)
		}
		return lipgloss.JoinVertical(
			lipgloss.Center,
			qrFullRender,
			m.controls.View(),
		)
	}
	if isSmallScreen {
		isQrPadded := lipgloss.Height(qrCode) < remainingHeight && lipgloss.Width(qrCode) < m.Width
		if isQrPadded {
			qrRenderPadHeight := max(0, remainingHeight-2-lipgloss.Height(qrCode))
			qrPad := strings.Repeat("\n", max(0, qrRenderPadHeight/2))
			return lipgloss.JoinVertical(
				lipgloss.Center,
				qrPad,
				qrCode,
				qrPad,
				controls)
		}
		return lipgloss.JoinVertical(
			lipgloss.Center,
			qrCode,
			controls)
	}

	padHeight := max(0, remainingHeight-lipgloss.Height(controls))
	padHString := strings.Repeat("\n", padHeight/2)
	text := red.Render("QR Code too large to display... Please adjust terminal dimensions or font.")
	padWidth := max(0, m.Width-lipgloss.Width(text))
	padWString := strings.Repeat(" ", padWidth/2)
	return lipgloss.JoinVertical(
		lipgloss.Left,
		padHString,
		lipgloss.JoinHorizontal(lipgloss.Left, padWString, red.Render("QR Code too large to display... Please adjust terminal dimensions or font.")),
		padHString,
		m.controls.View(),
	)
}
