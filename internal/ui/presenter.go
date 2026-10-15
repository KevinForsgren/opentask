package ui

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/kevinforsgren/opentask/internal/task"
)

func RenderTasks(w io.Writer, tasks []task.Task, styles Styles) error {
	if len(tasks) == 0 {
		_, err := io.WriteString(w, styles.Dim.Render("No tasks found.")+"\n")
		return err
	}
	headers := []string{"", "ID", "TITLE", "PRIORITY", "UPDATED"}
	rows := make([][]string, 0, len(tasks))
	for _, item := range tasks {
		mark, title := "[ ]", item.Title
		if item.Status == task.StatusDone {
			mark = "[✓]"
			title = styles.Dim.Render(title)
		}
		priority := priorityGlyph(item.Priority) + " " + string(item.Priority)
		priority = styles.Priority(item.Priority).Render(priority)
		rows = append(rows, []string{mark, strconv.FormatInt(item.ID, 10), title, priority, item.UpdatedAt.Local().Format("2006-01-02 15:04")})
	}
	widths := make([]int, len(headers))
	for i, header := range headers {
		widths[i] = lipgloss.Width(header)
	}
	for _, row := range rows {
		for i, value := range row {
			if width := lipgloss.Width(value); width > widths[i] {
				widths[i] = width
			}
		}
	}
	writeRow := func(row []string, separator bool) error {
		if separator {
			parts := make([]string, len(widths))
			for i, width := range widths {
				parts[i] = strings.Repeat("─", width)
			}
			_, err := io.WriteString(w, strings.Join(parts, "  ")+"\n")
			return err
		}
		parts := make([]string, len(row))
		for i, value := range row {
			parts[i] = value + strings.Repeat(" ", widths[i]-lipgloss.Width(value))
		}
		_, err := io.WriteString(w, strings.Join(parts, "  ")+"\n")
		return err
	}
	if err := writeRow(headers, false); err != nil {
		return err
	}
	if err := writeRow(nil, true); err != nil {
		return err
	}
	for _, row := range rows {
		if err := writeRow(row, false); err != nil {
			return err
		}
	}
	return nil
}

func RenderTask(w io.Writer, item task.Task, styles Styles) error {
	status := "[ ]"
	if item.Status == task.StatusDone {
		status = "[✓]"
	}
	_, err := fmt.Fprintf(w, "%s #%d %s (%s)\n", status, item.ID, item.Title, item.Priority)
	return err
}

func priorityGlyph(priority task.Priority) string {
	switch priority {
	case task.PriorityHigh:
		return "⚡"
	case task.PriorityMedium:
		return "•"
	default:
		return "·"
	}
}
