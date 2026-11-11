-- =======================================================
-- Часть 1: Проверка ROLLBACK (Атомарность — откат изменений)
-- =======================================================

-- Шаг 1. Смотрим исходный статус задачи с id = 1
SELECT id, title, status
FROM tasks
WHERE id = 1;

-- Шаг 2. Открываем транзакцию
BEGIN;

-- Шаг 3. Меняем статус на 'DONE' и добавляем тестовый комментарий
UPDATE tasks
SET status = 'DONE'
WHERE id = 1;

INSERT INTO comments (task_id, author_id, text)
VALUES (1, 1, 'Тестовый комментарий перед отменой');

-- Шаг 4. Проверяем внутри транзакции: здесь статус уже 'DONE' и комментарий виден
SELECT id, title, status
FROM tasks
WHERE id = 1;
SELECT *
FROM comments
WHERE task_id = 1;

-- Шаг 5. Отменяем транзакцию целиком
ROLLBACK;

-- Шаг 6. Проверяем после отката:
-- Статус вернулся в исходное состояние, а комментария в таблице нет!
SELECT id, title, status
FROM tasks
WHERE id = 1;
SELECT *
FROM comments
WHERE task_id = 1;


-- =======================================================
-- Часть 2: Проверка COMMIT (Фиксация изменений)
-- =======================================================

BEGIN;

UPDATE tasks
SET status = 'DONE'
WHERE id = 1;

-- Фиксируем изменения в базе на постоянной основе
COMMIT;

-- Проверяем: статус остался 'DONE'
SELECT id, title, status
FROM tasks
WHERE id = 1;