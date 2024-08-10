package tmux_tui

import (
	"fmt"
	"os/exec"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) initPanes() Model {
	return m
}

func (m Model) updatePanes(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.focusedPane != 3 {
			return nil, nil
		}
		switch msg.String() {
		case "ctrl+p", "k", tea.KeyUp.String():
			panes := m.visiblePanes()
			index := slices.IndexFunc(panes, func(e entity) bool { return e.id == m.focusedPaneId })
			if index > 0 && index < len(panes) {
				m.focusedPaneId = panes[index-1].id
			} else {
				m.focusedPaneId = panes[0].id
			}
			return m, tea.Batch(listEntitiesCmd, previewCmd(m))
		case "ctrl+n", "j", tea.KeyDown.String():
			panes := m.visiblePanes()
			index := slices.IndexFunc(panes, func(e entity) bool { return e.id == m.focusedPaneId })
			if index > -1 && index < len(panes)-1 {
				m.focusedPaneId = panes[index+1].id
			} else if index != len(panes)-1 {
				m.focusedPaneId = panes[0].id
			}
			return m, tea.Batch(listEntitiesCmd, previewCmd(m))
		}
	}

	return nil, nil
}

func (m Model) viewPanes() string {
	var panes []string

	for _, pane := range m.panes {
		if !m.showAll && pane.parent != m.focusedWindowId {
			continue
		}
		panesString := lipgloss.NewStyle().
			Width(m.windowWidth/3 - 4).
			MaxWidth(m.windowWidth/3 - 4).
			SetString(fmt.Sprintf("%d", pane.id))
		if pane.id == m.focusedPaneId {
			panesString = panesString.Foreground(lipgloss.Color("2"))
		}
		if m.appState == Swapping && pane.id == m.swapSrc && m.focusedPane == 3 {
			panesString = panesString.Background(lipgloss.Color("6")).Foreground(lipgloss.Color("0"))
		}
		panes = append(panes, panesString.String())
	}

	return lipgloss.JoinVertical(lipgloss.Top, panes...)
}

func goToPane(m Model) tea.Cmd {
	return func() tea.Msg {
		pane := m.paneWithId(m.focusedPaneId)
		window := m.windowWithId(pane.parent)
		c := exec.Command("tmux",
			"switch-client", "-t", fmt.Sprintf("$%d", window.parent), ";",
			"select-window", "-t", fmt.Sprintf("@%d", window.id), ";",
			"select-pane", "-t", fmt.Sprintf("%%%d", pane.id))
		c.Run()
		return tea.QuitMsg{}
	}
}

func deletePaneCmd(m Model) tea.Cmd {
	return func() tea.Msg {
		c := exec.Command("tmux", "kill-pane", "-t", fmt.Sprintf("%%%d", m.focusedPaneId))
		c.Run()
		return tickMsg{}
	}
}

func swapPanesCmd(m Model) tea.Cmd {
	return func() tea.Msg {
		src := fmt.Sprintf("%%%d", m.swapSrc)
		dst := fmt.Sprintf("%%%d", m.focusedPaneId)
		c := exec.Command("tmux", "swap-pane", "-s", src, "-t", dst)
		c.Run()
		return tickMsg{}
	}
}
