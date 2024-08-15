package tmux_tui

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/muesli/reflow/truncate"
)

type Frame struct {
	title    string
	contents string
	width    int
	height   int
	focused  bool
}

func NewFrame(m AppModel) Frame {
	return Frame{
		title:    "",
		contents: "",
		width:    m.terminal.width,
		height:   m.terminal.height,
		focused:  true,
	}
}

func (frame Frame) View(theme Theme) string {
	roundedBorder := lipgloss.RoundedBorder()

	// Account for paddings and borders
	width := frame.width - 2
	height := frame.height - 2

	// First line of the border, with the title
	truncated := truncate.String(fmt.Sprintf(" %s ", frame.title), uint(width-1))
	title := lipgloss.PlaceHorizontal(width-1, lipgloss.Left, truncated, lipgloss.WithWhitespaceChars(roundedBorder.Top))
	borderTop := fmt.Sprintf("%s%s%s%s", roundedBorder.TopLeft, roundedBorder.Top, title, roundedBorder.TopRight)

	contents := lipgloss.NewStyle().
		Foreground(theme.foreground).
		Background(theme.background).
		MaxWidth(width-4).
		MaxHeight(height).
		Align(lipgloss.Left, lipgloss.Top).
		Render(frame.contents)

	header := lipgloss.NewStyle().
		Background(theme.background).
		Foreground(theme.foreground).
		SetString(borderTop)
	pane := lipgloss.NewStyle().
		Background(theme.background).
		Border(lipgloss.RoundedBorder(), false, true, true, true).
		Foreground(theme.foreground).
		Height(height).
		PaddingLeft(1).
		PaddingRight(1).
		Width(width).
		SetString(contents)

	if frame.focused {
		header = header.Foreground(theme.accent)
		pane = pane.BorderForeground(theme.accent)
	}

	return lipgloss.JoinVertical(lipgloss.Top, header.String(), pane.String())
}

func (m AppModel) DrawGrid(preview, sessions, windows, frames, status Frame) string {
	w := m.terminal.width
	h := m.terminal.height
	h60 := h * 6 / 10
	h40 := h - h60
	w33 := w / 3
	lw33 := w - 2*w33

	preview.width = w
	preview.height = h60 - 3 // Makes room for the status bar
	sessions.width = w33
	sessions.height = h40
	windows.width = w33
	windows.height = h40
	frames.width = lw33
	frames.height = h40
	status.width = w
	status.height = 1

	previewRendered := preview.View(m.theme)
	sessionsRendered := sessions.View(m.theme)
	windowsRendered := windows.View(m.theme)
	framesRendered := frames.View(m.theme)
	statusRendered := status.View(m.theme)

	horizontalBox := lipgloss.JoinHorizontal(lipgloss.Left, sessionsRendered, windowsRendered, framesRendered)
	return lipgloss.JoinVertical(lipgloss.Top, previewRendered, horizontalBox, statusRendered)
}

type ListFrame struct {
	frame     Frame
	items     []TmuxEntity
	currentId int
	markedIds []int
	parentId   int
	filterText string
}

func (listFrame *ListFrame) Update() {
	visibleItems := listFrame.visibleItems()
	for _, item := range visibleItems {
		if item.id == listFrame.currentId {
			return
		}
	}
	if len(visibleItems) > 0 {
		listFrame.currentId = visibleItems[0].id
	} else {
		listFrame.currentId = -1
	}
}

func (listFrame *ListFrame) RenderContents(theme Theme) Frame {
	enumeratorStyle := lipgloss.NewStyle().Foreground(theme.accent).Background(theme.background)
	itemStyle := lipgloss.NewStyle().Foreground(theme.foreground).Background(theme.background)

	currentIndex := -1

	l := list.New().EnumeratorStyle(enumeratorStyle).ItemStyle(itemStyle)
	for i, item := range listFrame.visibleItems() {
		if slices.Contains(listFrame.markedIds, item.id) {
			l.Item(lipgloss.NewStyle().Foreground(theme.secondary).Render(fmt.Sprintf("[%d]: %s", item.id, item.name)))
		} else {
			l.Item(fmt.Sprintf("[%d]: %s", item.id, item.name))
		}
		if item.id == listFrame.currentId {
			currentIndex = i
		}
	}

	if currentIndex > -1 {
		enumerator := func(l list.Items, i int) string {
			if i == currentIndex {
				return "→ "
			}
			return ""
		}
		l = l.Enumerator(enumerator)
	}

	listFrame.frame.contents = l.String()
	return listFrame.frame
}

func (listFrame *ListFrame) SelectNext() {
	items := listFrame.visibleItems()
	if listFrame.currentId == -1 && len(items) > 0 {
		listFrame.currentId = 0
		return
	}
	for i, item := range items {
		if item.id == listFrame.currentId {
			if i+1 < len(items) {
				listFrame.currentId = items[i+1].id
			}
			return
		}
	}
	listFrame.currentId = -1
}

func (listFrame *ListFrame) SelectPrevious() {
	items := listFrame.visibleItems()
	if listFrame.currentId == -1 && len(items) > 0 {
		listFrame.currentId = items[len(items) - 1].id
		return
	}
	for i, item := range items {
		if item.id == listFrame.currentId {
			if i-1 >= 0 {
				listFrame.currentId = items[i-1].id
			}
			return
		}
	}
	listFrame.currentId = -1
}

func (listFrame *ListFrame) MarkSelection() {
	listFrame.markedIds = append(listFrame.markedIds, listFrame.currentId)
}

func (listFrame *ListFrame) UnmarkSelection() {
	index := slices.Index(listFrame.markedIds, listFrame.currentId)
	if index != -1 {
		listFrame.markedIds = slices.Delete(listFrame.markedIds, index, 1)
	}
}

func (listFrame *ListFrame) IsMarked(id int) bool {
	return slices.Contains(listFrame.markedIds, id)
}

func (listFrame *ListFrame) ClearMarks() {
	listFrame.markedIds = nil
}

func (listFrame ListFrame) ItemWithId(id int) *TmuxEntity {
	for _, item := range listFrame.items {
		if item.id == id {
			return &item
		}
	}
	return nil
}

func (listFrame *ListFrame) visibleItems() []TmuxEntity {
	var items []TmuxEntity
	filter := strings.ToLower(listFrame.filterText)
	for _, item := range listFrame.items {
		matchesParent := listFrame.parentId == -1 || item.parent == listFrame.parentId
		matchesFilter := len(filter) == 0 || strings.Contains(strings.ToLower(item.name), filter)
		if matchesParent && matchesFilter {
			items = append(items, item)
		}
	}
	return items
}
