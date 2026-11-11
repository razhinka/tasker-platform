package storage_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"tasker/internal/domain"
	"tasker/internal/storage"
)

// Table-Driven Test (Табличный тест) — главный стандарт тестирования в Go.
// Вместо @ParameterizedTest в JUnit 5 создаётся срез анонимных структур с тест-кейсами.
func TestMemoryStore_GetByID(t *testing.T) {
	// Подготовка хранилища с тестовыми данными
	store := storage.NewMemoryStore()
	task, _ := domain.NewTask(0, "Buy milk", "2% fat", domain.PriorityLow, time.Now())
	savedTask := store.Save(*task)

	// Описываем тест-кейсы
	tests := []struct {
		name        string
		id          int64
		wantErr     error
		expectFound bool
	}{
		{
			name:        "существующая задача",
			id:          savedTask.ID,
			wantErr:     nil,
			expectFound: true,
		},
		{
			name:        "несуществующая задача",
			id:          999,
			wantErr:     storage.ErrTaskNotFound,
			expectFound: false,
		},
	}

	for _, tt := range tests {
		// t.Run создаёт изолированный подтест с именем tt.name
		t.Run(tt.name, func(t *testing.T) {
			got, err := store.GetByID(tt.id)

			if tt.wantErr != nil {
				// Проверяем ошибку через errors.Is
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("ожидали ошибку %v, получили %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("не ожидали ошибку, но получили: %v", err)
			}

			if got.ID != savedTask.ID {
				t.Errorf("ожидали ID %d, получили %d", savedTask.ID, got.ID)
			}
		})
	}
}

// TODO: Напиши по аналогии табличный тест для Delete:
// 1. Успешное удаление существующей задачи
// 2. Попытка удаления несуществующей (ожидаем storage.ErrTaskNotFound)
// 3. Проверка, что после удаления задача больше не находится через GetByID
func TestMemoryStore_Delete(t *testing.T) {
	store := storage.NewMemoryStore()
	task, _ := domain.NewTask(0, "Buy milk", "2% fat", domain.PriorityLow, time.Now())
	savedTask := store.Save(*task)

	tests := []struct {
		name    string
		id      int64
		wantErr error
	}{
		{
			name:    "успешное удаление",
			id:      savedTask.ID,
			wantErr: nil,
		},
		{
			name:    "ошибка из-за несуществующего id",
			id:      999,
			wantErr: storage.ErrTaskNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.Delete(tt.id)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("ожидали ошибку %v, получили %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("не ожидали ошибку, но получили: %v", err)
			}
		})
	}
}

func TestMemoryStore_Concurrent(t *testing.T) {
	store := storage.NewMemoryStore()
	const workers = 50 // запускаем 50 параллельных горутин

	var wg sync.WaitGroup
	wg.Add(workers * 2) // 50 читателей + 50 писателей = 100 горутин

	// 50 горутин-писателей
	for i := 0; i < workers; i++ {
		go func(id int) {
			defer wg.Done()
			task, _ := domain.NewTask(0, "Title", "Desc", domain.PriorityLow, time.Now())
			store.Save(*task)
		}(i)
	}

	// 50 горутин-читателей
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			_ = store.ListAll()
		}()
	}

	// Ждем, пока все 100 горутин отработают
	wg.Wait()

	// Проверяем, что сохранились ровно 50 задач
	tasks := store.ListAll()
	if len(tasks) != workers {
		t.Fatalf("ожидали %d задач, но получили %d", workers, len(tasks))
	}
}
