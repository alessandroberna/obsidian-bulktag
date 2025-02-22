package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"

	"obsidian-tagfmt/internal/tag"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// needed to manage the navigation history
type stack struct {
	Push   func(int)
	Pop    func() int
	Length func() int
}

func newStack() stack {
	slice := make([]int, 0)
	return stack{
		Push: func(i int) {
			slice = append(slice, i)
		},
		Pop: func() int {
			res := slice[len(slice)-1]
			slice = slice[:len(slice)-1]
			return res
		},
		Length: func() int {
			return len(slice)
		},
	}
}

func (m *Model) pushView(selected, min, max int) {
	m.selectedStack.Push(selected)
	m.minStack.Push(min)
	m.maxStack.Push(max)
}

func (m *Model) popView() (int, int, int) {
	return m.selectedStack.Pop(), m.minStack.Pop(), m.maxStack.Pop()
}

type Model struct {
	id int

	Root         string   // root path
	AllowedTypes []string // allowed file types
	ShowFiles    bool     // true if showing files

	path      string          // current path
	entries   []os.DirEntry   // list of entries in current folder
	cursorPos int             // cursor position
	editMode  bool            // true if editing tag
	quitting  bool            // true if quitting
	textInput textinput.Model // text input for tag editing

	min           int
	max           int
	minStack      stack
	maxStack      stack
	selectedStack stack
	help          help.Model

	height     int
	autoHeight bool

	cursorGlyph string
	styles      Styles
	keyMap      keyMap
	inputkeyMap InputkeyMap

	Tags    *tag.FolderTag // list of tags
	Message string         // status or feedback message
}

type errorMsg struct {
	err error
}

var lastID int64

func nextID() int {
	return int(atomic.AddInt64(&lastID, 1))
}

func New() Model {
	ti := textinput.New()
	ti.TextStyle = DefaultStyles().CurrentTag
	ti.Cursor.Style = DefaultStyles().CurrentTag
	ti.CompletionStyle = DefaultStyles().CurrentTag
	ti.PromptStyle = DefaultStyles().PastTag
	ti.PlaceholderStyle = DefaultStyles().PastTag
	return Model{
		id:            nextID(),
		Root:          ".",
		AllowedTypes:  []string{},
		ShowFiles:     false,
		path:          ".",
		cursorPos:     0,
		editMode:      false,
		quitting:      false,
		textInput:     ti,
		min:           0,
		max:           0,
		help:          help.New(),
		Message:       "",
		minStack:      newStack(),
		maxStack:      newStack(),
		selectedStack: newStack(),
		height:        0,
		autoHeight:    true,
		cursorGlyph:   ">",
		styles:        DefaultStyles(),
		keyMap:        DefaultkeyMap(),
		inputkeyMap:   DefaultInputkeyMap(),
		Tags:          tag.NewTagGetter(".", nil),
	}
}

type readDirMsg struct {
	id      int
	entries []os.DirEntry
}

func (m *Model) readDir(path string) tea.Cmd {
	return func() tea.Msg {
		entries, err := os.ReadDir(path)
		if err != nil {
			return errorMsg{err}
		}
		var dirs []os.DirEntry
		var files []os.DirEntry
		for _, entry := range entries {
			if entry.IsDir() && evalEntry(entry.Name()) {
				dirs = append(dirs, entry)
			}
			if m.ShowFiles && !entry.IsDir() {
				files = append(files, entry)
			}
		}
		finalList := append(dirs, files...)
		return readDirMsg{id: m.id, entries: finalList}
	}
}

// Init initializes the file picker Model.
func (m Model) Init() tea.Cmd {
	return m.readDir(m.Root)
}

func (m Model) canSelect(file string) bool {
	if len(m.AllowedTypes) <= 0 {
		return true
	}

	for _, ext := range m.AllowedTypes {
		if strings.HasSuffix(file, ext) {
			return true
		}
	}
	return false
}

