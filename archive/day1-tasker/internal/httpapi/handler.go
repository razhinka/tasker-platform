package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"tasker/internal/domain"
	"tasker/internal/storage"
)

// TaskStore — интерфейс хранилища на стороне потребителя (consumer-side).
// Пакет httpapi сам диктует, какие методы хранилища ему нужны.
type TaskStore interface {
	Save(t domain.Task) domain.Task
	GetByID(id int64) (domain.Task, error)
	ListAll() []domain.Task
	Delete(id int64) error
}

type Handler struct {
	store TaskStore
}

func NewHandler(store TaskStore) *Handler {
	return &Handler{store: store}
}

// RegisterRoutes регистрирует маршруты в ServeMux
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /tasks", h.handleList)
	mux.HandleFunc("POST /tasks", h.handleCreate)
	mux.HandleFunc("GET /tasks/{id}", h.handleGetByID)
	mux.HandleFunc("DELETE /tasks/{id}", h.handleDelete)
}

// GET /tasks
func (h *Handler) handleList(w http.ResponseWriter, r *http.Request) {
	tasks := h.store.ListAll()
	writeJSON(w, http.StatusOK, tasks)
}

// GET /tasks/{id}
func (h *Handler) handleGetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid task id"})
		return
	}

	task, err := h.store.GetByID(id)
	if err != nil {
		if errors.Is(err, storage.ErrTaskNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "task not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}

	writeJSON(w, http.StatusOK, task)
}

// DTO для создания задачи
type CreateTaskRequest struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Priority    domain.Priority `json:"priority"`
	Deadline    time.Time       `json:"deadline"`
}

// POST /tasks
func (h *Handler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json body"})
		return
	}

	task, err := domain.NewTask(0, req.Title, req.Description, req.Priority, req.Deadline)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	saved := h.store.Save(*task)
	writeJSON(w, http.StatusCreated, saved)
}

// DELETE /tasks/{id}
func (h *Handler) handleDelete(w http.ResponseWriter, r *http.Request) {
	// TODO: Реализуй удаление задачи:
	// 1. Извлеки id через r.PathValue("id") и спарси через strconv.ParseInt.
	//    Если ошибка — верни 400 Bad Request: writeJSON(w, http.StatusBadRequest, ...)
	// 2. Вызови h.store.Delete(id).
	//    Если вернулась ошибка errors.Is(err, storage.ErrTaskNotFound) — верни 404:
	//    writeJSON(w, http.StatusNotFound, ...)
	// 3. Если всё успешно — верни статус 204 No Content:
	//    w.WriteHeader(http.StatusNoContent)

	id, err := strconv.ParseInt(
		r.PathValue("id"),
		10,
		64,
	)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, "really bad bro")
		return
	}

	err = h.store.Delete(id)
	if errors.Is(err, storage.ErrTaskNotFound) {
		writeJSON(w, http.StatusNotFound, "task not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Вспомогательная функция отправки JSON
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
