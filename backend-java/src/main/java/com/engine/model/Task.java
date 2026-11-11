package com.engine.model;

import com.fasterxml.jackson.annotation.JsonIgnore;
import jakarta.persistence.*;

import java.time.LocalDateTime;
import java.util.*;

@Entity
@Table(name = "tasks")
public class Task {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(name = "title", nullable = false)
    private String title;

    @Column(name = "description")
    private String description;

    @Enumerated(EnumType.STRING)
    @Column(name = "status")
    private Status status = Status.NOT_STARTED;

    @Enumerated(EnumType.STRING)
    @Column(name = "priority")
    private Priority priority = Priority.LOW;

    @ManyToOne
    @JoinColumn(name = "project_id")
    private Project project;

    @ManyToOne
    @JoinColumn(name = "assignee_id")
    private User assignee;

    @ManyToMany
    @JoinTable(
            name = "task_tags",
            joinColumns = @JoinColumn(name = "task_id"),
            inverseJoinColumns = @JoinColumn(name = "tag_id")
    )
    private Set<Tag> tags = new HashSet<>();

    @OneToMany(mappedBy = "task", cascade = CascadeType.ALL, orphanRemoval = true)
    private List<Comment> comments = new ArrayList<>();

    @Column(name = "deadline")
    private LocalDateTime deadline;

    public Task() {
    }

    private Task(Builder builder) {
        this.id = builder.id;
        this.title = builder.title;
        this.description = builder.description;
        this.project = builder.project;
        this.assignee = builder.assignee;
        this.status = builder.status;
        this.priority = builder.priority;
        this.tags = new HashSet<>(builder.tags);
        this.comments = new ArrayList<>(builder.comments);
        this.deadline = builder.deadline;
    }

    public static class Builder {
        private Long id;
        private String title;
        private Project project;
        private User assignee;
        private String description = "";
        private Status status = Status.NOT_STARTED;
        private Priority priority = Priority.LOW;
        private Set<Tag> tags = new HashSet<>();
        private List<Comment> comments = new ArrayList<>();
        private LocalDateTime deadline;

        public Builder(Long id, String title, Project project, User assignee) {
            this.id = id;
            this.title = Objects.requireNonNull(title, "title must not be null");
            this.project = Objects.requireNonNull(project, "project must not be null");
            this.assignee = Objects.requireNonNull(assignee, "assignee must not be null");
        }

        public Builder(Task task) {
            this.id = task.id;
            this.title = task.title;
            this.description = task.description;
            this.project = task.project;
            this.assignee = task.assignee;
            this.status = task.status;
            this.priority = task.priority;
            this.tags = new HashSet<>(task.tags);
            this.comments = new ArrayList<>(task.comments);
            this.deadline = task.deadline;
        }

        public Builder id(Long id) {
            this.id = id;
            return this;
        }

        public Builder title(String title) {
            this.title = title;
            return this;
        }

        public Builder project(Project project) {
            this.project = project;
            return this;
        }

        public Builder assignee(User assignee) {
            this.assignee = assignee;
            return this;
        }

        public Builder description(String description) {
            this.description = description;
            return this;
        }

        public Builder status(Status status) {
            this.status = status;
            return this;
        }

        public Builder priority(Priority priority) {
            this.priority = priority;
            return this;
        }

        public Builder tags(Set<Tag> tags) {
            this.tags = tags;
            return this;
        }

        public Builder comments(List<Comment> comments) {
            this.comments = comments;
            return this;
        }

        public Builder deadline(LocalDateTime deadline) {
            this.deadline = deadline;
            return this;
        }

        public Task build() {
            return new Task(this);
        }
    }

    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public Status getStatus() {
        return status;
    }

    public void setStatus(Status status) {
        this.status = status;
    }

    public Priority getPriority() {
        return priority;
    }

    public void setPriority(Priority priority) {
        this.priority = priority;
        return;
    }

    public Project getProject() {
        return project;
    }

    public void setProject(Project project) {
        this.project = project;
    }

    public User getAssignee() {
        return assignee;
    }

    public void setAssignee(User assignee) {
        this.assignee = assignee;
    }

    public Set<Tag> getTags() {
        return tags;
    }

    public void setTags(Set<Tag> tags) {
        this.tags = tags;
    }

    public List<Comment> getComments() {
        return comments;
    }

    public void setComments(List<Comment> comments) {
        this.comments = comments;
    }

    public Optional<LocalDateTime> getDeadline() {
        return Optional.ofNullable(deadline);
    }

    public void setDeadline(LocalDateTime deadline) {
        this.deadline = deadline;
    }

    @Override
    public String toString() {
        String assigneeName = (assignee != null) ? assignee.getUsername() : "Unassigned";
        String projectTitle = (project != null) ? project.getTitle() : "No Project";
        return "Task " + id + " \"" + title + "\"\nAssignee: " + assigneeName
                + "\nProject: " + projectTitle + "\nDescription: " + description
                + "\nStatus: " + status + "\nPriority: " + priority + "\nDeadline: " + deadline + "\n";
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (!(o instanceof Task other)) return false;
        return id != null && id.equals(other.id);
    }

    @Override
    public int hashCode() {
        return Objects.hash(id);
    }

    @JsonIgnore
    public boolean isDone() {
        return Status.DONE.equals(status);
    }
}