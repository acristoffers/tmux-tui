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

type (
	errorMsg          string
	previewMsg        string
	tickMsg           time.Time
	clearInputTextMsg struct{}

	listEntitiesMsg struct {
		sessions       []TmuxEntity
		windows        []TmuxEntity
		panes          []TmuxEntity
		currentSession int
		currentWindow  int
		currentPane    int
	}
)

const (
	None InputAction = iota
	Filter
	NewSession
	NewWindow
	RenameSession
	RenameWindow
)

type (
	terminal struct {
		width  int
		height int
	}

	TmuxEntity struct {
		id     int
		name   string
		parent int
	}

	InputAction int

	AppModel struct {
		Error string

		terminal terminal
		theme    Theme

		preview  Frame
		sessions ListFrame
		windows  ListFrame
		panes    ListFrame

		focusedFrame int

		showAll bool
		swapSrc int

		textInput   textinput.Model
		inputAction InputAction
		filter      string
	}
)

func NewApplication(theme Theme) *tea.Program {
	model := AppModel{
		Error:        "",
		terminal:     terminal{80, 80},
		theme:        theme,
		preview:      Frame{title: "Preview"},
		sessions:     ListFrame{frame: Frame{title: "[1] Sessions", focused: true}, parentId: -1},
		windows:      ListFrame{frame: Frame{title: "[2] Windows"}, parentId: -1},
		panes:        ListFrame{frame: Frame{title: "[3] Panes"}, parentId: -1},
		focusedFrame: 1,
		showAll:      false,
		swapSrc:      -1,
		inputAction:  None,
	}

	model.textInput = textinput.New()
	model.textInput.Focus()
	model.textInput.TextStyle = lipgloss.NewStyle().Foreground(theme.Foreground).Background(theme.Background)

	return tea.NewProgram(model, tea.WithAltScreen())
}

func (m AppModel) Init() tea.Cmd {
	return tea.Batch(tickCmd(), listEntitiesCmd)
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd = nil

	if m.swapSrc != -1 {
		goto swap_mode
	}

	if m.inputAction != None {
		goto input_mode
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyEsc.String():
			cmd = tea.Quit
		case "1":
			m.focusedFrame = 1
			m.sessions.frame.focused = true
			m.windows.frame.focused = false
			m.panes.frame.focused = false
			cmd = previewCmd(m)
		case "2":
			m.focusedFrame = 2
			m.sessions.frame.focused = false
			m.windows.frame.focused = true
			m.panes.frame.focused = false
			cmd = previewCmd(m)
		case "3":
			m.focusedFrame = 3
			m.sessions.frame.focused = false
			m.windows.frame.focused = false
			m.panes.frame.focused = true
			cmd = previewCmd(m)
		case tea.KeyEnter.String():
			switch m.focusedFrame {
			case 1:
				cmd = goToSessionCmd(m)
			case 2:
				cmd = goToWindowCmd(m)
			case 3:
				cmd = goToPaneCmd(m)
			}
		case "d":
			switch m.focusedFrame {
			case 1:
				cmd = deleteSessionCmd(m)
			case 2:
				cmd = deleteWindowCmd(m)
			case 3:
				cmd = deletePaneCmd(m)
			}
		case "h":
			if m.focusedFrame == 3 {
				cmd = splitPane(m, true)
			}
		case "v":
			if m.focusedFrame == 3 {
				cmd = splitPane(m, false)
			}
		case "r":
			switch m.focusedFrame {
			case 1:
				m.inputAction = RenameSession
				m.textInput.SetValue(m.sessions.ItemWithId(m.sessions.currentId).name)
				m.textInput.SetCursor(100)
			case 2:
				m.inputAction = RenameWindow
				m.textInput.SetValue(m.windows.ItemWithId(m.windows.currentId).name)
				m.textInput.SetCursor(100)
			}
		case "n":
			m.textInput.SetValue("")
			switch m.focusedFrame {
			case 1:
				m.inputAction = NewSession
			case 2:
				m.inputAction = NewWindow
			}
		case "N":
			switch m.focusedFrame {
			case 1:
				cmd = newSessionCmd(m)
			case 2:
				cmd = newWindowCmd(m)
			}
		case "/":
			m.inputAction = Filter
			m.textInput.SetValue(m.filter)
			m.textInput.SetCursor(100)
		case "s":
			switch m.focusedFrame {
			case 2:
				m.swapSrc = m.windows.currentId
				m.windows.MarkSelection()
			case 3:
				m.swapSrc = m.panes.currentId
				m.panes.MarkSelection()
			}
		}
	case clearInputTextMsg:
		m.textInput.SetValue("")
		cmd = listEntitiesCmd
	}
	goto common_bindings

