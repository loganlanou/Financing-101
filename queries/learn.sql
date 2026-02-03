-- name: ListLearningModules :many
SELECT id, title, description, category, sort_order, created_at
FROM learning_modules
ORDER BY category, sort_order;

-- name: ListLearningModulesByCategory :many
SELECT id, title, description, category, sort_order, created_at
FROM learning_modules
WHERE category = sqlc.arg('category')
ORDER BY sort_order;

-- name: GetLearningModule :one
SELECT id, title, description, category, sort_order, created_at
FROM learning_modules
WHERE id = sqlc.arg('id');

-- name: GetLessonsByModule :many
SELECT id, module_id, title, content, summary, sort_order, created_at
FROM lessons
WHERE module_id = sqlc.arg('module_id')
ORDER BY sort_order;

-- name: GetLesson :one
SELECT id, module_id, title, content, summary, sort_order, created_at
FROM lessons
WHERE id = sqlc.arg('id');

-- name: CountLessonsByModule :one
SELECT COUNT(*) as count
FROM lessons
WHERE module_id = sqlc.arg('module_id');

-- name: ListGlossaryTerms :many
SELECT id, term, definition, category, created_at
FROM glossary_terms
ORDER BY term;

-- name: ListGlossaryTermsByCategory :many
SELECT id, term, definition, category, created_at
FROM glossary_terms
WHERE category = sqlc.arg('category')
ORDER BY term;

-- name: GetGlossaryTerm :one
SELECT id, term, definition, category, created_at
FROM glossary_terms
WHERE term = sqlc.arg('term');

-- name: SearchGlossaryTerms :many
SELECT id, term, definition, category, created_at
FROM glossary_terms
WHERE term LIKE '%' || sqlc.arg('query') || '%'
   OR definition LIKE '%' || sqlc.arg('query') || '%'
ORDER BY term;

-- name: GetTodaysLearningTip :one
SELECT id, title, content, category, learn_url, active_date, created_at
FROM learning_tips
WHERE active_date = DATE('now')
LIMIT 1;

-- name: GetRandomLearningTip :one
SELECT id, title, content, category, learn_url, active_date, created_at
FROM learning_tips
ORDER BY RANDOM()
LIMIT 1;

-- name: ListLearningTips :many
SELECT id, title, content, category, learn_url, active_date, created_at
FROM learning_tips
ORDER BY created_at DESC
LIMIT sqlc.arg('limit');

-- name: InsertLearningModule :exec
INSERT INTO learning_modules (id, title, description, category, sort_order, created_at)
VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
    title=excluded.title,
    description=excluded.description,
    category=excluded.category,
    sort_order=excluded.sort_order;

-- name: InsertLesson :exec
INSERT INTO lessons (id, module_id, title, content, summary, sort_order, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
    module_id=excluded.module_id,
    title=excluded.title,
    content=excluded.content,
    summary=excluded.summary,
    sort_order=excluded.sort_order;

-- name: InsertGlossaryTerm :exec
INSERT INTO glossary_terms (id, term, definition, category, created_at)
VALUES (?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
    term=excluded.term,
    definition=excluded.definition,
    category=excluded.category;

-- name: InsertLearningTip :exec
INSERT INTO learning_tips (id, title, content, category, learn_url, active_date, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
    title=excluded.title,
    content=excluded.content,
    category=excluded.category,
    learn_url=excluded.learn_url,
    active_date=excluded.active_date;
