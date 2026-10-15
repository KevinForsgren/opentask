package ui

import (
	"io"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/kevinforsgren/opentask/internal/task"
)

type Styles struct {
	Header   lipgloss.Style
	Border   lipgloss.Style
	Dim      lipgloss.Style
	Success  lipgloss.Style
	Error    lipgloss.Style
	Priority func(task.Priority) lipgloss.Style
	NoColor  bool
}

func NewStyles(noColor bool) Styles {
	base := Styles{NoColor: noColor}
	base.Header = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	base.Border = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1)
	base.Dim = lipgloss.NewStyle().Faint(true)
	base.Success = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	base.Error = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	base.Priority = func(priority task.Priority) lipgloss.Style {
		switch priority {
		case task.PriorityHigh:
			return lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
		case task.PriorityMedium:
			return lipgloss.NewStyle().Foreground(lipgloss.Color("220"))
		default:
			return lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
		}
	}
	if noColor || !isTerminal(os.Stdout) {
		base.Header = lipgloss.NewStyle().Bold(true)
		base.Border = lipgloss.NewStyle()
		base.Dim = lipgloss.NewStyle()
		base.Success = lipgloss.NewStyle()
		base.Error = lipgloss.NewStyle()
		base.Priority = func(task.Priority) lipgloss.Style { return lipgloss.NewStyle() }
	}
	return base
}

func isTerminal(file *os.File) bool {
	info, err := file.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}

func PrintSuccess(w io.Writer, message string, styles Styles) error {
	_, err := io.WriteString(w, styles.Success.Render("✓ "+message)+"\n")
	return err
}

func PrintError(w io.Writer, message string, styles Styles) error {
	_, err := io.WriteString(w, styles.Error.Render("✗ "+message)+"\n")
	return err
}
