package notifier

import (
	"context"
	"errors"
	"sync"
	"time"

	"tasker-notifier/internal/domain"
)

var ErrQueueFull = errors.New("notification queue is full")

// Sender — функция отправки (например, HTTP запрос в Telegram/Webhook).
type Sender func(ctx context.Context, t domain.Task) error

type Service struct {
	jobs   chan domain.Task
	sender Sender
	wg     sync.WaitGroup
}

func NewService(queueSize int, sender Sender) *Service {
	return &Service{
		jobs:   make(chan domain.Task, queueSize),
		sender: sender,
	}
}

// Send кладёт задачу в очередь без блокировки.
func (s *Service) Send(t domain.Task) error {
	select {
	case s.jobs <- t:
		return nil
	default:
		// Если буфер канала полон, select мгновенно выполняет ветку default
		return ErrQueueFull
	}
}

// StartWorkerPool запускает пул из workersCount параллельных горутин-воркеров.
func (s *Service) StartWorkerPool(ctx context.Context, workersCount int, taskTimeout time.Duration) {
	for i := range workersCount {
		s.wg.Add(1)

		go func(workerID int) {
			defer s.wg.Done()

			for {
				select {
				case <-ctx.Done():
					return

				case task := <-s.jobs:
					taskCtx, cancel := context.WithTimeout(ctx, taskTimeout)
					s.sender(taskCtx, task)
					cancel()
				}
			}
		}(i)
	}
}

// Wait дожидается завершения всех воркеров пула
func (s *Service) Wait() {
	s.wg.Wait()
}
