package com.engine.repository;

import com.engine.model.Status;
import com.engine.model.Task;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

public interface TaskRepository extends JpaRepository<Task, Long> {

    List<Task> findByStatus(Status status);

    List<Task> findByProjectId(Long projectId);

    List<Task> findByAssigneeIdAndStatus(Long assigneeId, Status status);
}