swap_mode:
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyEsc.String():
			m.swapSrc = -1
			m.windows.ClearMarks()
			m.panes.ClearMarks()
		case "s", tea.KeySpace.String(), tea.KeyEnter.String():
			m.windows.ClearMarks()
			m.panes.ClearMarks()
			switch m.focusedFrame {
			case 2:
				cmd = swapWindowsCmd(m, m.swapSrc)
			case 3:
				cmd = swapPanesCmd(m, m.swapSrc)
			}
			m.swapSrc = -1
		}
	}
	goto common_bindings

input_mode:
	m.textInput, cmd = m.textInput.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyEsc.String():
			if m.inputAction == Filter {
				m.filter = ""
				m.sessions.filterText = ""
				m.windows.filterText = ""
				m.panes.filterText = ""
			}
			m.inputAction = None
			m.textInput.SetValue("")
		case tea.KeyEnter.String():
			switch m.inputAction {
			case NewSession:
				cmd = newSessionCmd(m)
			case RenameSession:
				cmd = renameSessionCmd(m)
			case NewWindow:
				cmd = newWindowCmd(m)
			case RenameWindow:
				cmd = renameWindowCmd(m)
			case Filter:
				m.filter = m.textInput.Value()
			}
			m.inputAction = None
		}
	}

	if m.inputAction == Filter {
		m.sessions.filterText = m.textInput.Value()
		m.windows.filterText = m.textInput.Value()
		m.panes.filterText = m.textInput.Value()
	}

	goto basic_handlers

common_bindings:
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			cmd = tea.Quit
		case "ctrl+p", "k", tea.KeyUp.String():
			switch m.focusedFrame {
			case 1:
				m.sessions.SelectPrevious()
			case 2:
				m.windows.SelectPrevious()
			case 3:
				m.panes.SelectPrevious()
			}
			cmd = listEntitiesCmd
		case "ctrl+n", "j", tea.KeyDown.String():
			switch m.focusedFrame {
			case 1:
				m.sessions.SelectNext()
			case 2:
				m.windows.SelectNext()
			case 3:
				m.panes.SelectNext()
			}
			cmd = listEntitiesCmd
		case "a":
			m.showAll = !m.showAll
			cmd = listEntitiesCmd
		}
	}

basic_handlers:
	switch msg := msg.(type) {
	case tickMsg:
		cmd = tea.Batch(tickCmd(), listEntitiesCmd)
	case tea.WindowSizeMsg:
		m.terminal.width = msg.Width
		m.terminal.height = msg.Height
	case listEntitiesMsg:
		m.sessions.items = msg.sessions
		m.windows.items = msg.windows
		m.panes.items = msg.panes
		if m.sessions.currentId == -1 && len(m.filter) == 0 {
			m.sessions.currentId = msg.currentSession
			m.windows.currentId = msg.currentWindow
			m.panes.currentId = msg.currentPane
		}
		cmd = previewCmd(m)
	case previewMsg:
		m.preview.contents = string(msg)
	}

	if m.showAll {
		m.windows.parentId = -1
		m.panes.parentId = -1
	} else {
		m.windows.parentId = m.sessions.currentId
		m.panes.parentId = m.windows.currentId
	}

	m.sessions.Update()
	m.windows.Update()
	m.panes.Update()

	return m, cmd
}

