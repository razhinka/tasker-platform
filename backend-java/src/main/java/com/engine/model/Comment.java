package com.engine.model;

import com.fasterxml.jackson.annotation.JsonIgnore;
import jakarta.persistence.*;
import java.time.LocalDateTime;
import java.util.Objects;

@Entity
@Table(name = "comments")
public class Comment {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @ManyToOne
    @JoinColumn(name = "author_id", nullable = false)
    private User author;

    @Column(name = "text", nullable = false)
    private String text;

    @Column(name = "created", nullable = false)
    private LocalDateTime created;

    @ManyToOne
    @JoinColumn(name = "task_id")
    @JsonIgnore
    private Task task;

    public Comment() {
    }

    public Comment(User author, String text, LocalDateTime created) {
        this.author = Objects.requireNonNull(author, "author must not be null");
        this.text = Objects.requireNonNull(text, "text must not be null");
        this.created = Objects.requireNonNull(created, "created must not be null");
    }

    public Long getId() { return id; }
    public void setId(Long id) { this.id = id; }

    public User getAuthor() { return author; }
    public void setAuthor(User author) { this.author = author; }

    public String getText() { return text; }
    public void setText(String text) { this.text = text; }

    public LocalDateTime getCreated() { return created; }
    public void setCreated(LocalDateTime created) { this.created = created; }

    public Task getTask() { return task; }
    public void setTask(Task task) { this.task = task; }

    public boolean isWrittenBy(User user) {
        return user != null && user.equals(author);
    }
}