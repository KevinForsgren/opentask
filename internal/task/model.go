package task

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Status string

const (
	StatusPending Status = "pending"
	StatusDone    Status = "done"
)

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

var (
	ErrNotFound      = errors.New("task not found")
	ErrInvalidInput  = errors.New("invalid task input")
	ErrInvalidID     = errors.New("task ID must be positive")
	ErrEmptyTitle    = errors.New("task title cannot be empty")
	ErrNoUpdateField = errors.New("at least one update field is required")
)

type Task struct {
	ID          int64
	Title       string
	Description string
	Status      Status
	Priority    Priority
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateInput struct {
	Title       string
	Description string
	Priority    Priority
}

type UpdateInput struct {
	Title       *string
	Description *string
	Priority    *Priority
}

type ListFilter struct {
	Status   *Status
	Priority *Priority
	Search   string
}

func (s Status) Valid() bool {
	return s == StatusPending || s == StatusDone
}

func (p Priority) Valid() bool {
	return p == PriorityLow || p == PriorityMedium || p == PriorityHigh
}

func ValidateCreate(input CreateInput) error {
	if strings.TrimSpace(input.Title) == "" {
		return ErrEmptyTitle
	}
	if input.Priority == "" {
		input.Priority = PriorityMedium
	}
	if !input.Priority.Valid() {
		return fmt.Errorf("%w: invalid priority %q", ErrInvalidInput, input.Priority)
	}
	return nil
}

func ValidateUpdate(input UpdateInput) error {
	if input.Title == nil && input.Description == nil && input.Priority == nil {
		return ErrNoUpdateField
	}
	if input.Title != nil && strings.TrimSpace(*input.Title) == "" {
		return ErrEmptyTitle
	}
	if input.Priority != nil && !input.Priority.Valid() {
		return fmt.Errorf("%w: invalid priority %q", ErrInvalidInput, *input.Priority)
	}
	return nil
}

func ValidateFilter(filter ListFilter) error {
	if filter.Status != nil && !filter.Status.Valid() {
		return fmt.Errorf("%w: invalid status %q", ErrInvalidInput, *filter.Status)
	}
	if filter.Priority != nil && !filter.Priority.Valid() {
		return fmt.Errorf("%w: invalid priority %q", ErrInvalidInput, *filter.Priority)
	}
	return nil
}
