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
		}
	}

	return nil, nil
}

func (m model) viewWindows() string {
	var windows []string

	for i, window := range m.windows {
		panesString := lipgloss.NewStyle().
			Width(m.windowWidth/3 - 4).
			MaxWidth(m.windowWidth/3 - 4).
			SetString(fmt.Sprintf("%d: %s", window.id, window.name))
		if i == m.focusedWindowsItem {
			panesString = panesString.Foreground(lipgloss.Color("2"))
		}
		if m.appState == Swapping && i == m.swapSrc && m.focusedPane == 2 {
			panesString = panesString.Background(lipgloss.Color("6")).Foreground(lipgloss.Color("0"))
		}
		windows = append(windows, panesString.String())
	}

	return lipgloss.JoinVertical(lipgloss.Top, windows...)
}

func listWindowsCmd(m model) tea.Cmd {
	return func() tea.Msg {
		windows := windowsListMsg{}

		c := exec.Command("tmux", "list-windows", "-aF", "#{session_id}\t#{window_id}\t#{window_name}")
		bytes, err := c.Output()
		if err != nil {
			return nil
		}
		str := string(bytes[:])

		scanner := bufio.NewScanner(strings.NewReader(str))
		for scanner.Scan() {
			parts := strings.Split(scanner.Text(), "\t")
			session, err := strconv.Atoi(strings.Replace(parts[0], "$", "", 1))
			if err != nil {
				continue
			}
			if !m.showAll && len(m.sessions) > 0 && session != m.sessions[m.focusedSessionsItem].id {
				continue
			}
			id, err := strconv.Atoi(strings.Replace(parts[1], "@", "", 1))
			name := parts[2]
			if err != nil {
				continue
			}
			windows.entities = append(windows.entities, entity{id, name, session})
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
		id := m.windows[m.focusedWindowsItem].id
		c := exec.Command("tmux", "list-windows", "-aF", "#{session_name}", "-f", fmt.Sprintf("#{m:#{window_id},@%d}", id))
		bytes, err := c.Output()
		if err != nil {
			return nil
		}
		str := strings.TrimSpace(string(bytes[:]))
		c = exec.Command("tmux", "switch-client", "-t", str)
		c.Run()
		c = exec.Command("tmux", "select-window", "-t", fmt.Sprintf("@%d", id))
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
		c := exec.Command("tmux", "new-window", "-n", m.textInput.Value(), "-t", fmt.Sprintf("%s:", m.sessions[m.focusedSessionsItem].name))
		c.Run()
		return goToSession(m)()
	}
}

func deleteWindowCmd(m model) tea.Cmd {
	return func() tea.Msg {
		c := exec.Command("tmux", "kill-window", "-t", fmt.Sprintf("@%d", m.windows[m.focusedWindowsItem].id))
		c.Run()
		return tickMsg{}
	}
}

func swapWindowsCmd(m model) tea.Cmd {
	return func() tea.Msg {
		src := fmt.Sprintf("@%d", m.windows[m.swapSrc].id)
		dst := fmt.Sprintf("@%d", m.windows[m.focusedWindowsItem].id)
		c := exec.Command("tmux", "swap-window", "-s", src, "-t", dst)
		c.Run()
		return tickMsg{}
	}
}
