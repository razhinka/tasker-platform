package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"tasker-notifier/internal/event"

	"github.com/segmentio/kafka-go"
)

type EventHandler func(ctx context.Context, e event.TaskEvent) error // почему мы пишем тим а потом просто функция без названия и тела?

type KafkaConsumer struct {
	reader  *kafka.Reader
	handler EventHandler
}

func NewKafkaConsumer(brokers []string, topic, groupID string, handler EventHandler) *KafkaConsumer { // ну тут всё понятно это конструктор
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return &KafkaConsumer{
		reader:  reader,
		handler: handler,
	}
}

func (c *KafkaConsumer) Start(ctx context.Context) error { // почему после func в скобках c *KafkaConsumer и только потом название функции?
	log.Println("Kafka Consumer успешно запущен, ожидание событий...")

	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				log.Println("Контекст отменён, завершаем Kafka Consumer...")
				return nil
			}
			log.Printf("Ошибка чтения из Kafka: %v", err)
			continue
		}

		var evt event.TaskEvent
		err = json.Unmarshal(msg.Value, &evt)
		if err != nil {
			log.Printf("Ошибка десериализации json сообщения из Kafka: %v", err)
			continue
		}
		c.handler(ctx, evt)
	}
}

func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}
