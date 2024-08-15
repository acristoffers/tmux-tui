package tmux_tui

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	background          lipgloss.Color
	foreground          lipgloss.Color
	accent              lipgloss.Color
	secondary           lipgloss.Color
	selectionBackground lipgloss.Color
}

var DraculaTheme = Theme{
	background:          lipgloss.Color("#282A36"),
	foreground:          lipgloss.Color("#E3E3DE"),
	accent:              lipgloss.Color("#50FA7B"),
	secondary:           lipgloss.Color("#FFB86C"),
	selectionBackground: lipgloss.Color("#BD93F9"),
}
