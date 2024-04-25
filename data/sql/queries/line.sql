-- name: ListLines :many
SELECT id, name FROM line
ORDER BY name;

-- name: GetStationsFromLine :many
SELECT
	st.id,
	st.name,
	st.x,
	st.y,
	st.z,
	st.color
FROM
	line ln
	JOIN station_line sl ON ln.id = sl.lineId
	JOIN station st ON sl.stationId = st.id
WHERE
	lineId = ?;

-- name: GetPointsFromLine :many
SELECT
	IFNULL(ep.x, st.x) x,
	IFNULL(ep.y, st.y) y,
	IFNULL(ep.z, st.z) z,
	cj.is_station
FROM
	line ln
	JOIN station_line sl ON ln.id = sl.lineId
	JOIN station st ON sl.stationId = st.id
	LEFT JOIN station_line sln ON sln.lineId = ln.id
		AND sln.odr = sl.odr + 1
	LEFT JOIN edge ed ON (ed.fromId = sl.stationId
		AND ed.toId = sln.stationId) OR (ed.fromId = sln.stationId
		AND ed.toId = sl.stationId)
	CROSS JOIN (
		SELECT
			1 is_station
	UNION
	SELECT
		0) cj
	LEFT JOIN edge_point ep ON ed.id = ep.edgeId
		AND cj.is_station = 0
WHERE
	ln.id = ?
	AND (cj.is_station = 1
		OR ep.id IS NOT NULL)
ORDER BY
	sl.odr,
	cj.is_station DESC,
	ep.odr;

-- name: GetLineById :one
SELECT id, name, color FROM line
WHERE id = ?
LIMIT 1;

-- name: GetLineByName :one
SELECT id, name, color FROM line
WHERE name = ?
LIMIT 1;

-- name: CreateLine :one
INSERT INTO line (name, color)
VALUES (?, ?)
RETURNING id;

-- name: CreateStationLine :one
INSERT INTO station_line (stationId, lineId, odr)
VALUES (?, ?, ?)
RETURNING id;

-- name: UpdateLine :one
UPDATE line
SET name = ?
WHERE id = ?
RETURNING id;

-- name: DeleteLine :exec
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
SELECT ?, ?, COALESCE(MAX(odr), 0) + 1
RETURNING id;

-- name: RemoveStationFromLine :exec
DELETE FROM station_line
WHERE stationId = ? AND lineId = ?;

-- name: DeleteAllLines :exec
DELETE FROM line;

-- name: DeleteAllStationLines :exec
DELETE FROM station_line;