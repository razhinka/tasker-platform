package com.engine.dto;

import com.engine.model.Priority;
import com.engine.model.Status;

import java.time.LocalDateTime;

public record TaskRequestDto(
        String title,
        String description,
        Status status,
        Priority priority,
        Long projectId,
        Long assigneeId,
        LocalDateTime deadline
) {}