-- 1. Справочник пользователей (10 пользователей)
-- С ограничением unq_username
INSERT INTO users (username)
SELECT 'user_' || i
FROM generate_series(1, 10) AS i
ON CONFLICT DO NOTHING;

-- 2. Справочник проектов (5 проектов)
INSERT INTO projects (title)
SELECT 'Project ' || i
FROM generate_series(1, 5) AS i;

-- 3. Справочник тегов (10 тегов)
-- С ограничением UNIQUE (tag)
INSERT INTO tags (tag)
VALUES ('backend'), ('frontend'), ('bug'), ('feature'), ('database'),
       ('refactoring'), ('testing'), ('devops'), ('urgent'), ('security')
ON CONFLICT (tag) DO NOTHING;

-- 4. Задачи (100 штук) со связью на users и projects
INSERT INTO tasks (title, deadline, priority, status, project_id, assignee_id, description)
SELECT
    'Task #' || i                                                                    AS title,
    NOW() + (random() * 14 || ' days')::interval                                     AS deadline,
    (ARRAY['LOW', 'MEDIUM', 'HIGH'])[1 + floor(random() * 3)::int]                   AS priority,
    (ARRAY['NOT_STARTED', 'IN_PROGRESS', 'DONE'])[1 + floor(random() * 3)::int]     AS status,
    1 + floor(random() * 5)::int                                                     AS project_id,   -- id от 1 до 5
    1 + floor(random() * 10)::int                                                    AS assignee_id,  -- id от 1 до 10
    'Description for task #' || i                                                    AS description
FROM generate_series(1, 100) AS i;

-- 5. Связи задач с тегами (task_tags: по 1-3 случайных тега на задачу)
-- DISTINCT предотвращает дублирование первичного ключа pk_task_tags (task_id, tag_id)
INSERT INTO task_tags (task_id, tag_id)
SELECT DISTINCT
    t.id                                AS task_id,
    1 + floor(random() * 10)::bigint    AS tag_id
FROM tasks t
         CROSS JOIN generate_series(1, 2)
ON CONFLICT DO NOTHING;

-- 6. Комментарии к задачам (по 1-2 комментария на часть задач)
INSERT INTO comments (task_id, author_id, text, created)
SELECT
    t.id                                            AS task_id,
    1 + floor(random() * 10)::bigint                AS author_id,
    'Auto-generated comment for task #' || t.id     AS text,
    NOW() - (random() * 5 || ' days')::interval     AS created
FROM tasks t
WHERE random() > 0.3; -- добавляем комментарии примерно для 70% задач