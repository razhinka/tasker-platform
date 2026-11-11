package storage

import (
	"context"
	"time"

	"tasker/internal/domain"
)

// StartOverdueChecker запускает фоновую горутину, которая периодически (каждые interval)
// опрашивает store, находит просроченные задачи и отправляет их в выходной канал.
// При отмене ctx горутина корректно завершается, закрывая выходной канал.
func StartOverdueChecker(ctx context.Context, store *MemoryStore, interval time.Duration) <-chan domain.Task {
	out := make(chan domain.Task)

	go func() {
		defer close(out)

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return

			case <-ticker.C:
				tt := store.ListAll()
				for _, t := range tt {
					if t.IsOverdue() {
						select {
						case out <- t:
						case <-ctx.Done():
							return
						}
					}
				}
			}
		}
	}()

	return out
}
