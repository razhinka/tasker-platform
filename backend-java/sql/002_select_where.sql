-- all projects
SELECT * FROM projects;
-- all tasks
SELECT * from tasks;

-- status is in progress
SELECT * FROM tasks
WHERE status = 'IN_PROGRESS';
-- status is done
SELECT * FROM tasks
WHERE status = 'DONE';
-- status is not started
SELECT * FROM tasks
WHERE status = 'NOT_STARTED';
-- All tasks with high prio and ordered by deadline
SELECT * FROM tasks
WHERE priority = 'HIGH'
ORDER BY deadline;
-- First 5 on the second page
SELECT * FROM tasks
LIMIT 5 OFFSET 3;
-- deadline in the next 7 days
SELECT * FROM tasks
WHERE (deadline is not null) and (deadline < CURRENT_TIMESTAMP + INTERVAL '7 days') ;
-- projects that contains java in title
SELECT * FROM projects
WHERE title ILIKE '%java%';
-- only alice's tasks
SELECT tasks.* FROM tasks
JOIN users ON tasks.assignee_id = users.id
WHERE username = 'alice';
-- counting tasks by status
SELECT status, count(*) FROM tasks
GROUP BY status;