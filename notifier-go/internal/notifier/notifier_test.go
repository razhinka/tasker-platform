package notifier_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"tasker-notifier/internal/domain"
	"tasker-notifier/internal/notifier"
)

func TestNotifier_WorkerPool(t *testing.T) {
	t.Run("пул воркеров параллельно обрабатывает задачи", func(t *testing.T) {
		var processedCount int64

		// Имитируем отправку, которая занимает 40 мс
		sender := func(ctx context.Context, task domain.Task) error {
			time.Sleep(40 * time.Millisecond)
			atomic.AddInt64(&processedCount, 1)
			return nil
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Создаем сервис с очередью на 10 элементов
		svc := notifier.NewService(10, sender)

		// Запускаем пул из 3 воркеров!
		// Если бы был 1 воркер, 3 задачи заняли бы: 3 * 40мс = 120мс.
		// А 3 параллельных воркера разберут их одновременно за ~40мс!
		svc.StartWorkerPool(ctx, 3, 100*time.Millisecond)

		start := time.Now()
		for i := 1; i <= 3; i++ {
			task, _ := domain.NewTask(int64(i), "T", "desc", domain.PriorityLow, time.Now())
			_ = svc.Send(*task)
		}

		// Даём 70 мс на обработку всех 3 задач параллельно
		time.Sleep(70 * time.Millisecond)

		elapsed := time.Since(start)
		if atomic.LoadInt64(&processedCount) != 3 {
			t.Fatalf("ожидали 3 обработанные задачи, получили %d", processedCount)
		}

		// Проверяем, что заняло меньше 100мс (значит они реально работали параллельно!)
		if elapsed > 100*time.Millisecond {
			t.Fatalf("задачи выполнялись слишком долго (%v), похоже они выполнялись последовательно, а не параллельно", elapsed)
		}

		// Проверяем корректную остановку пула
		cancel()   // отменяем контекст
		svc.Wait() // все воркеры должны чисто выйти
	})

	t.Run("таймаут контекста на долгую задачу", func(t *testing.T) {
		var timedOut int64

		// Задача зависает на 200мс, но таймаут задан 30мс
		sender := func(ctx context.Context, task domain.Task) error {
			select {
			case <-time.After(200 * time.Millisecond):
				return nil
			case <-ctx.Done():
				atomic.AddInt64(&timedOut, 1)
				return ctx.Err()
			}
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		svc := notifier.NewService(5, sender)
		svc.StartWorkerPool(ctx, 1, 30*time.Millisecond) // таймаут всего 30мс

		task, _ := domain.NewTask(100, "Slow", "desc", domain.PriorityHigh, time.Now())
		_ = svc.Send(*task)

		time.Sleep(60 * time.Millisecond)

		if atomic.LoadInt64(&timedOut) != 1 {
			t.Fatalf("таймаут не отменил задачу вовремя")
		}

		cancel()
		svc.Wait()
	})

	t.Run("переполнение очереди возвращает ErrQueueFull", func(t *testing.T) {
		svc := notifier.NewService(1, nil)
		task, _ := domain.NewTask(99, "Task", "desc", domain.PriorityLow, time.Now())

		if err := svc.Send(*task); err != nil {
			t.Fatalf("первая задача должна была войти: %v", err)
		}

		err := svc.Send(*task)
		if !errors.Is(err, notifier.ErrQueueFull) {
			t.Fatalf("ожидали ErrQueueFull, получили %v", err)
		}
	})
}
