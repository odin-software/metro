-- +goose Up
-- +goose StatementBegin
CREATE TABLE passenger (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    current_station_id INTEGER NOT NULL,
    destination_station_id INTEGER NOT NULL,
    current_train_id INTEGER,
    state TEXT NOT NULL,
    sentiment REAL NOT NULL DEFAULT 100.0,
    spawn_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(current_station_id) REFERENCES station(id),
    FOREIGN KEY(destination_station_id) REFERENCES station(id),
    FOREIGN KEY(current_train_id) REFERENCES train(id)
);

CREATE INDEX idx_passenger_current_station ON passenger(current_station_id);
CREATE INDEX idx_passenger_destination_station ON passenger(destination_station_id);
CREATE INDEX idx_passenger_current_train ON passenger(current_train_id);
CREATE INDEX idx_passenger_state ON passenger(state);

CREATE TABLE passenger_event (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    passenger_id TEXT NOT NULL,
    event_type TEXT NOT NULL,
    station_id INTEGER,
    train_id INTEGER,
    sentiment REAL,
    metadata TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(passenger_id) REFERENCES passenger(id),
    FOREIGN KEY(station_id) REFERENCES station(id),
    FOREIGN KEY(train_id) REFERENCES train(id)
);

CREATE INDEX idx_passenger_event_passenger_id ON passenger_event(passenger_id);
CREATE INDEX idx_passenger_event_type ON passenger_event(event_type);
CREATE INDEX idx_passenger_event_created_at ON passenger_event(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS passenger_event;
DROP TABLE IF EXISTS passenger;
-- +goose StatementEnd
