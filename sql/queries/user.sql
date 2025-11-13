-- name: GetAllUsers :many
SELECT * FROM users;

-- name: CreateUser :exec
INSERT INTO users (first_name, last_name, email)
VALUES ($1, $2, $3);
