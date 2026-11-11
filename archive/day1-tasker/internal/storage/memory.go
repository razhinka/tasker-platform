package storage

import (
	"errors"
	"sync"

	"tasker/internal/domain"
)

var ErrTaskNotFound = errors.New("task not found")

type MemoryStore struct {
	tasks  map[int64]domain.Task
	nextID int64
	mu     sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		tasks:  make(map[int64]domain.Task),
		nextID: 1,
	}
}

func (m *MemoryStore) Save(t domain.Task) domain.Task {
	m.mu.Lock()
	defer m.mu.Unlock()

	if t.ID == 0 {
		t.ID = m.nextID
		m.nextID++
	}
	m.tasks[t.ID] = t
	return t
}

func (m *MemoryStore) GetByID(id int64) (domain.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	task, ok := m.tasks[id]
	if !ok {
		return domain.Task{}, ErrTaskNotFound
	}
	return task, nil
}

func (m *MemoryStore) ListAll() []domain.Task {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]domain.Task, 0, len(m.tasks))
	for _, t := range m.tasks {
		result = append(result, t)
	}
	return result
}

func (m *MemoryStore) Delete(id int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.tasks[id]
	if !ok {
		return ErrTaskNotFound
	}
	delete(m.tasks, id)
	return nil
}
