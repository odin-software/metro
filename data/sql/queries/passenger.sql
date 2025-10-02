-- name: CreatePassenger :one
INSERT INTO passenger (id, name, current_station_id, destination_station_id, state, sentiment, spawn_time)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetPassengerById :one
SELECT * FROM passenger
WHERE id = ?
LIMIT 1;

-- name: GetAllActivePassengers :many
SELECT * FROM passenger
ORDER BY spawn_time DESC;

-- name: GetPassengersByStation :many
SELECT * FROM passenger
WHERE current_station_id = ? AND state = 'waiting'
ORDER BY spawn_time ASC;

-- name: GetPassengersByTrain :many
SELECT * FROM passenger
WHERE current_train_id = ? AND state = 'riding'
ORDER BY spawn_time ASC;

-- name: UpdatePassengerState :exec
UPDATE passenger
SET state = ?, sentiment = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: UpdatePassengerBoarding :exec
UPDATE passenger
SET state = 'riding', current_train_id = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: UpdatePassengerDisembarking :exec
UPDATE passenger
SET state = ?, current_train_id = NULL, current_station_id = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeletePassenger :exec
DELETE FROM passenger
WHERE id = ?;

-- name: DeleteAllPassengers :exec
DELETE FROM passenger;

-- name: CountActivePassengers :one
SELECT COUNT(*) as count FROM passenger;

-- Passenger Event queries

-- name: CreatePassengerEvent :one
INSERT INTO passenger_event (passenger_id, event_type, station_id, train_id, sentiment, metadata)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetPassengerEvents :many
SELECT * FROM passenger_event
WHERE passenger_id = ?
ORDER BY created_at ASC;

-- name: GetRecentEvents :many
SELECT * FROM passenger_event
ORDER BY created_at DESC
LIMIT ?;

-- name: DeletePassengerEvents :exec
DELETE FROM passenger_event
WHERE passenger_id = ?;
