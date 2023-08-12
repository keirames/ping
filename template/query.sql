-- name: GetMessagesByRoomID :many
SELECT * FROM messages m
WHERE m.room_id = $1
ORDER BY m.created_at ASC
LIMIT 10
OFFSET $2;

-- name: GetRooms :many
SELECT r.* FROM chat_rooms r
INNER JOIN users_and_chat_rooms uacr
ON uacr.room_id = r.id
WHERE uacr.user_id = $1
LIMIT 10
OFFSET $2;

-- name: IsMessageExist :one
SELECT 1 FROM messages m
WHERE m.id = $1
AND m.room_id = $2;

-- name: IsRoomExist :one
SELECT 1 FROM chat_rooms r
WHERE r.id = $1;

-- name: IsMemberOfRoom :one
SELECT 1 FROM users_and_chat_rooms uacr
WHERE uacr.user_id = $1
AND uacr.room_id = $2;

-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = $1 LIMIT 1;

-- name: CreateMessage :one
INSERT INTO messages (
  id, content, type, room_id, user_id
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: ListAuthors :many
SELECT * FROM authors
ORDER BY name;

-- name: CreateAuthor :one
INSERT INTO authors (
  name, bio
) VALUES (
  $1, $2
)
RETURNING *;

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = $1;