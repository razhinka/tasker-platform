package com.engine.controller;

import com.engine.dto.TaskRequestDto;
import com.engine.model.Comment;
import com.engine.model.Task;
import com.engine.repository.CommentRepository;
import com.engine.service.TaskService;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/tasks")
public class TaskController {

    private final TaskService taskService;
    private final CommentRepository commentRepository;

    public TaskController(TaskService taskService, CommentRepository commentRepository) {
        this.taskService = taskService;
        this.commentRepository = commentRepository;
    }

    @GetMapping
    public List<Task> getAllTasks() {
        return taskService.getAllTasks();
    }

    @GetMapping("/{id}")
    public Task getTaskById(@PathVariable Long id) {
        return taskService.getTaskById(id);
    }

    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    public Task createTask(@RequestBody TaskRequestDto dto) {
        return taskService.createTaskFromDto(dto);
    }

    @PutMapping("/{id}")
    public Task updateTask(@PathVariable Long id, @RequestBody TaskRequestDto dto) {
        return taskService.updateTaskFromDto(id, dto);
    }

    @DeleteMapping("/{id}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    public void deleteTask(@PathVariable Long id) {
        taskService.deleteTask(id);
    }

    @GetMapping("/{id}/comments")
    public List<Comment> getTaskComments(@PathVariable Long id) {
        return commentRepository.findByTaskId(id);
    }
}