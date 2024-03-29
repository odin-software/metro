-- name: GetAllTrains :many
SELECT id, name, x, y, z FROM train 
ORDER BY id;

-- name: GetTrainById :one
SELECT id, name, x, y, z FROM train
WHERE id = ?
LIMIT 1;

-- name: GetTrainByName :one
SELECT id, name, x, y, z FROM train
WHERE name = ?
LIMIT 1;

-- name: CreateTrain :one
INSERT INTO train (name, x, y, z, currentId, makeId, lineId)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING id;

-- name: UpdateTrain :one
UPDATE train
SET name = ?, x = ?, y = ?, z = ?, currentId = ?, makeId = ?, lineId = ?
WHERE id = ?
RETURNING id;

-- name: DeleteTrain :exec
DELETE FROM train 
WHERE id = ?;

-- name: GetTrainLine :one
SELECT lineId FROM station_line
WHERE stationId = ?;

-- name: GetTrainMake :one
SELECT makeId FROM train
WHERE id = ?;

-- name: GetTrainNext :one
SELECT nextId FROM train
WHERE id = ?;

-- name: SetTrainNext :exec
UPDATE train
SET nextId = ?
WHERE id = ?;