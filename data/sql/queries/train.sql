-- name: GetAllTrains :many
SELECT id, name, x, y, z FROM train
ORDER BY id;

-- name: GetAllTrainsFull :many
SELECT
	tr.id,
	tr.name,
	tr.x,
	tr.y,
	tr.z,
	tr.forward,
	mk.color,
	st.id as stationId,
	ln.name as lineName,
	mk.name as makeName
FROM train tr
JOIN line ln ON tr.lineId = ln.id
JOIN make mk ON tr.makeId = mk.id
JOIN station st ON tr.currentId = st.id;

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
SET name = ?, x = ?, y = ?, z = ?, currentId = ?, nextId = ?
WHERE name = ?
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

-- name: ChangeTrainToLine :exec
UPDATE train
SET lineId = ?, currentId = NULL, nextId = NULL
WHERE id = ?;

-- name: DeleteAllTrains :exec
DELETE FROM train;
