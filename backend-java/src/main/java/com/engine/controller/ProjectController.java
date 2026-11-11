package com.engine.controller;

import com.engine.dto.ProjectRequestDto;
import com.engine.model.Project;
import com.engine.model.Task;
import com.engine.service.ProjectService;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/projects")
public class ProjectController {

    private final ProjectService projectService;

    public ProjectController(ProjectService projectService) {
        this.projectService = projectService;
    }

    @GetMapping
    public List<Project> getAllProjects() {
        return projectService.getAllProjects();
    }

    @GetMapping("/{id}")
    public Project getProjectById(@PathVariable Long id) {
        return projectService.getProjectById(id);
    }

    @GetMapping("/{id}/tasks")
    public List<Task> getTasksByProjectId(@PathVariable Long id) {
        return projectService.getTasksByProjectId(id);
    }

    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    public Project createProject(@RequestBody ProjectRequestDto dto) {
        return projectService.createProject(dto);
    }

    @PutMapping("/{id}")
    public Project updateProject(@PathVariable Long id, @RequestBody ProjectRequestDto dto) {
        return projectService.updateProject(id, dto);
    }

    @DeleteMapping("/{id}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    public void deleteProject(@PathVariable Long id) {
        projectService.deleteProject(id);
    }
}