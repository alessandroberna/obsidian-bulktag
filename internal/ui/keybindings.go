package ui

import (
	"os"
	"path/filepath"

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

func basicKeyHandler(m *Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keyMap.GoToTop):
		m.cursorPos = 0
		m.min = 0
		m.max = m.height - 1
	case key.Matches(msg, m.keyMap.GoToLast):
		m.cursorPos = len(m.entries) - 1
		m.min = len(m.entries) - m.height
		m.max = len(m.entries) - 1
	case key.Matches(msg, m.keyMap.Down):
		m.cursorPos++
		if m.cursorPos >= len(m.entries) {
			m.cursorPos = len(m.entries) - 1
		}
		if m.cursorPos > m.max {
			m.min++
			m.max++
		}
	case key.Matches(msg, m.keyMap.Up):
		m.cursorPos--
		if m.cursorPos < 0 {
			m.cursorPos = 0
		}
		if m.cursorPos < m.min {
			m.min--
			m.max--
		}
	case key.Matches(msg, m.keyMap.PageDown):
		m.cursorPos += m.height
		if m.cursorPos >= len(m.entries) {
			m.cursorPos = len(m.entries) - 1
		}
		m.min += m.height
		m.max += m.height

		if m.max >= len(m.entries) {
			m.max = len(m.entries) - 1
			m.min = m.max - m.height
		}
	case key.Matches(msg, m.keyMap.PageUp):
		m.cursorPos -= m.height
		if m.cursorPos < 0 {
			m.cursorPos = 0
		}
		m.min -= m.height
		m.max -= m.height

		if m.min < 0 {
			m.min = 0
			m.max = m.min + m.height
		}
	case key.Matches(msg, m.keyMap.Back):
		m.path = filepath.Dir(m.path)
		if m.selectedStack.Length() > 0 {
			m.cursorPos, m.min, m.max = m.popView()
		} else {
			m.cursorPos = 0
			m.min = 0
			m.max = m.height - 1
		}
		return m, m.readDir(m.path)
	case key.Matches(msg, m.keyMap.Open):
		if len(m.entries) == 0 {
			break
		}

		f := m.entries[m.cursorPos]
		info, err := f.Info()
		if err != nil {
			break
		}
		isDir := f.IsDir()
		if info.Mode()&os.ModeSymlink != 0 { // if it's a symlink
			symlinkPath, _ := filepath.EvalSymlinks(filepath.Join(m.path, f.Name()))
			info, err := os.Stat(symlinkPath)
			if err != nil {
				break
			}
			if info.IsDir() {
				isDir = true
			}
		}
		if isDir {
			m.path = filepath.Join(m.path, f.Name())
			m.pushView(m.cursorPos, m.min, m.max)
			m.cursorPos = 0
			m.min = 0
			m.max = m.height - 1
			return m, m.readDir(m.path)
		}
		// might want to add a markdown renderer here
	}
	return m, nil
}
