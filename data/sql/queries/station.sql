-- name: ListStations :many
SELECT id, name, x, y, z FROM station ORDER BY id;

-- name: GetStationById :one
SELECT id, name, x, y, z FROM station WHERE id = ? ORDER BY id;

-- name: GetStationByName :one
SELECT id, name, x, y, z FROM station
WHERE name = ? 
ORDER BY id;

-- name: CreateStation :one
INSERT INTO station (name, x, y, z) VALUES (?, ?, ?, ?);

-- name: UpdateStation :one
UPDATE station
SET name = ?, x = ?, y = ?, z = ?
WHERE id = ?;

-- name: DeleteStation :one
DELETE FROM station WHERE id = ?;

-- name: TrainsAtStation :many
SELECT id, name FROM train WHERE currentId = ? ORDER BY id;

-- name: TrainsToStation :many
SELECT id, name FROM train WHERE nextId = ? ORDER BY id;

-- name: GetStationLines :many
SELECT id FROM station_line WHERE stationId = ? ORDER BY id;