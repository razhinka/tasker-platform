SELECT tasks.priority, COUNT(*)
FROM tasks
GROUP BY priority;

SELECT p.title, COUNT(*)
FROM tasks
         JOIN projects p on p.id = tasks.project_id
GROUP BY p.title
having count(*) > 2;
-- ERROR: function avg(timestamp without time zone) does not exist
SELECT status,
       TIMESTAMP 'epoch' + AVG(EXTRACT(EPOCH FROM deadline)) * INTERVAL '1 second' AS avg_deadline
FROM tasks
WHERE deadline IS NOT NULL
GROUP BY status;

SELECT projects.title AS project_title, COUNT(*)
FROM tasks
         JOIN projects ON tasks.project_id = projects.id
GROUP BY projects.title
ORDER BY COUNT(*);