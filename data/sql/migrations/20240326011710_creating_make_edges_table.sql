-- +goose Up
-- +goose StatementBegin
SELECT 'UP creating make edges table migration';
CREATE TABLE make (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    acceleration REAL,
    top_speed REAL,
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
    order INTEGER NOT NULL,
    x REAL NOT NULL,
    y REAL NOT NULL,
    z REAL NOT NULL,
    FOREIGN KEY(edgeId) REFERENCES edge(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'DOWN creating make edges table migration';
DROP TABLE make;
DROP TABLE edge;
DROP TABLE edge_point;
-- +goose StatementEnd
