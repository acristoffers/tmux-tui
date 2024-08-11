package tmux_tui

import (
	"bufio"
	"fmt"
	"os/exec"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func NewApplication() *tea.Program {
	model := Model{
		windowHeight:     60,
		windowWidth:      80,
		focusedPane:      1,
		focusedSessionId: -1,
		focusedWindowId:  -1,
		focusedPaneId:    -1,
		preview:          "",
		appState:         MainWindow,
		inputAction:      None,
		showAll:          false,
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

func (m Model) Init() tea.Cmd {
	return tea.Batch(tickCmd(), listEntitiesCmd)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errorMsg:
		m.Error = string(msg)
		return m, tea.Quit
	}

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
		return m.UpdateTextInput(msg)
	}

	return m, nil
}

func (m Model) UpdateNormal(msg tea.Msg) (Model, tea.Cmd) {
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
		case "N":
			if m.focusedPane == 1 {
				return m, newSessionCmd(m)
			} else if m.focusedPane == 2 {
				return m, newWindowCmd(m)
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
			return m, listEntitiesCmd
		case "s":
			switch m.focusedPane {
			case 2:
				m.appState = Swapping
				m.swapSrc = m.focusedWindowId
			case 3:
				m.appState = Swapping
				m.swapSrc = m.focusedPaneId
			}
		case "h":
			if m.focusedPane == 3 {
				return m, splitPane(m, true)
			}
		case "v":
			if m.focusedPane == 3 {
				return m, splitPane(m, false)
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
		return m, tea.Batch(tickCmd(), listEntitiesCmd)
	case previewMsg:
		m.preview = string(msg)
	case listEntitiesMsg:
		m.sessions = msg.sessions
		m.windows = msg.windows
		m.panes = msg.panes

		if m.focusedSessionId == -1 {
			m.focusedSessionId = msg.currentSession
			m.focusedWindowId = msg.currentWindow
			m.focusedPaneId = msg.currentPane
			return m, previewCmd(m)
		}

		if !slices.ContainsFunc(m.sessions, func(e entity) bool { return e.id == m.focusedSessionId }) {
			if slices.ContainsFunc(m.sessions, func(e entity) bool { return e.id == msg.currentSession }) {
				m.focusedSessionId = msg.currentSession
			} else {
				m.focusedSessionId = m.sessions[0].id
			}
		}

		windows := m.visibleWindows()
		if !slices.ContainsFunc(windows, func(e entity) bool { return e.id == m.focusedWindowId }) {
			if slices.ContainsFunc(windows, func(e entity) bool { return e.id == msg.currentWindow }) {
				m.focusedWindowId = msg.currentWindow
			} else {
				m.focusedWindowId = windows[0].id
			}
		}

		panes := m.visiblePanes()
		if !slices.ContainsFunc(panes, func(e entity) bool { return e.id == m.focusedPaneId }) {
			if slices.ContainsFunc(panes, func(e entity) bool { return e.id == msg.currentPane }) {
				m.focusedPaneId = msg.currentPane
			} else {
				m.focusedPaneId = panes[0].id
			}
		}

		return m, previewCmd(m)
	case clearInputTextMsg:
		m.textInput.SetValue("")
		return m, listEntitiesCmd
	}

	return m, nil
}

func (m Model) UpdateSwapping(msg tea.Msg) (Model, tea.Cmd) {
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
			return m, nil
		case "s", tea.KeySpace.String(), tea.KeyEnter.String():
			m.appState = MainWindow
			switch m.focusedPane {
			case 2:
				if m.swapSrc == m.focusedWindowId {
					return m, nil
				}
				return m, swapWindowsCmd(m)
			case 3:
				if m.swapSrc == m.focusedPaneId {
					return m, nil
				}
				return m, swapPanesCmd(m)
			}
		}
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
	case tickMsg:
		return m, tea.Batch(tickCmd(), listEntitiesCmd)
	case previewMsg:
		m.preview = string(msg)
	}

	return m, nil
}

func (m Model) UpdateTextInput(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyEnter.String():
			m.appState = MainWindow
			if len(m.textInput.Value()) == 0 {
				m.inputAction = None
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

	return m, cmd
}

func (m Model) View() string {
	w := m.windowWidth
	h := m.windowHeight

	switch m.appState {
	case MainWindow, Swapping:
		h60 := h * 6 / 10
		h40 := h - h60

		w33 := w / 3
		lw33 := w - 2*w33

		foregroundColor := lipgloss.Color("15")
		if !lipgloss.HasDarkBackground() {
			foregroundColor = lipgloss.Color("0")
		}

		preview := lipgloss.NewStyle().
			Foreground(foregroundColor).
			Background(lipgloss.NoColor{}).
			MaxWidth(w-4).
			MaxHeight(h60-2).
			Align(lipgloss.Left, lipgloss.Top).
			Render(m.preview)
		previewPane := Pane("Preview", w, h60, false, preview)

		sessionsPane := Pane("[1] Sessions", w33, h40-3, m.focusedPane == 1, m.viewSessions())
		windowsPane := Pane("[2] Windows", w33, h40-3, m.focusedPane == 2, m.viewWindows())
		panesPane := Pane("[3] Panes", lw33, h40-3, m.focusedPane == 3, m.viewPanes())

		if m.appState == Swapping {
			sessionsPane = Pane("Sessions", w33, h40-3, m.focusedPane == 1, m.viewSessions())
			windowsPane = Pane("Windows", w33, h40-3, m.focusedPane == 2, m.viewWindows())
			panesPane = Pane("Panes", lw33, h40-3, m.focusedPane == 3, m.viewPanes())
		}

		statusPane := Pane("Status", w, 3, false, statusLine(m))

		horizontalBox := lipgloss.JoinHorizontal(lipgloss.Left, sessionsPane, windowsPane, panesPane)
		verticalBox := lipgloss.JoinVertical(lipgloss.Top, previewPane, horizontalBox, statusPane)

		return verticalBox
	case TextInput:
		return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, m.textInput.View())
	}

	return ""
}

