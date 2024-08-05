package tmux_tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func Pane(title string, width int, height int, selected bool, contents string) string {
	border := lipgloss.RoundedBorder()

	// Account for paddings and header
	width -= 2
	height -= 2

	titleSize := width - 3 - lipgloss.Width(title)
	if titleSize < 0 {
		titleSize = width - 3
		title = title[:titleSize-5] + "..."
	}
	headerLeft := border.TopLeft + border.Top + " "
	headerRight := " " + strings.Repeat(border.Top, width-3-len(title)) + border.TopRight
	header := lipgloss.NewStyle().
		SetString(headerLeft + title + headerRight)

	pane := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), false, true, true, true).
		PaddingLeft(1).
		PaddingRight(1).
		Height(height).
		Width(width).
		SetString(contents)

	if selected {
		borderColor := lipgloss.Color("2")
		header = header.Foreground(borderColor)
		pane = pane.BorderForeground(borderColor)
	}

	return lipgloss.JoinVertical(lipgloss.Top, header.String(), pane.String())
}
