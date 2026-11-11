package com.engine.service;

import com.engine.dto.ProjectRequestDto;
import com.engine.model.Project;
import com.engine.model.Task;
import com.engine.repository.ProjectRepository;
import com.engine.repository.TaskRepository;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;

@Service
@Transactional(readOnly = true)
public class ProjectService {

    private final ProjectRepository projectRepository;
    private final TaskRepository taskRepository;

    public ProjectService(ProjectRepository projectRepository, TaskRepository taskRepository) {
        this.projectRepository = projectRepository;
        this.taskRepository = taskRepository;
    }

    public List<Project> getAllProjects() {
        return projectRepository.findAll();
    }

    public Project getProjectById(Long id) {
        return projectRepository.findById(id)
                .orElseThrow(() -> new IllegalArgumentException("Project not found with id: " + id));
    }

    public List<Task> getTasksByProjectId(Long projectId) {
        if (!projectRepository.existsById(projectId)) {
            throw new IllegalArgumentException("Project not found with id: " + projectId);
        }
        return taskRepository.findByProjectId(projectId);
    }

    @Transactional
    public Project createProject(ProjectRequestDto dto) {
        Project project = new Project(dto.title());
        return projectRepository.save(project);
    }

    @Transactional
    public Project updateProject(Long id, ProjectRequestDto dto) {
        Project project = getProjectById(id);
        project.setTitle(dto.title());
        return projectRepository.save(project);
    }

    @Transactional
    public void deleteProject(Long id) {
        if (!projectRepository.existsById(id)) {
            throw new IllegalArgumentException("Project not found with id: " + id);
        }
        projectRepository.deleteById(id);
    }
}