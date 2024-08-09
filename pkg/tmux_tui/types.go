package tmux_tui

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
)

type entity struct {
	id   int
	name string
  parent int
}

type AppState int

const (
	MainWindow AppState = iota
	TextInput
  Swapping
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

	showAll bool

	sessions []entity
	windows  []entity
	panes    []entity

	preview string

	userInput   string
	appState    AppState
	inputAction InputAction

  swapSrc int

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

type showAllMsg struct{}

type swapStartMsg struct{}
type swapExecuteMsg struct{}
