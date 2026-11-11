package com.engine.dto;

import java.time.Instant;
import java.time.LocalDateTime;

public record TaskEvent(
    String eventId,
    String eventType, // "TASK_CREATED", "TASK_UPDATED", "TASK_DELETED"
    Long taskId,
    String title,
    String description,
    String status,
    String priority,
    String assignee,
    LocalDateTime deadline,
    Instant timestamp
) {}