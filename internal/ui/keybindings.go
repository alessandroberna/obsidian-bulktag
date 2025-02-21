package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type KeyMap struct {
	// from bubbles/filepicker
	GoToTop  key.Binding
	GoToLast key.Binding
	Down     key.Binding
	Up       key.Binding
	PageUp   key.Binding
	PageDown key.Binding
	Back     key.Binding
	Open     key.Binding
	//Select   key.Binding

	// custom
	Quit     key.Binding
	EditTag  key.Binding
	ApplyTag key.Binding
	Help     key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		GoToTop:  key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "go to first")),
		GoToLast: key.NewBinding(key.WithKeys("G"), key.WithHelp("G", "go to last")),
		Down:     key.NewBinding(key.WithKeys("s", "j", "down", "ctrl+n"), key.WithHelp("↓/j/s/ctrl+n", "move down")),
		Up:       key.NewBinding(key.WithKeys("w", "k", "up", "ctrl+p"), key.WithHelp("↑/k/w/ctrl+p", "move up")),
		PageUp:   key.NewBinding(key.WithKeys("K", "pgup"), key.WithHelp("pgup/K", "page up")),
		PageDown: key.NewBinding(key.WithKeys("J", "pgdown"), key.WithHelp("pgdown/J", "page down")),
		Back:     key.NewBinding(key.WithKeys("a", "h", "backspace", "left", "esc"), key.WithHelp("←/h/a/backspace/esc", "move back")),
		Open:     key.NewBinding(key.WithKeys("d", "l", "right", "enter"), key.WithHelp("→/l/d/enter", "open")),
		//Select:   key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")), // not used
		Quit:     key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
		EditTag:  key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "edit current folder's tag")),
		ApplyTag: key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "apply tag to current folder and children")),
		Help:     key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
	}
}

func basicKeyHandler(m *Model, msg tea.KeyMsg) (tea.Cmd) {
	switch {
	case key.Matches(msg, m.keyMap.GoToTop):
		return handleGoToTop(m)
	case key.Matches(msg, m.keyMap.GoToLast):
		return handleGoToLast(m)
	case key.Matches(msg, m.keyMap.Down):
		return handleDown(m)
	case key.Matches(msg, m.keyMap.Up):
		return handleUp(m)
	case key.Matches(msg, m.keyMap.PageDown):
		return handlePageDown(m)
	case key.Matches(msg, m.keyMap.PageUp):
		return handlePageUp(m)
	case key.Matches(msg, m.keyMap.Back):
		return handleBack(m)
	case key.Matches(msg, m.keyMap.Open):
		return handleOpen(m)
	case key.Matches (msg, m.keyMap.Quit):
		return handleQuit(m)
	}
	return nil
}
