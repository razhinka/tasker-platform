package storage_test

import (
	"context"
	"testing"
	"time"

	"tasker/internal/domain"
	"tasker/internal/storage"
)

func TestStartOverdueChecker(t *testing.T) {
	store := storage.NewMemoryStore()

	// Создаём просроченную задачу (дедлайн 1 час назад)
	overdueTask, _ := domain.NewTask(0, "Overdue 1", "desc", domain.PriorityHigh, time.Now().Add(-1*time.Hour))
	store.Save(*overdueTask)

	// Создаём актуальную задачу (дедлайн через 2 часа)
	normalTask, _ := domain.NewTask(0, "Normal", "desc", domain.PriorityLow, time.Now().Add(2*time.Hour))
	store.Save(*normalTask)

	// Создаём контекст с отменой
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запускаем чекер с частым интервалом (10 мс)
	out := storage.StartOverdueChecker(ctx, store, 10*time.Millisecond)

	// Ожидаем получение просроченной задачи из канала
	select {
	case task := <-out:
		if task.Title != "Overdue 1" {
			t.Fatalf("ожидали задачу 'Overdue 1', получили '%s'", task.Title)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("таймаут: чекер не прислал просроченную задачу")
	}

	// Отменяем контекст и проверяем, что канал закрылся (воркер завершился)
	cancel()

	// Из закрытого канала чтение возвращает нулевое значение и ok == false
	select {
	case _, ok := <-out:
		if ok {
			t.Fatal("канал должен был закрыться после отмены контекста")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("таймаут ожидания закрытия канала после отмены ctx")
	}
}
