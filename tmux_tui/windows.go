package tmux_tui

import (
	"fmt"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

func goToWindowCmd(m AppModel) tea.Cmd {
	return func() tea.Msg {
		window := m.windows.ItemWithId(m.windows.currentId)
		c := exec.Command("tmux",
			"switch-client", "-t", fmt.Sprintf("$%d", window.parent), ";",
			"select-window", "-t", fmt.Sprintf("@%d", window.id))
		c.Run()
		return tea.QuitMsg{}
	}
}

func renameWindowCmd(m AppModel) tea.Cmd {
	return func() tea.Msg {
		c := exec.Command("tmux", "rename-window", "-t", fmt.Sprintf("@%d", m.windows.currentId), m.textInput.Value())
		c.Run()
		return clearInputTextMsg{}
	}
}

func newWindowCmd(m AppModel) tea.Cmd {
	return func() tea.Msg {
		if len(m.textInput.Value()) > 0 {
			c := exec.Command("tmux", "new-window", "-n", m.textInput.Value(), "-t", fmt.Sprintf("$%d:", m.sessions.currentId))
			c.Run()
		} else {
			c := exec.Command("tmux", "new-window", "-t", fmt.Sprintf("$%d:", m.sessions.currentId))
			c.Run()
		}
		return clearInputTextMsg{}
	}
}

func deleteWindowCmd(m AppModel) tea.Cmd {
	return func() tea.Msg {
		c := exec.Command("tmux", "kill-window", "-t", fmt.Sprintf("@%d", m.windows.currentId))
		c.Run()
		return tickMsg{}
	}
}

func swapWindowsCmd(m AppModel, src int) tea.Cmd {
	return func() tea.Msg {
		src := fmt.Sprintf("@%d", src)
		dst := fmt.Sprintf("@%d", m.windows.currentId)
		c := exec.Command("tmux", "swap-window", "-s", src, "-t", dst)
		c.Run()
		return tickMsg{}
	}
}
