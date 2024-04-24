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
-- +goose StatementEnd