package ui

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	marginBottom  = 5
	fileSizeWidth = 7
	paddingLeft   = 2
)

// Styles defines the possible customizations for styles in the file picker.
type Styles struct {
	DisabledCursor   lipgloss.Style
	Cursor           lipgloss.Style
	Symlink          lipgloss.Style
	Directory        lipgloss.Style
	File             lipgloss.Style
	DisabledFile     lipgloss.Style
	Permission       lipgloss.Style
	Selected         lipgloss.Style
	DisabledSelected lipgloss.Style
	FileSize         lipgloss.Style
	EmptyDirectory   lipgloss.Style

	// custom:
	PastTag    lipgloss.Style
	CurrentTag lipgloss.Style
	UiString   lipgloss.Style
	Help       lipgloss.Style
	Error      lipgloss.Style
}

// DefaultStyles defines the default styling for the file picker.
func DefaultStyles() Styles {
	return DefaultStylesWithRenderer(lipgloss.DefaultRenderer())
}

// DefaultStylesWithRenderer defines the default styling for the file picker,
// with a given Lip Gloss renderer.
func DefaultStylesWithRenderer(r *lipgloss.Renderer) Styles {
	return Styles{
		DisabledCursor:   r.NewStyle().Foreground(lipgloss.Color("247")),
		Cursor:           r.NewStyle().Foreground(lipgloss.Color("212")),
		Symlink:          r.NewStyle().Foreground(lipgloss.Color("36")),
		Directory:        r.NewStyle().Foreground(lipgloss.Color("99")),
		File:             r.NewStyle(),
		DisabledFile:     r.NewStyle().Foreground(lipgloss.Color("243")),
		DisabledSelected: r.NewStyle().Foreground(lipgloss.Color("247")),
		Permission:       r.NewStyle().Foreground(lipgloss.Color("244")),
		Selected:         r.NewStyle().Foreground(lipgloss.Color("212")).Bold(true), // Currently selected folder
		FileSize:         r.NewStyle().Foreground(lipgloss.Color("240")).Width(fileSizeWidth).Align(lipgloss.Right),
		EmptyDirectory:   r.NewStyle().Foreground(lipgloss.Color("240")).PaddingLeft(paddingLeft),
		PastTag:          r.NewStyle().Foreground(lipgloss.Color("244")),
		CurrentTag:       r.NewStyle().Foreground(lipgloss.Color("212")),
		UiString:         r.NewStyle().Foreground(lipgloss.Color("150")),
		Error:            r.NewStyle().Foreground(lipgloss.Color("9")),
		Help:             r.NewStyle().Foreground(lipgloss.Color("250")),
	}
}
