/**
  The task table must have indexes on the project_id
  because queries will most often perform JOINs on this index
 */
CREATE INDEX idx_tasks_project_id ON tasks (project_id);

CREATE INDEX inx_tasks_assignee_id ON tasks (assignee_id);

CREATE INDEX idx_tasks_deadline ON tasks (deadline);

CREATE INDEX idx_tasks_status_priority ON tasks (status, priority);
