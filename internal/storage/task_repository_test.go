package storage

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/kevinforsgren/opentask/internal/task"
)

func TestTaskRepositoryCRUDAndFilters(t *testing.T) {
	ctx := context.Background()
	db, err := Open(ctx, filepath.Join(t.TempDir(), "tasks.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := NewTaskRepository(db)

	created, err := repo.Create(ctx, task.CreateInput{Title: "Plan release", Description: "project notes", Priority: task.PriorityHigh})
	if err != nil {
		t.Fatal(err)
	}
	if created.Status != task.StatusPending {
		t.Fatalf("status = %q, want pending", created.Status)
	}
	if _, err := repo.Create(ctx, task.CreateInput{Title: "Archive notes", Priority: task.PriorityLow}); err != nil {
		t.Fatal(err)
	}
	high := task.PriorityHigh
	items, err := repo.List(ctx, task.ListFilter{Priority: &high, Search: "release"})
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].ID != created.ID {
		t.Fatalf("filtered tasks = %#v", items)
	}
	updatedTitle := "Plan launch"
	updated, err := repo.Update(ctx, created.ID, task.UpdateInput{Title: &updatedTitle})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Title != updatedTitle {
		t.Fatalf("title = %q, want %q", updated.Title, updatedTitle)
	}
	if _, err := repo.SetStatus(ctx, created.ID, task.StatusDone); err != nil {
		t.Fatal(err)
	}
	if err := repo.Delete(ctx, created.ID); err != nil {
		t.Fatal(err)
	}
	if err := repo.Delete(ctx, created.ID); err != task.ErrNotFound {
		t.Fatalf("second delete error = %v, want %v", err, task.ErrNotFound)
	}
}
