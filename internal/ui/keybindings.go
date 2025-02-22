package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type keyMap struct {
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

func DefaultkeyMap() keyMap {
	return keyMap{
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
		ApplyTag: key.NewBinding(key.WithKeys("m"), key.WithHelp("m", "apply tag to current folder and children")),
		Help:     key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.GoToTop, k.GoToLast, k.Up, k.Down},
		{k.PageUp, k.PageDown, k.Back, k.Open},
		{k.Help, k.Quit, k.EditTag, k.ApplyTag},
	}
}

type InputkeyMap struct {
	Accept key.Binding
	Cancel key.Binding
}

func DefaultInputkeyMap() InputkeyMap {
	return InputkeyMap{
		Accept: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "accept")),
		Cancel: key.NewBinding(key.WithKeys("esc")),
	}
}

func basicKeyHandler(m *Model, msg tea.KeyMsg) tea.Cmd {
	if m.editMode {
		return editModeHandler(m, msg)
	}
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
	case key.Matches(msg, m.keyMap.Quit):
		return handleQuit(m)
	case key.Matches(msg, m.keyMap.EditTag):
		return handleEditTag(m)
	case key.Matches(msg, m.keyMap.ApplyTag):
		return handleApplyTag(m)
	case key.Matches(msg, m.keyMap.Help):
		return handleHelp(m)
	}
	return nil
}

func editModeHandler(m *Model, msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, m.inputkeyMap.Accept):
		return handleAccept(m)
	case key.Matches(msg, m.inputkeyMap.Cancel):
		return handleCancel(m)
	default:
		return handleInput(m, msg)
	}
}
