-- name: ListLines :many
SELECT id, name FROM line ORDER BY id;

-- name: GetLineById :one
SELECT id, name FROM line
WHERE id = ?
LIMIT 1;

-- name: GetLineByName :one
SELECT id, name FROM line
WHERE name = ?
LIMIT 1;

-- name: CreateLine :one
INSERT INTO line (name)
VALUES (?);

-- name: UpdateLine :one
UPDATE line
SET name = ?
WHERE id = ?;

-- name: DeleteLine :one
DELETE FROM line 
WHERE id = ?;

-- name: GetLineStations :many
SELECT stationId FROM station_line
WHERE lineId = ?
ORDER by odr;

-- name: GetLineTrains :many
SELECT id, name FROM train
WHERE lineId = ?
ORDER BY id;

-- name: AddStationToLine :one
INSERT INTO station_line (stationId, lineId, odr)
SELECT ?, ?, COALESCE(MAX(odr), 0) + 1;

-- name: RemoveStationFromLine :one
DELETE FROM station_line
WHERE stationId = ? AND lineId = ?;