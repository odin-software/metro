-- name: CreateScheduleEntry :one
INSERT INTO schedule (train_id, station_id, scheduled_time, sequence_order)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetScheduleForTrain :many
SELECT * FROM schedule
WHERE train_id = ?
ORDER BY sequence_order ASC;

-- name: GetScheduleForStation :many
SELECT * FROM schedule
WHERE station_id = ?
ORDER BY scheduled_time ASC;

-- name: GetNextScheduledStop :one
SELECT * FROM schedule
WHERE train_id = ? AND sequence_order > ?
ORDER BY sequence_order ASC
LIMIT 1;

-- name: GetScheduleByTrainAndStation :one
SELECT * FROM schedule
WHERE train_id = ? AND station_id = ?
LIMIT 1;

-- name: DeleteScheduleForTrain :exec
DELETE FROM schedule WHERE train_id = ?;

-- name: GetAllSchedules :many
SELECT * FROM schedule
ORDER BY train_id, sequence_order;
