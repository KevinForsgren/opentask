package task

import (
	"context"
	"fmt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (Task, error) {
	if input.Priority == "" {
		input.Priority = PriorityMedium
	}
	if err := ValidateCreate(input); err != nil {
		return Task{}, err
	}
	return s.repo.Create(ctx, input)
}

func (s *Service) List(ctx context.Context, filter ListFilter) ([]Task, error) {
	if err := ValidateFilter(filter); err != nil {
		return nil, err
	}
	return s.repo.List(ctx, filter)
}

func (s *Service) Get(ctx context.Context, id int64) (Task, error) {
	if id <= 0 {
		return Task{}, ErrInvalidID
	}
	return s.repo.Get(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (Task, error) {
	if id <= 0 {
		return Task{}, ErrInvalidID
	}
	if err := ValidateUpdate(input); err != nil {
		return Task{}, err
	}
	return s.repo.Update(ctx, id, input)
}

func (s *Service) Done(ctx context.Context, id int64) (Task, error) {
	return s.SetStatus(ctx, id, StatusDone)
}

func (s *Service) Toggle(ctx context.Context, id int64) (Task, error) {
	current, err := s.Get(ctx, id)
	if err != nil {
		return Task{}, err
	}
	status := StatusDone
	if current.Status == StatusDone {
		status = StatusPending
	}
	return s.SetStatus(ctx, id, status)
}

func (s *Service) SetStatus(ctx context.Context, id int64, status Status) (Task, error) {
	if id <= 0 {
		return Task{}, ErrInvalidID
	}
	if !status.Valid() {
		return Task{}, fmt.Errorf("%w: invalid status %q", ErrInvalidInput, status)
	}
	return s.repo.SetStatus(ctx, id, status)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrInvalidID
	}
	return s.repo.Delete(ctx, id)
}
