package consumer_test

import (
	"encoding/json"
	"testing"

	"tasker-notifier/internal/event"
)

func TestTaskEvent_JSONDeserialization(t *testing.T) {
	// Реальный JSON, который сгенерировал наш Spring Boot!
	realKafkaJSON := []byte(`{"eventId":"aa77f197-8fe8-47cd-a237-c247ba3d15fe","eventType":"TASK_CREATED","taskId":103,"title":"Event-Driven Архитектура","description":"Связка Java и Go","status":"NOT_STARTED","priority":"HIGH","assignee":"user_1","deadline":null,"timestamp":1794386396.870369100}`)

	var evt event.TaskEvent
	err := json.Unmarshal(realKafkaJSON, &evt)
	if err != nil {
		t.Fatalf("ошибка парсинга реального JSON от Spring Boot: %v", err)
	}

	if evt.TaskID != 103 {
		t.Errorf("ожидали taskId 103, получили %d", evt.TaskID)
	}
	if evt.EventType != "TASK_CREATED" {
		t.Errorf("ожидали eventType TASK_CREATED, получили %s", evt.EventType)
	}
	if evt.Assignee == nil || *evt.Assignee != "user_1" {
		t.Errorf("ожидали assignee 'user_1'")
	}
	if evt.Timestamp.IsZero() {
		t.Errorf("timestamp не должен быть пустым")
	}
}
