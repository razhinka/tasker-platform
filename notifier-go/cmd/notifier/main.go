package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"tasker-notifier/internal/consumer"
	"tasker-notifier/internal/event"
)

func main() {
	log.Println("🚀 Запуск микросервиса tasker-notifier (Go + Kafka)...")

	// 1. Создаем контекст, слушающий сигналы ОС (Ctrl+C, SIGTERM)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 2. Определяем функцию обработки каждого входящего события
	eventHandler := func(ctx context.Context, e event.TaskEvent) error {
		assignee := "не назначен"
		if e.Assignee != nil {
			assignee = *e.Assignee
		}

		switch e.EventType {
		case "TASK_CREATED":
			log.Printf("🔔 [НОВАЯ ЗАДАЧА] #%d %q | Приоритет: %s | Исполнитель: %s\n",
				e.TaskID, e.Title, e.Priority, assignee)
		case "TASK_UPDATED":
			log.Printf("✏️  [ОБНОВЛЕНИЕ] Задача #%d %q переведена в статус %s\n",
				e.TaskID, e.Title, e.Status)
		case "TASK_DELETED":
			log.Printf("🗑️  [УДАЛЕНИЕ] Задача #%d %q была удалена\n",
				e.TaskID, e.Title)
		default:
			log.Printf("ℹ️  [СОБЫТИЕ] #%d: %s\n", e.TaskID, e.EventType)
		}
		return nil
	}

	// 3. Создаем Kafka Consumer
	brokers := []string{"localhost:9092"}
	topic := "task.events"
	groupID := "tasker-notifier-group"

	c := consumer.NewKafkaConsumer(brokers, topic, groupID, eventHandler)
	defer c.Close()

	// 4. Запускаем консьюмер в отдельной горутине
	errChan := make(chan error, 1)
	go func() {
		errChan <- c.Start(ctx)
	}()

	// 5. Ожидаем завершения: либо отмена по Ctrl+C, либо фатальная ошибка
	select {
	case <-ctx.Done():
		log.Println("Получен сигнал завершения. Останавливаем Kafka Consumer...")
	case err := <-errChan:
		if err != nil {
			log.Fatalf("Фатальная ошибка консьюмера: %v", err)
		}
	}

	log.Println("👋 Микросервис tasker-notifier корректно остановлен.")
}
