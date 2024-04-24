-- name: ListMakes :many
SELECT name, description, acceleration, top_speed, color
FROM make;