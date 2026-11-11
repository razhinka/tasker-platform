package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tasker/internal/httpapi"
	"tasker/internal/storage"
)

func main() {
	// 1. Инициализируем хранилище
	store := storage.NewMemoryStore()

	// 2. Создаем контекст, который слушает сигналы операционной системы (SIGINT / Ctrl+C, SIGTERM)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 3. Запускаем фоновый чекер просроченных задач (проверяет каждые 3 секунды)
	overdueChan := storage.StartOverdueChecker(ctx, store, 3*time.Second)
	go func() {
		for task := range overdueChan {
			log.Printf("[ВНИМАНИЕ] Просрочена задача! ID=%d, Title=%q\n", task.ID, task.Title)
		}
	}()

	// 4. Настраиваем маршрутизацию
	mux := http.NewServeMux()
	handler := httpapi.NewHandler(store)
	handler.RegisterRoutes(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// 5. Запускаем HTTP-сервер в отдельной горутине, чтобы он не блокировал main
	go func() {
		log.Println("🚀 Сервер запущен на http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Ошибка HTTP-сервера: %v", err)
		}
	}()

	// 6. Блокируем main горутину до тех пор, пока пользователь не нажмет Ctrl+C
	<-ctx.Done()
	log.Println("\n🛑 Получен сигнал остановки (Ctrl+C). Начинаем Graceful Shutdown...")

	// 7. Даем серверу до 5 секунд, чтобы завершить активные запросы клиентов
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Принудительная остановка сервера: %v", err)
	}

	log.Println("👋 Сервер успешно и безопасно остановлен.")
}
