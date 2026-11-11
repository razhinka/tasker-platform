package httpapi_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"tasker/internal/domain"
	"tasker/internal/httpapi"
	"tasker/internal/storage"
)

func setupTestServer() (http.Handler, *storage.MemoryStore) {
	store := storage.NewMemoryStore()
	handler := httpapi.NewHandler(store)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	return mux, store
}

func TestHandler_Endpoints(t *testing.T) {
	router, store := setupTestServer()

	t.Run("POST /tasks — создание задачи", func(t *testing.T) {
		body := `{"title": "Test Task", "description": "Desc", "priority": "HIGH"}`
		req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("ожидали статус 201 Created, получили %d. Ответ: %s", rec.Code, rec.Body.String())
		}

		var created domain.Task
		if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
			t.Fatalf("не удалось распарсить ответ: %v", err)
		}

		if created.ID == 0 || created.Title != "Test Task" {
			t.Fatalf("некорректная задача в ответе: %+v", created)
		}
	})

	t.Run("GET /tasks/{id} — получение задачи", func(t *testing.T) {
		task, _ := domain.NewTask(0, "Existing", "desc", domain.PriorityLow, time.Now())
		saved := store.Save(*task)

		idStr := strconv.FormatInt(saved.ID, 10)
		req := httptest.NewRequest(http.MethodGet, "/tasks/"+idStr, nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("ожидали 200 OK, получили %d", rec.Code)
		}

		var got domain.Task
		_ = json.Unmarshal(rec.Body.Bytes(), &got)
		if got.ID != saved.ID {
			t.Fatalf("ожидали ID %d, получили %d", saved.ID, got.ID)
		}
	})

	t.Run("DELETE /tasks/{id} — удаление задачи", func(t *testing.T) {
		task, _ := domain.NewTask(0, "To Delete", "desc", domain.PriorityLow, time.Now())
		saved := store.Save(*task)

		idStr := strconv.FormatInt(saved.ID, 10)
		req := httptest.NewRequest(http.MethodDelete, "/tasks/"+idStr, nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusNoContent {
			t.Fatalf("ожидали статус 204 No Content, получили %d", rec.Code)
		}

		// 2. Проверяем, что повторное удаление даёт 404 Not Found
		rec404 := httptest.NewRecorder()
		router.ServeHTTP(rec404, req)

		if rec404.Code != http.StatusNotFound {
			t.Fatalf("ожидали статус 404 Not Found при повторном удалении, получили %d", rec404.Code)
		}
	})
}
