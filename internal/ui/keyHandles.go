package ui

import (
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

func updateTextInput(m *Model) {
	if m.Tags.Tag == "" {
		m.textInput.Reset()
	} else {
		m.textInput.SetValue(m.Tags.Tag)
		m.textInput.CursorEnd()
	}
}

func handleGoToTop(m *Model) tea.Cmd {
	m.cursorPos = 0
	m.min = 0
	m.max = m.height - 1
	return nil
}

func handleGoToLast(m *Model) tea.Cmd {
	m.cursorPos = len(m.entries) - 1
	m.min = len(m.entries) - m.height
	m.max = len(m.entries) - 1
	return nil
}

func handleDown(m *Model) tea.Cmd {
	m.cursorPos++
	if m.cursorPos >= len(m.entries) {
		m.cursorPos = len(m.entries) - 1
	}
	if m.cursorPos > m.max {
		m.min++
		m.max++
	}
	return nil
}

func handleUp(m *Model) tea.Cmd {
	m.cursorPos--
	if m.cursorPos < 0 {
		m.cursorPos = 0
	}
	if m.cursorPos < m.min {
		m.min--
		m.max--
	}
	return nil
}

func handlePageDown(m *Model) tea.Cmd {
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
	return nil
}

func handlePageUp(m *Model) tea.Cmd {
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
	return nil
}

func handleBack(m *Model) tea.Cmd {
	m.path = filepath.Dir(m.path)
	if m.selectedStack.Length() > 0 {
		m.cursorPos, m.min, m.max = m.popView()
	} else {
		m.cursorPos = 0
		m.min = 0
		m.max = m.height - 1
	}
	if m.Tags.Parent != nil {
		m.Tags = m.Tags.Parent
	}
	updateTextInput(m)
	return m.readDir(m.path)
}

func handleOpen(m *Model) tea.Cmd {
	if len(m.entries) == 0 {
		return nil
	}

	f := m.entries[m.cursorPos]
	info, err := f.Info()
	if err != nil {
		return nil
	}
	isDir := f.IsDir()
	if info.Mode()&os.ModeSymlink != 0 { // if it's a symlink
		symlinkPath, _ := filepath.EvalSymlinks(filepath.Join(m.path, f.Name()))
		info, err := os.Stat(symlinkPath)
		if err != nil {
			return nil
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
		m.Tags = newTagGetter(m.path, m.Tags)
		updateTextInput(m)
		return m.readDir(m.path)
	}
	return nil
}

func handleQuit(m *Model) tea.Cmd {
	m.quitting = true
	return tea.Quit
}

func handleEditTag(m *Model) tea.Cmd {
	m.editMode = true
	m.textInput.Focus()
	return nil
}

func handleAccept(m *Model) tea.Cmd {
	m.editMode = false
	m.Tags.Tag = m.textInput.Value()
	return nil
}

func handleCancel(m *Model) tea.Cmd {
	m.editMode = false
	m.textInput.Reset()
	return nil
}

func handleInput(m *Model, msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return cmd
}