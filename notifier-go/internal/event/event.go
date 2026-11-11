package event

import (
	"encoding/json"
	"strings"
	"time"
)

// FlexibleTime парсит как числовой timestamp от Jackson (1794386396.87),
// так и строки RFC3339 ("2026-11-11T07:50:00Z") и LocalDateTime ("2026-11-20T18:00:00").
type FlexibleTime struct {
	time.Time
}

func (ft *FlexibleTime) UnmarshalJSON(b []byte) error {
	raw := string(b)
	if raw == "null" || raw == `""` || raw == "" {
		return nil
	}

	// 1. Если Jackson передал число epoch seconds (например 1794386396.870369100)
	var sec float64
	if err := json.Unmarshal(b, &sec); err == nil {
		secInt := int64(sec)
		nsec := int64((sec - float64(secInt)) * 1e9)
		ft.Time = time.Unix(secInt, nsec)
		return nil
	}

	// 2. Если строка: снимаем кавычки
	s := strings.Trim(raw, `"`)

	// 2.1. Пробуем стандартный RFC3339
	t, err := time.Parse(time.RFC3339, s)
	if err == nil {
		ft.Time = t
		return nil
	}

	// 2.2. Пробуем формат Java LocalDateTime ("2006-01-02T15:04:05")
	t, err = time.Parse("2006-01-02T15:04:05", s)
	if err == nil {
		ft.Time = t
		return nil
	}

	return err
}

// TaskEvent полностью совпадает по контракту с Java record (com.engine.dto.TaskEvent).
type TaskEvent struct {
	EventID     string        `json:"eventId"`
	EventType   string        `json:"eventType"` // "TASK_CREATED", "TASK_UPDATED", "TASK_DELETED"
	TaskID      int64         `json:"taskId"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Status      string        `json:"status"`
	Priority    string        `json:"priority"`
	Assignee    *string       `json:"assignee"`
	Deadline    *FlexibleTime `json:"deadline"`
	Timestamp   FlexibleTime  `json:"timestamp"`
}
