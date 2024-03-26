-- +goose Up
-- +goose StatementBegin
SELECT 'UP initial migration';
CREATE TABLE station (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    x REAL,
    y REAL,
    z REAL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE line (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
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
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(currentId) REFERENCES station(id),
    FOREIGN KEY(nextId) REFERENCES station(id) 
);

CREATE TABLE station_line (
    id INTEGER PRIMARY KEY,
    stationId INTEGER,
    lineId INTEGER,
    FOREIGN KEY(stationId) REFERENCES station(id),
    FOREIGN KEY(lineId) REFERENCES station(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'DOWN initial migration';
DROP TABLE station;
DROP TABLE line;
DROP TABLE train;
DROP TABLE station_line;
-- +goose StatementEnd
