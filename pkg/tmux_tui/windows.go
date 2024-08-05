package tmux_tui

import (
	"bufio"
	"fmt"
	"math"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) initWindows() model {
	return m
}

func (m model) updateWindows(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case windowsListMsg:
		if len(m.windows) != len(msg.entities) {
			m.focusedWindowsItem = 0
			for i, v := range msg.entities {
				if v.id == msg.current {
					m.focusedWindowsItem = i
				}
			}
		}
		m.windows = msg.entities
		return m, listPanesCmd(m)
	case tea.KeyMsg:
		if m.focusedPane != 2 {
			return nil, nil
		}
		switch msg.String() {
		case "ctrl+p", "k", tea.KeyUp.String():
			m.focusedWindowsItem = int(math.Max(0, float64(m.focusedWindowsItem)-1))
			return m, tea.Batch(listPanesCmd(m), previewCmd(m))
		case "ctrl+n", "j", tea.KeyDown.String():
			m.focusedWindowsItem = int(math.Min(float64(len(m.windows)-1), float64(m.focusedWindowsItem)+1))
			return m, tea.Batch(listPanesCmd(m), previewCmd(m))
		case tea.KeyEnter.String():
			return m, goToWindow(m)
		}
	}

	return nil, nil
}

func (m model) viewWindows() string {
	output := ""

	for i, window := range m.windows {
		windowString := lipgloss.NewStyle().
			MaxWidth(m.windowWidth/3 - 5).
			SetString(fmt.Sprintf("%d: %s", window.id, window.name))
		if i == m.focusedWindowsItem {
			windowString = windowString.Foreground(lipgloss.Color("2"))
		}
		output = lipgloss.JoinVertical(lipgloss.Top, output, windowString.String())
	}

	return output
}

func listWindowsCmd(m model) tea.Cmd {
	return func() tea.Msg {
		windows := windowsListMsg{}

		c := exec.Command("tmux", "list-windows", "-aF", "#{session_id}:#{window_id}:#{window_name}")
		bytes, err := c.Output()
		if err != nil {
			return nil
		}
		str := string(bytes[:])

		scanner := bufio.NewScanner(strings.NewReader(str))
		for scanner.Scan() {
			parts := strings.Split(scanner.Text(), ":")
			session, err := strconv.Atoi(strings.Replace(parts[0], "$", "", 1))
			if err != nil {
				continue
			}
			if len(m.sessions) > 0 && session != m.sessions[m.focusedSessionsItem].id {
				continue
			}
			id, err := strconv.Atoi(strings.Replace(parts[1], "@", "", 1))
			name := parts[2]
			if err != nil {
				continue
			}
			windows.entities = append(windows.entities, entity{id, name})
		}

		sort.Slice(windows.entities, func(i, j int) bool {
			return windows.entities[i].id < windows.entities[j].id
		})

		c = exec.Command("tmux", "display-message", "-p", "#{window_id}")
		bytes, err = c.Output()
		if err != nil {
			windows.current = 0
			return windows
		}
		str = string(bytes[:])
		str = strings.TrimSpace(string(bytes[:]))
		windows.current, err = strconv.Atoi(strings.Replace(str, "@", "", 1))
		if err != nil {
			windows.current = 0
		}

		return windows
	}
}

func goToWindow(m model) tea.Cmd {
	return func() tea.Msg {
		_ = goToSession(m)()
		c := exec.Command("tmux", "select-window", "-t", fmt.Sprintf("@%d", m.windows[m.focusedWindowsItem].id))
		c.Run()
		return tea.QuitMsg{}
	}
}

func renameWindowCmd(m model) tea.Cmd {
	return func() tea.Msg {
		c := exec.Command("tmux", "rename-window", "-t", fmt.Sprintf("@%d", m.windows[m.focusedWindowsItem].id), m.textInput.Value())
		c.Run()
		return goToWindow(m)()
	}
}

func newWindowCmd(m model) tea.Cmd {
	return func() tea.Msg {
		c := exec.Command("tmux", "new-window", "-n", m.textInput.Value())
		c.Run()
		return tea.QuitMsg{}
	}
}

func deleteWindowCmd(m model) tea.Cmd {
	return func() tea.Msg {
		c := exec.Command("tmux", "kill-window", "-t", fmt.Sprintf("@%d", m.windows[m.focusedWindowsItem].id))
		c.Run()
		return tickMsg{}
	}
}
