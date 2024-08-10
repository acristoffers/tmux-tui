package tmux_tui

import (
	"fmt"
	"os/exec"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) initSessions() Model {
	return m
}

func (m Model) updateSessions(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.focusedPane != 1 {
			return nil, nil
		}
		switch msg.String() {
		case "ctrl+p", "k", tea.KeyUp.String():
			sessions := m.sessions
			index := slices.IndexFunc(sessions, func(e entity) bool { return e.id == m.focusedSessionId })
			if index > 0 && index < len(sessions) {
				m.focusedSessionId = sessions[index-1].id
			} else {
				m.focusedSessionId = sessions[0].id
			}
			return m, tea.Batch(listEntitiesCmd, previewCmd(m))
		case "ctrl+n", "j", tea.KeyDown.String():
			sessions := m.sessions
			index := slices.IndexFunc(sessions, func(e entity) bool { return e.id == m.focusedSessionId })
			if index > -1 && index < len(sessions)-1 {
				m.focusedSessionId = sessions[index+1].id
			} else if index != len(sessions)-1 {
				m.focusedSessionId = sessions[0].id
			}
			return m, tea.Batch(listEntitiesCmd, previewCmd(m))
		}
	}

	return nil, nil
}

func (m Model) viewSessions() string {
	var sessions []string

	for _, session := range m.sessions {
		panesString := lipgloss.NewStyle().
			Width(m.windowWidth/3 - 4).
			MaxWidth(m.windowWidth/3 - 5).
			SetString(fmt.Sprintf("%d: %s", session.id, session.name))
		if session.id == m.focusedSessionId {
			panesString = panesString.Foreground(lipgloss.Color("2"))
		}
		sessions = append(sessions, panesString.String())
	}

	return lipgloss.JoinVertical(lipgloss.Top, sessions...)
}

func goToSession(m Model) tea.Cmd {
	return func() tea.Msg {
		c := exec.Command("tmux", "switch-client", "-t", fmt.Sprintf("$%d", m.focusedSessionId))
		c.Run()
		return tea.QuitMsg{}
	}
}

func renameSessionCmd(m Model) tea.Cmd {
	return func() tea.Msg {
		_ = goToSession(m)()
		c := exec.Command("tmux", "rename-session", "-t", m.sessionWithId(m.focusedSessionId).name, m.textInput.Value())
		c.Run()
		return tea.QuitMsg{}
	}
}

func newSessionCmd(m Model) tea.Cmd {
	return func() tea.Msg {
    c := exec.Command("tmux",
			"new-session", "-ds", m.textInput.Value(), ";",
			"switch-client", "-t", m.textInput.Value())
		c.Run()
		return tea.QuitMsg{}
	}
}

func deleteSessionCmd(m Model) tea.Cmd {
	return func() tea.Msg {
		c := exec.Command("tmux", "kill-session", "-t", m.sessionWithId(m.focusedSessionId).name)
		c.Run()
		return tickMsg{}
	}
}
