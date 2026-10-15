package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/kevinforsgren/opentask/internal/storage"
	"github.com/kevinforsgren/opentask/internal/task"
	"github.com/kevinforsgren/opentask/internal/ui"
)

type appContext struct {
	service *task.Service
	styles  ui.Styles
	dbPath  string
}

var (
	dbPath  string
	noColor bool
	app     appContext
)

func Execute(ctx context.Context) error {
	root := newRootCommand(ctx)
	root.SetOut(os.Stdout)
	root.SetErr(os.Stderr)
	return root.Execute()
}

func newRootCommand(ctx context.Context) *cobra.Command {
	root := &cobra.Command{
		Use:           "task",
		Short:         "A fast, local personal task tracker",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			if dbPath == "" {
				dbPath = os.Getenv("TASK_DB_PATH")
			}
			if dbPath == "" {
				dbPath = "tasks.db"
			}
			db, err := storage.Open(ctx, dbPath)
			if err != nil {
				return err
			}
			cmd.Context()
			app = appContext{service: task.NewService(storage.NewTaskRepository(db)), styles: ui.NewStyles(noColor), dbPath: dbPath}
			return nil
		},
	}
	root.PersistentFlags().StringVar(&dbPath, "db", "", "database path (default: TASK_DB_PATH or ./tasks.db)")
	root.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable color and styling")
	root.AddCommand(newAddCommand(), newListCommand(), newUpdateCommand(), newDoneCommand(), newToggleCommand(), newDeleteCommand(), newDashboardCommand())
	_ = ctx
	return root
}

func ExitCode(err error) int {
	if errors.Is(err, task.ErrNotFound) || errors.Is(err, task.ErrInvalidID) || errors.Is(err, task.ErrInvalidInput) || errors.Is(err, task.ErrEmptyTitle) || errors.Is(err, task.ErrNoUpdateField) {
		return 2
	}
	return 1
}

func parseID(value string) (int64, error) {
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 {
		return 0, task.ErrInvalidID
	}
	return id, nil
}

func parseStatus(value string) (*task.Status, error) {
	status := task.Status(value)
	if !status.Valid() {
		return nil, fmt.Errorf("%w: invalid status %q", task.ErrInvalidInput, value)
	}
	return &status, nil
}

func parsePriority(value string) (*task.Priority, error) {
	priority := task.Priority(value)
	if !priority.Valid() {
		return nil, fmt.Errorf("%w: invalid priority %q", task.ErrInvalidInput, value)
	}
	return &priority, nil
}
