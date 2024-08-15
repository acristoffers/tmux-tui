package tmux_tui

import (
	"fmt"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

func goToSessionCmd(m AppModel) tea.Cmd {
	return func() tea.Msg {
		c := exec.Command("tmux", "switch-client", "-t", fmt.Sprintf("$%d", m.sessions.currentId))
		c.Run()
		return tea.QuitMsg{}
	}
}

func renameSessionCmd(m AppModel) tea.Cmd {
	return func() tea.Msg {
		_ = goToSessionCmd(m)()
		session := fmt.Sprintf("$%d", m.sessions.currentId)
		c := exec.Command("tmux", "rename-session", "-t", session, m.textInput.Value())
		c.Run()
		return clearInputTextMsg{}
	}
}

func newSessionCmd(m AppModel) tea.Cmd {
	return func() tea.Msg {
		if len(m.textInput.Value()) > 0 {
			c := exec.Command("tmux",
				"new-session", "-ds", m.textInput.Value(), ";",
				"switch-client", "-t", m.textInput.Value())
			c.Run()
		} else {
			c := exec.Command("tmux", "new-session", "-d")
			c.Run()
		}
		return clearInputTextMsg{}
	}
}

func deleteSessionCmd(m AppModel) tea.Cmd {
	return func() tea.Msg {
		c := exec.Command("tmux", "kill-session", "-t", m.sessions.ItemWithId(m.sessions.currentId).name)
		c.Run()
		return tickMsg{}
	}
}
