package tmux_tui

import (
	"fmt"
	"os/exec"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) initWindows() Model {
	return m
}

func (m Model) updateWindows(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.focusedPane != 2 {
			return nil, nil
		}
		switch msg.String() {
		case "ctrl+p", "k", tea.KeyUp.String():
			windows := m.visibleWindows()
			index := slices.IndexFunc(windows, func(e entity) bool { return e.id == m.focusedWindowId })
			if index > 0 && index < len(windows) {
				m.focusedWindowId = windows[index-1].id
			} else {
				m.focusedWindowId = windows[0].id
			}
			return m, tea.Batch(listEntitiesCmd, previewCmd(m))
		case "ctrl+n", "j", tea.KeyDown.String():
			windows := m.visibleWindows()
			index := slices.IndexFunc(windows, func(e entity) bool { return e.id == m.focusedWindowId })
			if index > -1 && index < len(windows)-1 {
				m.focusedWindowId = windows[index+1].id
			} else if index != len(windows)-1 {
				m.focusedWindowId = windows[0].id
			}
			return m, tea.Batch(listEntitiesCmd, previewCmd(m))
		}
	}

	return nil, nil
}

func (m Model) viewWindows() string {
	var windows []string

	for _, window := range m.windows {
		if !m.showAll && window.parent != m.focusedSessionId {
			continue
		}
		panesString := lipgloss.NewStyle().
			Width(m.windowWidth/3 - 4).
			MaxWidth(m.windowWidth/3 - 4).
			SetString(fmt.Sprintf("%d: %s", window.id, window.name))
		if m.appState == Swapping && window.id == m.swapSrc && m.focusedPane == 2 {
			panesString = panesString.
				Background(lipgloss.Color("4")).
				Foreground(lipgloss.Color("15")).
				SetString(fmt.Sprintf("[src] %d: %s", window.id, window.name))
		}
		if window.id == m.focusedWindowId {
			panesString = panesString.Foreground(lipgloss.Color("2"))
		}
		windows = append(windows, panesString.String())
	}

	return lipgloss.JoinVertical(lipgloss.Top, windows...)
}

func goToWindow(m Model) tea.Cmd {
	return func() tea.Msg {
		window := m.windowWithId(m.focusedWindowId)
		c := exec.Command("tmux",
			"switch-client", "-t", fmt.Sprintf("$%d", window.parent), ";",
			"select-window", "-t", fmt.Sprintf("@%d", window.id))
		c.Run()
		return tea.QuitMsg{}
	}
}

func renameWindowCmd(m Model) tea.Cmd {
	return func() tea.Msg {
		c := exec.Command("tmux", "rename-window", "-t", fmt.Sprintf("@%d", m.focusedWindowId), m.textInput.Value())
		c.Run()
		return clearInputTextMsg{}
	}
}

func newWindowCmd(m Model) tea.Cmd {
	return func() tea.Msg {
		c := exec.Command("tmux", "new-window", "-n", m.textInput.Value(), "-t", fmt.Sprintf("%s:", m.sessionWithId(m.focusedSessionId).name))
		c.Run()
		return clearInputTextMsg{}
	}
}

func deleteWindowCmd(m Model) tea.Cmd {
	return func() tea.Msg {
		c := exec.Command("tmux", "kill-window", "-t", fmt.Sprintf("@%d", m.focusedWindowId))
		c.Run()
		return tickMsg{}
	}
}

func swapWindowsCmd(m Model) tea.Cmd {
	return func() tea.Msg {
		src := fmt.Sprintf("@%d", m.swapSrc)
		dst := fmt.Sprintf("@%d", m.focusedWindowId)
		c := exec.Command("tmux", "swap-window", "-s", src, "-t", dst)
		c.Run()
		return tickMsg{}
	}
}
