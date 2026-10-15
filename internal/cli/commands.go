package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/kevinforsgren/opentask/internal/task"
	"github.com/kevinforsgren/opentask/internal/ui"
)

func newAddCommand() *cobra.Command {
	var description, priority string
	cmd := &cobra.Command{
		Use:   `add "title"`,
		Short: "Create a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			item, err := app.service.Create(cmd.Context(), task.CreateInput{Title: args[0], Description: description, Priority: task.Priority(priority)})
			if err != nil {
				return err
			}
			return ui.PrintSuccess(cmd.OutOrStdout(), fmt.Sprintf("Created task #%d", item.ID), app.styles)
		},
	}
	cmd.Flags().StringVarP(&description, "description", "d", "", "optional task description")
	cmd.Flags().StringVarP(&priority, "priority", "p", "medium", "priority: low, medium, or high")
	return cmd
}

func newListCommand() *cobra.Command {
	var status, priority, search string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List tasks",
		RunE: func(cmd *cobra.Command, _ []string) error {
			var filter task.ListFilter
			var err error
			if status != "" {
				filter.Status, err = parseStatus(status)
				if err != nil {
					return err
				}
			}
			if priority != "" {
				filter.Priority, err = parsePriority(priority)
				if err != nil {
					return err
				}
			}
			filter.Search = search
			items, err := app.service.List(cmd.Context(), filter)
			if err != nil {
				return err
			}
			return ui.RenderTasks(cmd.OutOrStdout(), items, app.styles)
		},
	}
	cmd.Flags().StringVar(&status, "status", "", "filter by status: pending or done")
	cmd.Flags().StringVarP(&priority, "priority", "p", "", "filter by priority: low, medium, or high")
	cmd.Flags().StringVarP(&search, "search", "s", "", "search title and description")
	return cmd
}

func newUpdateCommand() *cobra.Command {
	var title, description, priority string
	var descriptionSet bool
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseID(args[0])
			if err != nil {
				return err
			}
			input := task.UpdateInput{}
			if cmd.Flags().Changed("title") {
				input.Title = &title
			}
			if cmd.Flags().Changed("description") || descriptionSet {
				input.Description = &description
			}
			if cmd.Flags().Changed("priority") {
				parsed, err := parsePriority(priority)
				if err != nil {
					return err
				}
				input.Priority = parsed
			}
			item, err := app.service.Update(cmd.Context(), id, input)
			if err != nil {
				return err
			}
			return ui.PrintSuccess(cmd.OutOrStdout(), fmt.Sprintf("Updated task #%d", item.ID), app.styles)
		},
	}
	cmd.Flags().StringVar(&title, "title", "", "new title")
	cmd.Flags().StringVarP(&description, "description", "d", "", "new description (empty clears it)")
	cmd.Flags().StringVarP(&priority, "priority", "p", "", "new priority")
	cmd.Flags().BoolVar(&descriptionSet, "clear-description", false, "clear the description")
	return cmd
}

func newDoneCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "done <id>",
		Short: "Mark a task done",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseID(args[0])
			if err != nil {
				return err
			}
			item, err := app.service.Done(cmd.Context(), id)
			if err != nil {
				return err
			}
			return ui.PrintSuccess(cmd.OutOrStdout(), fmt.Sprintf("Completed task #%d", item.ID), app.styles)
		},
	}
}

func newToggleCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "toggle <id>",
		Short: "Toggle a task between pending and done",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseID(args[0])
			if err != nil {
				return err
			}
			item, err := app.service.Toggle(cmd.Context(), id)
			if err != nil {
				return err
			}
			return ui.PrintSuccess(cmd.OutOrStdout(), fmt.Sprintf("Toggled task #%d to %s", item.ID, item.Status), app.styles)
		},
	}
}

func newDeleteCommand() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseID(args[0])
			if err != nil {
				return err
			}
			if !force {
				return fmt.Errorf("deleting tasks requires --force")
			}
			if err := app.service.Delete(cmd.Context(), id); err != nil {
				return err
			}
			return ui.PrintSuccess(cmd.OutOrStdout(), fmt.Sprintf("Deleted task #%d", id), app.styles)
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "confirm deletion")
	return cmd
}

func newDashboardCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "dashboard",
		Short: "Open the interactive task dashboard",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return ui.RunDashboard(cmd.Context(), app.service, app.styles)
		},
	}
}

func normalizeDescription(value string) string {
	return strings.TrimSpace(value)
}