func evalEntry(entry string) bool {
	// ignores hidden directories and obsidian image folders (attachments, resources)
	// TODO: make this configurable
	if strings.HasPrefix(entry, ".") || entry == "attachments" || entry == "resources" {
		return false
	}
	return true
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errorMsg:
		m.Message = fmt.Sprintf("Error reading directory: %v", msg.err)
	case readDirMsg:
		if msg.id == m.id {
			m.entries = msg.entries
			m.max = max(m.max, m.height-1)
		}
	case tea.WindowSizeMsg:
		if m.autoHeight {
			m.height = msg.Height - marginBottom
		}
		m.help.Width = msg.Width
		m.max = m.height - 1
	case tea.KeyMsg:
		cmd := basicKeyHandler(&m, msg)
		return m, cmd
	}
	return m, nil
}

func styledConditionalSlashJoin(string1 string, string2 string, style1 lipgloss.Style, style2 lipgloss.Style) string {
	if string1 == "" {
		return style2.Render(string2)
	}
	if string2 == "" {
		return style1.Render(string1)
	}
	return style1.Render(string1+"/") + style2.Render(string2)
}

func conditionalSlashAppend(s string) string {
	if s == "" {
		return ""
	}
	return s + "/"
}

// View returns the view of the file picker.
func (m Model) View() string {
	// TODO: split into smaller functions
	if m.quitting {
		return ""
	}


	var s strings.Builder
	if len(m.entries) == 0 {
		if m.ShowFiles {
			s.WriteString(m.styles.EmptyDirectory.Render("This directory is empty."))
		} else {
			//s.WriteString(m.styles.EmptyDirectory.Render("This directory is empty. Press 'f' to show files."))
			s.WriteString(m.styles.EmptyDirectory.Render("No subdirectories found."))
		//return m.styles.EmptyDirectory.Height(m.height).MaxHeight(m.height).String()
		}
	} else {
		for i, f := range m.entries {
			if i < m.min || i > m.max-2 {
				continue
			}

			var symlinkPath string
			info, _ := f.Info()
			isSymlink := info.Mode()&os.ModeSymlink != 0
			name := f.Name()

			if isSymlink {
				symlinkPath, _ = filepath.EvalSymlinks(filepath.Join(m.path, name))
			}

			disabled := !m.canSelect(name) && !f.IsDir()

			if m.cursorPos == i { //nolint:nestif
				selected := ""
				selected += " " + name
				if isSymlink {
					selected += " → " + symlinkPath
				}
				if disabled {
					s.WriteString(m.styles.DisabledSelected.Render(m.cursorGlyph) + m.styles.DisabledSelected.Render(selected))
				} else {
					s.WriteString(m.styles.Cursor.Render(m.cursorGlyph) + m.styles.Selected.Render(selected))
				}
				s.WriteRune('\n')
				continue
			}

			style := m.styles.File
			if f.IsDir() {
				style = m.styles.Directory
			} else if isSymlink {
				style = m.styles.Symlink
			} else if disabled {
				style = m.styles.DisabledFile
			}

			fileName := style.Render(name)
			s.WriteString(m.styles.Cursor.Render(" "))
			if isSymlink {
				fileName += " → " + symlinkPath
			}
			s.WriteString(" " + fileName)
			s.WriteRune('\n')
		}
	}
	if m.editMode {
		//s.WriteString("\nEditing tag: " + styledConditionalSlashJoin(m.Tags.parentTagsStr(), m.textInput.View(), m.styles.PastTag, m.styles.CurrentTag))
		s.WriteString(m.styles.UiString.Render("\nEditing tag: "))
		m.textInput.Prompt = conditionalSlashAppend(m.Tags.ParentTagsStr())
		s.WriteString(m.textInput.View() + "\n")
	} else {
		s.WriteString(m.styles.UiString.Render("\nCurrent Tag: ") + styledConditionalSlashJoin(m.Tags.ParentTagsStr(), m.Tags.Tag, m.styles.PastTag, m.styles.CurrentTag) + "\n")
	}

	if m.Message != "" {
		s.WriteString(m.styles.Error.Render(m.Message) + "\n")
	}
	s.WriteString(m.help.View(m.keyMap))
	// Add padding to the bottom of the list
	for i := lipgloss.Height(s.String()); i <= m.height; i++ {
		s.WriteRune('\n')
	}

	return s.String()
}
