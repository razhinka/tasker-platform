SELECT tasks.*
FROM tasks
         JOIN projects ON tasks.project_id = projects.id
WHERE projects.title = 'Java Course';

SELECT tasks.*, users.username, projects.title
FROM tasks
         JOIN users ON tasks.assignee_id = users.id
         JOIN projects ON tasks.project_id = projects.id
WHERE username = 'alice';

SELECT tasks.*, t.tag
FROM tasks
         JOIN task_tags tt on tasks.id = tt.task_id
         JOIN tags t on tt.tag_id = t.id;

SELECT tasks.*
FROM tasks
WHERE assignee_id IS NULL;

SELECT p.title, COUNT(*)
FROM tasks
         JOIN projects p on tasks.project_id = p.id
GROUP BY p.title;

