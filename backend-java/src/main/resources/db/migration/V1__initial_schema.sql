/**
  Author: Razhinka
  Create date: 14.09.2026
 */

CREATE TABLE users
(
    id       BIGSERIAL PRIMARY KEY,
    username VARCHAR(31) NOT NULL,
    CONSTRAINT unq_username UNIQUE (username)
);

create table projects
(
    id    BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL
);

CREATE TABLE tasks
(
    id          BIGSERIAL PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    status      VARCHAR(32)  NOT NULL DEFAULT 'NOT STARTED',
    priority    VARCHAR(16)  NOT NULL DEFAULT 'LOW',
    deadline    TIMESTAMP,
    description TEXT,
    project_id  BIGINT,
    assignee_id BIGINT REFERENCES users (id),
    CONSTRAINT fk_tasks_projects FOREIGN KEY (project_id) REFERENCES projects (id),
    CONSTRAINT fk_tasks_assignee FOREIGN KEY (assignee_id) REFERENCES users (id),
    CONSTRAINT valid_status CHECK (status IN ('NOT_STARTED', 'IN_PROGRESS', 'DONE')),
    CONSTRAINT valid_priority CHECK (priority IN ('LOW', 'MEDIUM', 'HIGH'))
);

CREATE TABLE comments
(
    id        BIGSERIAL PRIMARY KEY,
    task_id   BIGINT,
    author_id BIGINT REFERENCES users (id) NOT NULL,
    text      TEXT                         NOT NULL,
    created   timestamp                    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (task_id) REFERENCES tasks (id) ON DELETE CASCADE
);

CREATE TABLE tags
(
    id  BIGSERIAL PRIMARY KEY,
    tag VARCHAR(31) NOT NULL UNIQUE
);

CREATE TABLE task_tags
(
    task_id BIGINT,
    tag_id  BIGINT,
    CONSTRAINT pk_task_tags PRIMARY KEY (task_id, tag_id)
)