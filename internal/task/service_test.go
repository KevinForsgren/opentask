package task

import (
	"context"
	"testing"
)

type memoryRepository struct {
	next  int64
	items map[int64]Task
}

func (r *memoryRepository) Create(_ context.Context, input CreateInput) (Task, error) {
	r.next++
	item := Task{ID: r.next, Title: input.Title, Priority: input.Priority, Status: StatusPending}
	r.items[item.ID] = item
	return item, nil
}
func (r *memoryRepository) List(_ context.Context, _ ListFilter) ([]Task, error) { return nil, nil }
func (r *memoryRepository) Get(_ context.Context, id int64) (Task, error) {
	item, ok := r.items[id]
	if !ok {
		return Task{}, ErrNotFound
	}
	return item, nil
}
func (r *memoryRepository) Update(_ context.Context, id int64, input UpdateInput) (Task, error) {
	item, err := r.Get(context.Background(), id)
	if err != nil {
		return Task{}, err
	}
	if input.Title != nil {
		item.Title = *input.Title
	}
	r.items[id] = item
	return item, nil
}
func (r *memoryRepository) SetStatus(_ context.Context, id int64, status Status) (Task, error) {
	item, err := r.Get(context.Background(), id)
	if err != nil {
		return Task{}, err
	}
	item.Status = status
	r.items[id] = item
	return item, nil
}
func (r *memoryRepository) Delete(_ context.Context, id int64) error {
	if _, ok := r.items[id]; !ok {
		return ErrNotFound
	}
	delete(r.items, id)
	return nil
}

func TestServiceValidationAndToggle(t *testing.T) {
	repo := &memoryRepository{items: map[int64]Task{}}
	service := NewService(repo)
	ctx := context.Background()

	item, err := service.Create(ctx, CreateInput{Title: "Ship it"})
	if err != nil {
		t.Fatal(err)
	}
	if item.Priority != PriorityMedium {
		t.Fatalf("default priority = %q, want %q", item.Priority, PriorityMedium)
	}
	if _, err := service.Create(ctx, CreateInput{Title: " "}); err != ErrEmptyTitle {
		t.Fatalf("empty title error = %v, want %v", err, ErrEmptyTitle)
	}
	item, err = service.Toggle(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if item.Status != StatusDone {
		t.Fatalf("toggled status = %q, want %q", item.Status, StatusDone)
	}
}
