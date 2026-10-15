package ui

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevinforsgren/opentask/internal/task"
)

type Dashboard struct {
	ctx     context.Context
	service *task.Service
	styles  Styles
	tasks   []task.Task
	cursor  int
	err     error
	width   int
}

type tasksLoadedMsg struct {
	tasks []task.Task
	err   error
}

func RunDashboard(ctx context.Context, service *task.Service, styles Styles) error {
	model := Dashboard{ctx: ctx, service: service, styles: styles}
	program := tea.NewProgram(model, tea.WithAltScreen())
	_, err := program.Run()
	return err
}

func (m Dashboard) Init() tea.Cmd {
	return m.load
}

func (m Dashboard) load() tea.Msg {
	tasks, err := m.service.List(m.ctx, task.ListFilter{})
	return tasksLoadedMsg{tasks: tasks, err: err}
}

func (m Dashboard) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "r":
			return m, m.load
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.tasks)-1 {
				m.cursor++
			}
		case "enter", "d":
			if len(m.tasks) > 0 {
				id := m.tasks[m.cursor].ID
				return m, func() tea.Msg {
					_, err := m.service.Toggle(m.ctx, id)
					return tasksLoadedMsg{tasks: mustList(m.service, m.ctx), err: err}
				}
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case tasksLoadedMsg:
		m.tasks, m.err = msg.tasks, msg.err
		if m.cursor >= len(m.tasks) {
			m.cursor = max(0, len(m.tasks)-1)
		}
	}
	return m, nil
}

func (m Dashboard) View() string {
	if m.err != nil {
		return m.styles.Error.Render("Error: "+m.err.Error()) + "\n"
	}
	header := m.styles.Header.Render("OpenTask Dashboard")
	help := m.styles.Dim.Render("↑/↓ navigate  enter/d mark done  r refresh  q quit")
	body := ""
	for i, item := range m.tasks {
		cursor := "  "
		if i == m.cursor {
			cursor = "› "
		}
		status := "[ ]"
		if item.Status == task.StatusDone {
			status = "[✓]"
		}
		line := fmt.Sprintf("%s%s #%d %s", cursor, status, item.ID, item.Title)
		if i == m.cursor {
			line = m.styles.Header.Render(line)
		} else if item.Status == task.StatusDone {
			line = m.styles.Dim.Render(line)
		}
		body += line + "\n"
	}
	if body == "" {
		body = m.styles.Dim.Render("No tasks found.") + "\n"
	}
	return header + "\n\n" + body + "\n" + help
}

func mustList(service *task.Service, ctx context.Context) []task.Task {
	tasks, _ := service.List(ctx, task.ListFilter{})
	return tasks
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
