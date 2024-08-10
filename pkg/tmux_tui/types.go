package tmux_tui

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
)

type entity struct {
	id     int
	name   string
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

type Model struct {
	windowWidth  int
	windowHeight int
	focusedPane  int

	focusedSessionId int
	focusedWindowId  int
	focusedPaneId    int

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

  Error string
}

type tickMsg time.Time
type previewMsg string

type listEntitiesMsg struct {
	sessions       []entity
	windows        []entity
	panes          []entity
	currentSession int
	currentWindow  int
	currentPane    int
}

type showAllMsg struct{}

type swapStartMsg struct{}
type swapExecuteMsg struct{}

type errorMsg string
