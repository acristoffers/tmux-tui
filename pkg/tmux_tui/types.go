package tmux_tui

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
)

type entity struct {
	id   int
	name string
}

type AppState int

const (
	MainWindow AppState = iota
	TextInput
)

type InputAction int

const (
	None InputAction = iota
	RenameSession
	RenameWindow
	NewSession
	NewWindow
)

type model struct {
	windowWidth  int
	windowHeight int
	focusedPane  int

	focusedSessionsItem int
	focusedWindowsItem  int
	focusedPanesItem    int

	sessions []entity
	windows  []entity
	panes    []entity

	preview string

	userInput   string
	appState    AppState
	inputAction InputAction

	textInput textinput.Model
}

type tickMsg time.Time
type previewMsg string

type panesListMsg struct {
	entities []entity
	current  int
}

type sessionsListMsg struct {
	entities []entity
	current  int
}

type windowsListMsg struct {
	entities []entity
	current  int
}
