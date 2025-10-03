-- +goose Up
-- +goose StatementBegin
CREATE TABLE station (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    x REAL,
    y REAL,
    z REAL,
    color VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE line (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    color VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE train (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    x REAL NOT NULL,
    y REAL NOT NULL,
    z REAL NOT NULL,
    currentId INTEGER,
    nextId INTEGER,
    makeId INTEGER,
    lineId INTEGER,
    forward INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(currentId) REFERENCES station(id),
    FOREIGN KEY(nextId) REFERENCES station(id),
    FOREIGN KEY(makeId) REFERENCES make(id),
    FOREIGN KEY(lineId) REFERENCES line(id)
);
CREATE TABLE station_line (
    id INTEGER PRIMARY KEY,
    stationId INTEGER,
    lineId INTEGER,
    odr INTEGER,
    FOREIGN KEY(stationId) REFERENCES station(id),
    FOREIGN KEY(lineId) REFERENCES station(id)
);
CREATE TABLE make (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    acceleration REAL,
    top_speed REAL,
    color VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE edge (
    id INTEGER PRIMARY KEY,
    fromId INTEGER NOT NULL,
    toId INTEGER NOT NULL,
    FOREIGN KEY(fromId) REFERENCES station(id),
    FOREIGN KEY(toId) REFERENCES station(id)
);
CREATE TABLE edge_point (
    id INTEGER PRIMARY KEY,
    edgeId INTEGER NOT NULL,
    odr INTEGER NOT NULL,
    x REAL NOT NULL,
    y REAL NOT NULL,
    z REAL NOT NULL,
    FOREIGN KEY(edgeId) REFERENCES edge(id)
);
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
CREATE TABLE schedule (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    train_id INTEGER NOT NULL,
    station_id INTEGER NOT NULL,
    scheduled_time INTEGER NOT NULL,
    sequence_order INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(train_id) REFERENCES train(id) ON DELETE CASCADE,
    FOREIGN KEY(station_id) REFERENCES station(id) ON DELETE CASCADE
);
CREATE INDEX idx_schedule_train ON schedule(train_id);
CREATE INDEX idx_schedule_station ON schedule(station_id);
CREATE INDEX idx_schedule_time ON schedule(scheduled_time);
CREATE INDEX idx_schedule_train_sequence ON schedule(train_id, sequence_order);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE station;
DROP TABLE line;
DROP TABLE train;
DROP TABLE station_line;
DROP TABLE make;
DROP TABLE edge;
DROP TABLE edge_point;
DROP TABLE passenger_event;
DROP TABLE passenger;
-- +goose StatementEnd
