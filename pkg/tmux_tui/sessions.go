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

func (m model) initSessions() model {
	return m
}

func (m model) updateSessions(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case sessionsListMsg:
		if len(m.sessions) != len(msg.entities) {
			m.focusedSessionsItem = 0
			for i, v := range msg.entities {
				if v.id == msg.current {
					m.focusedSessionsItem = i
				}
			}
		}
		m.sessions = msg.entities
		if len(m.preview) == 0 {
			return m, tea.Batch(listWindowsCmd(m), previewCmd(m))
		}
		return m, listWindowsCmd(m)
	case tea.KeyMsg:
		if m.focusedPane != 1 {
			return nil, nil
		}
		switch msg.String() {
		case "ctrl+p", "k", tea.KeyUp.String():
			m.focusedSessionsItem = int(math.Max(0, float64(m.focusedSessionsItem)-1))
			return m, tea.Batch(listWindowsCmd(m), previewCmd(m))
		case "ctrl+n", "j", tea.KeyDown.String():
			m.focusedSessionsItem = int(math.Min(float64(len(m.sessions)-1), float64(m.focusedSessionsItem)+1))
			return m, tea.Batch(listWindowsCmd(m), previewCmd(m))
		case tea.KeyEnter.String():
			return m, goToSession(m)
		}
	}

	return nil, nil
}

func (m model) viewSessions() string {
	output := ""

	for i, session := range m.sessions {
		sessionString := lipgloss.NewStyle().
			MaxWidth(m.windowWidth/3 - 5).
			SetString(fmt.Sprintf("%d: %s", session.id, session.name))
		if i == m.focusedSessionsItem {
			sessionString = sessionString.Foreground(lipgloss.Color("2"))
		}
		output = lipgloss.JoinVertical(lipgloss.Top, output, sessionString.String())
	}

	return output
}

func listSessionsCmd() tea.Msg {
	sessions := sessionsListMsg{}

	c := exec.Command("tmux", "list-sessions", "-F", "#{session_id}:#{session_name}")
	bytes, err := c.Output()
	if err != nil {
		return nil
	}
	str := string(bytes[:])

	scanner := bufio.NewScanner(strings.NewReader(str))
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ":")
		id, err := strconv.Atoi(strings.Replace(parts[0], "$", "", 1))
		name := parts[1]
		if err != nil {
			continue
		}
		sessions.entities = append(sessions.entities, entity{id, name})
	}

	sort.Slice(sessions.entities, func(i, j int) bool {
		return sessions.entities[i].id < sessions.entities[j].id
	})

	c = exec.Command("tmux", "display-message", "-p", "#{session_id}")
	bytes, err = c.Output()
	if err != nil {
		sessions.current = 0
		return sessions
	}
	str = strings.TrimSpace(string(bytes[:]))
	sessions.current, err = strconv.Atoi(strings.Replace(str, "$", "", 1))
	if err != nil {
		sessions.current = 0
	}

	return sessions
}

func goToSession(m model) tea.Cmd {
	return func() tea.Msg {
		c := exec.Command("tmux", "switch-client", "-t", m.sessions[m.focusedSessionsItem].name)
		c.Run()
		return tea.QuitMsg{}
	}
}

func renameSessionCmd(m model) tea.Cmd {
	return func() tea.Msg {
		_ = goToSession(m)()
		c := exec.Command("tmux", "rename-session", "-t", m.sessions[m.focusedSessionsItem].name, m.textInput.Value())
		c.Run()
		return tea.QuitMsg{}
	}
}

func newSessionCmd(m model) tea.Cmd {
	return func() tea.Msg {
		c := exec.Command("tmux", "new-session", "-ds", m.textInput.Value())
		c.Run()
		c = exec.Command("tmux", "switch-client", "-t", m.textInput.Value())
		c.Run()
		return tea.QuitMsg{}
	}
}

func deleteSessionCmd(m model) tea.Cmd {
	return func() tea.Msg {
		c := exec.Command("tmux", "kill-session", "-t", m.sessions[m.focusedSessionsItem].name)
		c.Run()
		return tickMsg{}
	}
}