func (m AppModel) View() string {
	preview := m.preview

	sessions := m.sessions.RenderContents(m.theme)
	windows := m.windows.RenderContents(m.theme)
	panes := m.panes.RenderContents(m.theme)

	if m.swapSrc == -1 {
		sessions.title = "[1] Sessions"
		windows.title = "[2] Windows"
		panes.title = "[3] Panes"
	} else {
		sessions.title = "Sessions"
		windows.title = "Windows"
		panes.title = "Panes"
	}

	var status Frame
	if m.inputAction == None {
		status = m.StatusBar()
	} else if m.inputAction == Filter {
		m.textInput.Width = m.terminal.width - 4
		status = Frame{
			title:    "Filter",
			contents: m.textInput.View(),
			width:    m.terminal.width,
			height:   m.terminal.height,
			focused:  true,
		}
	} else {
		m.textInput.Width = m.terminal.width - 4
		status = Frame{
			title:    "New name",
			contents: m.textInput.View(),
			width:    m.terminal.width,
			height:   m.terminal.height,
			focused:  true,
		}
	}

	return m.DrawGrid(preview, sessions, windows, panes, status)
}

func (m AppModel) StatusBar() Frame {
	frame := Frame{title: "Status"}

	style := lipgloss.NewStyle().Foreground(m.theme.Foreground).Background(m.theme.Background)

	left := []string{"Quit: q"}

	if m.swapSrc == -1 {
		left = append(left, style.Render("Go to: <enter>"))
		left = append(left, style.Render("Delete: d"))
		left = append(left, style.Render("Swap: s"))

		if m.focusedFrame != 3 {
			left = append(left, style.Render("New: n"))
			left = append(left, style.Render("New (nameless): N"))
			left = append(left, style.Render("Rename: r"))
		} else {
			left = append(left, style.Render("Vertical split: v"))
			left = append(left, style.Render("Horizontal split: h"))
		}
	} else {
		left = append(left, style.Foreground(m.theme.Accent).Render("Swap: s/<space>/<enter>"))
		left = append(left, style.Render("Cancel: <esc>"))
	}

	if m.showAll {
		left = append(left, style.Foreground(m.theme.Accent).Render("Show all: a"))
	} else {
		left = append(left, style.Render("Show all: a"))
	}

	if len(m.textInput.Value()) > 0 {
		left = append(left, style.Foreground(m.theme.Accent).Render("Filter: /"))
	} else {
		left = append(left, style.Render("Filter: /"))
	}

	rightString := style.Foreground(m.theme.Secondary).Render(strings.TrimSpace(Version))

	separator := style.Render(" | ")
	maxWidth := uint(m.terminal.width - 7 - lipgloss.Width(rightString))
	leftString := left[0]
	for i, v := range left {
		if i == 0 {
			continue
		}
		newWidth := lipgloss.Width(leftString) + 3 + lipgloss.Width(v)
		if newWidth <= int(maxWidth) {
			leftString = fmt.Sprintf("%s%s%s", leftString, separator, v)
		}
	}
	leftString = style.Width(int(maxWidth)).Render(leftString)

	space := style.Render(" ")
	frame.contents = leftString + space + rightString

	return frame
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func listEntitiesCmd() tea.Msg {
	// Fetches info about all sessions, windows and panes at once
	c := exec.Command("tmux",
		"list-panes", "-aF", "#{session_id}\t#{window_id}\t#{pane_id}\t#{session_name}\t#{window_name}\t#{pane_current_command}", ";",
		"display-message", "-p", "#{session_id}\t#{window_id}\t#{pane_id}")
	bytes, err := c.Output()
	if err != nil {
		return nil
	}

	sessions := []TmuxEntity{}
	windows := []TmuxEntity{}
	panes := []TmuxEntity{}

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
		pane_name := parts[5]

		sessions = append(sessions, TmuxEntity{session_id, session_name, -1})
		windows = append(windows, TmuxEntity{window_id, window_name, session_id})
		panes = append(panes, TmuxEntity{pane_id, pane_name, window_id})
	}

	if len(sessions) == 0 {
		return errorMsg("No sessions found. Is tmux running?")
	}

	eq := func(a, b TmuxEntity) bool {
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

func previewCmd(m AppModel) tea.Cmd {
	return func() tea.Msg {
		id := ""
		switch m.focusedFrame {
		case 1:
			id = fmt.Sprintf("$%d", m.sessions.currentId)
		case 2:
			id = fmt.Sprintf("@%d", m.windows.currentId)
		case 3:
			id = fmt.Sprintf("%%%d", m.panes.currentId)
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
