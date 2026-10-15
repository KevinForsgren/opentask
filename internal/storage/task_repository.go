package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/kevinforsgren/opentask/internal/task"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, input task.CreateInput) (task.Task, error) {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	result, err := r.db.ExecContext(ctx, `INSERT INTO tasks (title, description, priority, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
		strings.TrimSpace(input.Title), input.Description, input.Priority, now, now)
	if err != nil {
		return task.Task{}, fmt.Errorf("create task: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return task.Task{}, fmt.Errorf("read created task ID: %w", err)
	}
	return r.Get(ctx, id)
}

func (r *TaskRepository) List(ctx context.Context, filter task.ListFilter) ([]task.Task, error) {
	query := `SELECT id, title, description, status, priority, created_at, updated_at FROM tasks WHERE 1=1`
	args := make([]any, 0, 3)
	if filter.Status != nil {
		query += ` AND status = ?`
		args = append(args, *filter.Status)
	}
	if filter.Priority != nil {
		query += ` AND priority = ?`
		args = append(args, *filter.Priority)
	}
	if strings.TrimSpace(filter.Search) != "" {
		query += ` AND (title LIKE ? ESCAPE '\' OR description LIKE ? ESCAPE '\')`
		search := "%" + escapeLike(strings.TrimSpace(filter.Search)) + "%"
		args = append(args, search, search)
	}
	query += ` ORDER BY CASE status WHEN 'pending' THEN 0 ELSE 1 END, CASE priority WHEN 'high' THEN 0 WHEN 'medium' THEN 1 ELSE 2 END, updated_at DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	defer rows.Close()
	var tasks []task.Task
	for rows.Next() {
		item, err := scanTask(rows)
		if err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		tasks = append(tasks, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tasks: %w", err)
	}
	return tasks, nil
}

func (r *TaskRepository) Get(ctx context.Context, id int64) (task.Task, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, title, description, status, priority, created_at, updated_at FROM tasks WHERE id = ?`, id)
	item, err := scanTask(row)
	if errors.Is(err, sql.ErrNoRows) {
		return task.Task{}, task.ErrNotFound
	}
	if err != nil {
		return task.Task{}, fmt.Errorf("get task: %w", err)
	}
	return item, nil
}

func (r *TaskRepository) Update(ctx context.Context, id int64, input task.UpdateInput) (task.Task, error) {
	sets := make([]string, 0, 4)
	args := make([]any, 0, 4)
	if input.Title != nil {
		sets = append(sets, "title = ?")
		args = append(args, strings.TrimSpace(*input.Title))
	}
	if input.Description != nil {
		sets = append(sets, "description = ?")
		args = append(args, *input.Description)
	}
	if input.Priority != nil {
		sets = append(sets, "priority = ?")
		args = append(args, *input.Priority)
	}
	sets = append(sets, "updated_at = ?")
	args = append(args, time.Now().UTC().Format(time.RFC3339Nano), id)
	result, err := r.db.ExecContext(ctx, "UPDATE tasks SET "+strings.Join(sets, ", ")+" WHERE id = ?", args...)
	if err != nil {
		return task.Task{}, fmt.Errorf("update task: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return task.Task{}, fmt.Errorf("check updated task: %w", err)
	}
	if affected == 0 {
		return task.Task{}, task.ErrNotFound
	}
	return r.Get(ctx, id)
}

func (r *TaskRepository) SetStatus(ctx context.Context, id int64, status task.Status) (task.Task, error) {
	result, err := r.db.ExecContext(ctx, `UPDATE tasks SET status = ?, updated_at = ? WHERE id = ?`,
		status, time.Now().UTC().Format(time.RFC3339Nano), id)
	if err != nil {
		return task.Task{}, fmt.Errorf("set task status: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return task.Task{}, fmt.Errorf("check task status: %w", err)
	}
	if affected == 0 {
		return task.Task{}, task.ErrNotFound
	}
	return r.Get(ctx, id)
}

func (r *TaskRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM tasks WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check deleted task: %w", err)
	}
	if affected == 0 {
		return task.ErrNotFound
	}
	return nil
}

type scanner interface {
	Scan(...any) error
}

func scanTask(s scanner) (task.Task, error) {
	var item task.Task
	var createdAt, updatedAt string
	if err := s.Scan(&item.ID, &item.Title, &item.Description, &item.Status, &item.Priority, &createdAt, &updatedAt); err != nil {
		return task.Task{}, err
	}
	var err error
	item.CreatedAt, err = time.Parse(time.RFC3339Nano, createdAt)
	if err != nil {
		return task.Task{}, err
	}
	item.UpdatedAt, err = time.Parse(time.RFC3339Nano, updatedAt)
	return item, err
}

func escapeLike(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	value = strings.ReplaceAll(value, `%`, `\%`)
	return strings.ReplaceAll(value, `_`, `\_`)
}
