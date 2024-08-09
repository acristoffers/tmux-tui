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

func (m model) initPanes() model {
	return m
}

func (m model) updatePanes(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case panesListMsg:
		if len(m.panes) != len(msg.entities) {
			m.focusedPanesItem = 0
			for i, v := range msg.entities {
				if v.id == msg.current {
					m.focusedPanesItem = i
				}
			}
		}
		m.panes = msg.entities
		return m, nil
	case tea.KeyMsg:
		if m.focusedPane != 3 {
			return nil, nil
		}
		switch msg.String() {
		case "ctrl+p", "k", tea.KeyUp.String():
			m.focusedPanesItem = int(math.Max(0, float64(m.focusedPanesItem)-1))
			return m, previewCmd(m)
		case "ctrl+n", "j", tea.KeyDown.String():
			m.focusedPanesItem = int(math.Min(float64(len(m.panes)-1), float64(m.focusedPanesItem)+1))
			return m, previewCmd(m)
		case tea.KeyEnter.String():
			return m, goToPane(m)
		}
	}

	return nil, nil
}

func (m model) viewPanes() string {
	output := ""

	for i, pane := range m.panes {
		panesString := lipgloss.NewStyle().
			MaxWidth(m.windowWidth/3 - 5).
			SetString(fmt.Sprintf("%d", pane.id))
		if i == m.focusedPanesItem {
			panesString = panesString.Foreground(lipgloss.Color("2"))
		}
		output = lipgloss.JoinVertical(lipgloss.Top, output, panesString.String())
	}

	return output
}

func listPanesCmd(m model) tea.Cmd {
	return func() tea.Msg {
		panes := panesListMsg{}

		c := exec.Command("tmux", "list-panes", "-aF", "#{window_id}:#{pane_id}")
		bytes, err := c.Output()
		if err != nil {
			return nil
		}
		str := string(bytes[:])

		scanner := bufio.NewScanner(strings.NewReader(str))
		for scanner.Scan() {
			parts := strings.Split(scanner.Text(), ":")
			window, err := strconv.Atoi(strings.Replace(parts[0], "@", "", 1))
			if err != nil {
				continue
			}
			if !m.showAll && len(m.windows) > 0 && window != m.windows[m.focusedWindowsItem].id {
				continue
			}
			id, err := strconv.Atoi(strings.Replace(parts[1], "%", "", 1))
			if err != nil {
				continue
			}
			panes.entities = append(panes.entities, entity{id, ""})
		}

		sort.Slice(panes.entities, func(i, j int) bool {
			return panes.entities[i].id < panes.entities[j].id
		})

		c = exec.Command("tmux", "display-message", "-p", "#{pane_id}")
		bytes, err = c.Output()
		if err != nil {
			panes.current = 0
			return panes
		}
		str = string(bytes[:])
		str = strings.TrimSpace(string(bytes[:]))
		panes.current, err = strconv.Atoi(strings.Replace(str, "%", "", 1))
		if err != nil {
			panes.current = 0
		}

		return panes
	}
}

func goToPane(m model) tea.Cmd {
	return func() tea.Msg {
		_ = goToWindow(m)()
		c := exec.Command("tmux", "select-pane", "-t", fmt.Sprintf("%%%d", m.panes[m.focusedPanesItem].id))
		c.Run()
		return tea.QuitMsg{}
	}
}

func deletePaneCmd(m model) tea.Cmd {
	return func() tea.Msg {
		c := exec.Command("tmux", "kill-pane", "-t", fmt.Sprintf("%%%d", m.panes[m.focusedPanesItem].id))
		c.Run()
		return tickMsg{}
	}
}
