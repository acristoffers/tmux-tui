package tmux_tui

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func NewApplication() *tea.Program {
	model := model{
		windowHeight:        60,
		windowWidth:         80,
		focusedPane:         1,
		focusedSessionsItem: 0,
		focusedWindowsItem:  0,
		focusedPanesItem:    0,
		preview:             "",
		appState:            MainWindow,
		inputAction:         None,
		showAll:             false,
	}

	model.textInput = textinput.New()
	model.textInput.Placeholder = "New name"
	model.textInput.Focus()
	model.textInput.CharLimit = 30
	model.textInput.Width = 30

	return tea.NewProgram(model, tea.WithAltScreen())
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tickCmd(), listSessionsCmd)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.appState != TextInput {
		r, c := m.updateSessions(msg)
		if r != nil {
			return r, c
		}

		r, c = m.updateWindows(msg)
		if r != nil {
			return r, c
		}

		r, c = m.updatePanes(msg)
		if r != nil {
			return r, c
		}
	}

	switch m.appState {
	case MainWindow:
		return m.UpdateNormal(msg)
	case Swapping:
		return m.UpdateSwapping(msg)
	case TextInput:
		m, cmd := m.UpdateTextInput(msg)
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) UpdateNormal(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", tea.KeyEsc.String():
			return m, tea.Quit
		case "1":
			m.focusedPane = 1
			return m, previewCmd(m)
		case "2":
			m.focusedPane = 2
			return m, previewCmd(m)
		case "3":
			m.focusedPane = 3
			return m, previewCmd(m)
		case "r":
			if m.focusedPane == 1 {
				m.inputAction = RenameSession
				m.appState = TextInput
				return m, nil
			} else if m.focusedPane == 2 {
				m.inputAction = RenameWindow
				m.appState = TextInput
				return m, nil
			}
		case "n":
			if m.focusedPane == 1 {
				m.inputAction = NewSession
				m.appState = TextInput
				return m, nil
			} else if m.focusedPane == 2 {
				m.inputAction = NewWindow
				m.appState = TextInput
				return m, nil
			}
		case "d":
			switch m.focusedPane {
			case 1:
				return m, deleteSessionCmd(m)
			case 2:
				return m, deleteWindowCmd(m)
			case 3:
				return m, deletePaneCmd(m)
			}
		case "a":
			m.showAll = !m.showAll
			return m, listSessionsCmd
		case "s":
			switch m.focusedPane {
			case 2:
				m.appState = Swapping
				m.swapSrc = m.focusedWindowsItem
			case 3:
				m.appState = Swapping
				m.swapSrc = m.focusedPanesItem
			}
		case tea.KeyEnter.String():
			switch m.focusedPane {
			case 1:
				return m, goToSession(m)
			case 2:
				return m, goToWindow(m)
			case 3:
				return m, goToPane(m)
			}
		}
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
	case tickMsg:
		return m, tea.Batch(tickCmd(), listSessionsCmd)
	case previewMsg:
		m.preview = string(msg)
	}

	return m, nil
}

func (m model) UpdateSwapping(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case tea.KeyEsc.String():
			m.appState = MainWindow
			return m, nil
		case "a":
			m.showAll = !m.showAll
			return m, listSessionsCmd
		case "s", tea.KeySpace.String(), tea.KeyEnter.String():
			m.appState = MainWindow
			switch m.focusedPane {
			case 2:
				if m.swapSrc == m.focusedWindowsItem {
					return m, nil
				}
				return m, swapWindowsCmd(m)
			case 3:
				if m.swapSrc == m.focusedPanesItem {
					return m, nil
				}
				return m, swapPanesCmd(m)
			}
		}
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
	case tickMsg:
		return m, tea.Batch(tickCmd(), listSessionsCmd)
	case previewMsg:
		m.preview = string(msg)
	}

	return m, nil
}

func (m model) UpdateTextInput(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyEnter.String():
			m.appState = MainWindow
			if len(m.textInput.Value()) == 0 {
				m.inputAction = None
				m.textInput.SetValue("")
				return m, nil
			}
			switch m.inputAction {
			case RenameSession:
				return m, renameSessionCmd(m)
			case RenameWindow:
				return m, renameWindowCmd(m)
			case NewSession:
				return m, newSessionCmd(m)
			case NewWindow:
				return m, newWindowCmd(m)
			}
		case tea.KeyEsc.String():
			m.inputAction = None
			m.appState = MainWindow
			m.textInput.SetValue("")
		}
	}
	return m, nil
}

func (m model) View() string {
	w := m.windowWidth
	h := m.windowHeight

	switch m.appState {
	case MainWindow, Swapping:
		h60 := h * 6 / 10
		h40 := h - h60

		w33 := w / 3
		lw33 := w - 2*w33

		preview := lipgloss.NewStyle().MaxWidth(w-4).MaxHeight(h60-4).Align(lipgloss.Left, lipgloss.Top).Render(m.preview)
		previewPane := Pane("Preview", w, h60, false, preview)

		sessionsPane := Pane("[1] Sessions", w33, h40-3, m.focusedPane == 1, m.viewSessions())
		windowsPane := Pane("[2] Windows", w33, h40-3, m.focusedPane == 2, m.viewWindows())
		panesPane := Pane("[3] Panes", lw33, h40-3, m.focusedPane == 3, m.viewPanes())

		statusPane := Pane("Status", w, 3, false, statusLine(m))

		horizontalBox := lipgloss.JoinHorizontal(lipgloss.Left, sessionsPane, windowsPane, panesPane)
		verticalBox := lipgloss.JoinVertical(lipgloss.Top, previewPane, horizontalBox, statusPane)

		return verticalBox
	case TextInput:
		return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, m.textInput.View())
	}

	return ""
}

func previewCmd(m model) tea.Cmd {
	return func() tea.Msg {
		id := ""
		switch m.focusedPane {
		case 1:
			id = fmt.Sprintf("$%d", m.sessions[m.focusedSessionsItem].id)
		case 2:
			id = fmt.Sprintf("@%d", m.windows[m.focusedWindowsItem].id)
		case 3:
			id = fmt.Sprintf("%%%d", m.panes[m.focusedPanesItem].id)
		}
		c := exec.Command("tmux", "capture-pane", "-ep", "-t", id)
		bytes, err := c.Output()
		if err != nil {
			return nil
		}
		preview := string(bytes[:])
		return previewMsg(preview)
	}
}

func statusLine(m model) string {
	left := []string{"Quit: q", "Go to: <enter>", "Delete: d"}

	if m.focusedPane != 3 {
		left = append(left, "New: n")
		left = append(left, "Rename: r")
	}

	if m.showAll {
		item := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render("Show all: a")
		left = append(left, item)
	} else {
		left = append(left, "Show all: a")
	}

	if m.appState == Swapping {
		left = append(left, lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render("Swap: s"))
	} else {
		left = append(left, "Swap: s")
	}

	leftString := strings.Join(left, " | ")
	rightString := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render(strings.TrimSpace(Version))

	leftString = lipgloss.PlaceHorizontal(m.windowWidth-5-lipgloss.Width(rightString), lipgloss.Left, leftString)

	return leftString + " " + rightString
}