func (m Model) visibleWindows() []entity {
	if m.showAll {
		return m.windows
	}

	es := []entity{}
	for _, v := range m.windows {
		if v.parent == m.focusedSessionId {
			es = append(es, v)
		}
	}
	return es
}

func (m Model) visiblePanes() []entity {
	if m.showAll {
		return m.panes
	}

	es := []entity{}
	for _, v := range m.panes {
		if v.parent == m.focusedWindowId {
			es = append(es, v)
		}
	}
	return es
}

func (m Model) sessionWithId(id int) entity {
	e := entity{}
	for _, v := range m.sessions {
		if v.id == id {
			e = v
			break
		}
	}
	return e
}

func (m Model) windowWithId(id int) entity {
	e := entity{}
	for _, v := range m.windows {
		if v.id == id {
			e = v
			break
		}
	}
	return e
}

func (m Model) paneWithId(id int) entity {
	e := entity{}
	for _, v := range m.panes {
		if v.id == id {
			e = v
			break
		}
	}
	return e
}

func previewCmd(m Model) tea.Cmd {
	return func() tea.Msg {
		id := ""
		switch m.focusedPane {
		case 1:
			id = fmt.Sprintf("$%d", m.focusedSessionId)
		case 2:
			id = fmt.Sprintf("@%d", m.focusedWindowId)
		case 3:
			id = fmt.Sprintf("%%%d", m.focusedPaneId)
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

func statusLine(m Model) string {
	left := []string{"Quit: q"}

	if m.appState != Swapping {
    left = append(left, "Go to: <enter>")
    left = append(left, "Delete: d")

		if m.focusedPane != 3 {
			left = append(left, "New: n")
			left = append(left, "New (nameless): N")
			left = append(left, "Rename: r")
		} else {
			left = append(left, "Vertical split: v")
			left = append(left, "Horizontal split: h")
		}
	}

	if m.showAll {
		item := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render("Show all: a")
		left = append(left, item)
	} else {
		left = append(left, "Show all: a")
	}

	if m.appState == Swapping {
		left = append(left, lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render("Swap: s/<space>/<enter>"))
		left = append(left, "Cancel: <esc>")
	} else if m.focusedPane != 1 {
		left = append(left, "Swap: s")
	}

	leftString := strings.Join(left, " | ")
	rightString := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render(strings.TrimSpace(Version))

	leftString = lipgloss.PlaceHorizontal(m.windowWidth-5-lipgloss.Width(rightString), lipgloss.Left, leftString)

	return leftString + " " + rightString
}

func listEntitiesCmd() tea.Msg {
	// Fetches info about all sessions, windows and panes at once
	c := exec.Command("tmux",
		"list-panes", "-aF", "#{session_id}\t#{window_id}\t#{pane_id}\t#{session_name}\t#{window_name}", ";",
		"display-message", "-p", "#{session_id}\t#{window_id}\t#{pane_id}")
	bytes, err := c.Output()
	if err != nil {
		return nil
	}

	sessions := []entity{}
	windows := []entity{}
	panes := []entity{}

	currentSession := 0
	currentWindow := 0
	currentPane := 0

	scanner := bufio.NewScanner(strings.NewReader(string(bytes[:])))
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "\t")

		session_id, err := strconv.Atoi(strings.Replace(parts[0], "$", "", 1))
		if err != nil {
			continue
		}

		window_id, err := strconv.Atoi(strings.Replace(parts[1], "@", "", 1))
		if err != nil {
			continue
		}

		pane_id, err := strconv.Atoi(strings.Replace(parts[2], "%", "", 1))
		if err != nil {
			continue
		}

		if len(parts) == 3 {
			currentSession = session_id
			currentWindow = window_id
			currentPane = pane_id
			continue
		}

		session_name := parts[3]
		window_name := parts[4]

		sessions = append(sessions, entity{session_id, session_name, -1})
		windows = append(windows, entity{window_id, window_name, session_id})
		panes = append(panes, entity{pane_id, "", window_id})
	}

	cmp := func(a, b entity) int {
		if a.id < b.id {
			return -1
		} else if a.id == b.id {
			return 0
		} else {
			return 1
		}
	}

	if len(sessions) == 0 {
		return errorMsg("No sessions found. Is tmux running?")
	}

	slices.SortFunc(sessions, cmp)
	slices.SortFunc(windows, cmp)
	slices.SortFunc(panes, cmp)

	eq := func(a, b entity) bool {
		return a.id == b.id
	}

	sessions = slices.CompactFunc(sessions, eq)
	windows = slices.CompactFunc(windows, eq)
	panes = slices.CompactFunc(panes, eq)

	msg := listEntitiesMsg{}

	msg.sessions = sessions
	msg.windows = windows
	msg.panes = panes

	msg.currentSession = currentSession
	msg.currentWindow = currentWindow
	msg.currentPane = currentPane

	return msg
}
