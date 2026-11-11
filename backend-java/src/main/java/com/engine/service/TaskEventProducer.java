package com.engine.service;

import com.engine.dto.TaskEvent;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Service;

@Service
public class TaskEventProducer {

    private static final Logger log = LoggerFactory.getLogger(TaskEventProducer.class);
    public static final String TOPIC = "task.events";

    private final KafkaTemplate<String, Object> kafkaTemplate;

    public TaskEventProducer(KafkaTemplate<String, Object> kafkaTemplate) {
        this.kafkaTemplate = kafkaTemplate;
    }

    public void sendTaskEvent(TaskEvent event) {
        String key = event.taskId() != null ? String.valueOf(event.taskId()) : event.eventId();
        log.info("Publishing TaskEvent to Kafka topic [{}]: type={}, taskId={}", TOPIC, event.eventType(), event.taskId());
        kafkaTemplate.send(TOPIC, key, event);
    }
}
