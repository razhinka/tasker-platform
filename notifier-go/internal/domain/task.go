package domain

import (
	"errors"
	"strings"
	"time"
)

type Status string

const (
	StatusTodo       Status = "TODO"
	StatusInProgress Status = "IN_PROGRESS"
	StatusDone       Status = "DONE"
)

type Priority string

const (
	PriorityLow    Priority = "LOW"
	PriorityMedium Priority = "MEDIUM"
	PriorityHigh   Priority = "HIGH"
)

type Task struct {
	ID          int64
	Title       string
	Description string
	Status      Status
	Priority    Priority
	CreatedAt   time.Time
	Deadline    time.Time
}

func NewTask(id int64, title, description string, priority Priority, deadline time.Time) (*Task, error) {
	if strings.TrimSpace(title) == "" {
		return nil, errors.New("title cannot be empty")
	}

	return &Task{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      StatusTodo,
		Priority:    priority,
		CreatedAt:   time.Now(),
		Deadline:    deadline,
	}, nil
}

func (t *Task) MarkDone() {
	t.Status = StatusDone
}

// IsOverdue проверяет, просрочена ли задача на текущий момент.
// Обрати внимание: получатель (t Task) по значению, так как метод ничего не мутирует.
func (t Task) IsOverdue() bool {
	if t.Status == StatusDone || t.Deadline.IsZero() {
		return false
	}
	return time.Now().After(t.Deadline)
}
