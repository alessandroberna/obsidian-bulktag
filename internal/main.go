package internal

import (
	"obsidian-bulktag/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

type ConfigSt struct {
	Root      string
	ShowFiles bool
}

var Config ConfigSt

func Main() error {
	fp := ui.New()
	fp.Root = Config.Root
	fp.ShowFiles = Config.ShowFiles

	tm := tea.NewProgram(fp)
	tm.Run()
	return nil
}
