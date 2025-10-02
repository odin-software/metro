-- +goose Up
-- +goose StatementBegin
-- Schedule table stores the expected arrival times for trains at stations
-- scheduled_time is in seconds since midnight (e.g., 28800 = 8:00 AM)
CREATE TABLE schedule (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    train_id INTEGER NOT NULL,
    station_id INTEGER NOT NULL,
    scheduled_time INTEGER NOT NULL, -- seconds since midnight (0-86399)
    sequence_order INTEGER NOT NULL, -- order in the route (1 = first stop, 2 = second, etc.)
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(train_id) REFERENCES train(id) ON DELETE CASCADE,
    FOREIGN KEY(station_id) REFERENCES station(id) ON DELETE CASCADE
);

-- Indexes for faster lookups
CREATE INDEX idx_schedule_train ON schedule(train_id);
CREATE INDEX idx_schedule_station ON schedule(station_id);
CREATE INDEX idx_schedule_time ON schedule(scheduled_time);
CREATE INDEX idx_schedule_train_sequence ON schedule(train_id, sequence_order);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE schedule;
-- +goose StatementEnd
