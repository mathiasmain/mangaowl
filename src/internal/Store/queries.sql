-- name: GetTODOByID :one
SELECT *
FROM Todo
WHERE id = $1
LIMIT 1;

-- name: GetTODOS :many
SELECT *
FROM Todo
ORDER BY createdAt;

-- name: CreateTODO :exec
INSERT INTO Todo (body)
VALUES ($1);

-- name: UpdateTODO :exec
UPDATE Todo
set completedAt = $2
WHERE id = $1;

-- name: DeleteTODO :exec
DELETE FROM Todo
WHERE id = $1;