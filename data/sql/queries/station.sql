-- name: ListStations :many
SELECT * FROM station ORDER BY id;

-- name: GetStationById :one
SELECT * FROM station WHERE id = ? ORDER BY id;

-- name: GetStationByName :one
SELECT * FROM station
WHERE name = ? 
ORDER BY id;