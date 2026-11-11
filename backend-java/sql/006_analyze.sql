INSERT INTO tasks (title, deadline, priority, status)
SELECT 'Task #' || i                                                AS title,
       NOW() + (random() * 10 || ' days')::interval                           AS deadline,
       (ARRAY ['LOW', 'MEDIUM', 'HIGH'])[random(1, 3)]              AS priority,
       (ARRAY ['DONE', 'IN_PROGRESS', 'NOT_STARTED'])[random(1, 3)] AS status
FROM generate_series(1, 100) AS i;

EXPLAIN ANALYZE
SELECT *
FROM tasks
WHERE deadline = '2026-09-14 08:49:59.855372';

EXPLAIN ANALYZE
SELECT *
FROM tasks
WHERE status = 'IN_PROGRESS'
  AND priority = 'HIGH';

EXPLAIN ANALYZE
SELECT *
FROM tasks
WHERE priority = 'HIGH';

CREATE INDEX idx_tasks_deadline ON tasks (deadline);

CREATE INDEX idx_tasks_status_priority ON tasks (status, priority);