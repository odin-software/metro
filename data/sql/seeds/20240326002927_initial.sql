-- +goose Up
-- +goose StatementBegin
SELECT 'UP initial seed';
-- Stations data
INSERT INTO station (id, name, x, y, z) VALUES (1, 'Station 1', 50.0, 350.0, 0.0);
INSERT INTO station (id, name, x, y, z) VALUES (2, 'Station 2', 250.0, 200.0, 0.0);
INSERT INTO station (id, name, x, y, z) VALUES (3, 'Station 3', 150.0, 100.0, 0.0);
INSERT INTO station (id, name, x, y, z) VALUES (4, 'Station 4', 500.0, 50.0, 0.0);
INSERT INTO station (id, name, x, y, z) VALUES (5, 'Station 5', 650.0, 150.0, 0.0);
INSERT INTO station (id, name, x, y, z) VALUES (6, 'Station 6', 200.0, 400.0, 0.0);
INSERT INTO station (id, name, x, y, z) VALUES (7, 'Station 7', 200.0, 500.0, 0.0);
INSERT INTO station (id, name, x, y, z) VALUES (8, 'Station 8', 400.0, 450.0, 0.0);
INSERT INTO station (id, name, x, y, z) VALUES (9, 'Station 9', 500.0, 350.0, 0.0);
INSERT INTO station (id, name, x, y, z) VALUES (10, 'Station 10', 650.0, 300.0, 0.0);
INSERT INTO station (id, name, x, y, z) VALUES (11, 'Station 11', 450.0, 150.0, 0.0);
INSERT INTO station (id, name, x, y, z) VALUES (12, 'Station 12', 700.0, 50.0, 0.0);
-- Lines data
INSERT INTO line (1, name) VALUES ('Line 1');
INSERT INTO line (2, name) VALUES ('Line 2');
INSERT INTO line (3, name) VALUES ('Line 3');
INSERT INTO line (4, name) VALUES ('Line 4');
-- Make data
INSERT INTO make (id, name, description, acceleration, top_speed) VALUES (1, '4-Legged-chu', 'A type of fast train.', 0.003, 1.0);
INSERT INTO make (id, name, description, acceleration, top_speed) VALUES (2, '1-Legged-chu', 'Another type of fast train.', 0.004, 0.7);
-- Trains data
INSERT INTO train (name, x, y, z, makeId, currentId, lineId) VALUES ('Cha', 50.0, 350.0, 0.0, 1, 1, 1);
INSERT INTO train (name, x, y, z, makeId, currentId, lineId) VALUES ('Che', 500.0, 50.0, 0.0, 2, 4, 4);
INSERT INTO train (name, x, y, z, makeId, currentId, lineId) VALUES ('Chi', 250.0, 200.0, 0.0, 2, 2, 1);
INSERT INTO train (name, x, y, z, makeId, currentId, lineId) VALUES ('Cho', 400.0, 450.0, 0.0, 1, 8, 3);
INSERT INTO train (name, x, y, z, makeId, currentId, lineId) VALUES ('Chu', 250.0, 200.0, 0.0, 2, 2, 2);
-- Station lines data
INSERT INTO station_line (stationId, lineId, order) VALUES (1, 1, 1);
INSERT INTO station_line (stationId, lineId, order) VALUES (2, 1, 2);
INSERT INTO station_line (stationId, lineId, order) VALUES (4, 1, 3);
INSERT INTO station_line (stationId, lineId, order) VALUES (12, 1, 4);
INSERT INTO station_line (stationId, lineId, order) VALUES (3, 2, 1);
INSERT INTO station_line (stationId, lineId, order) VALUES (2, 2, 2);
INSERT INTO station_line (stationId, lineId, order) VALUES (6, 2, 3);
INSERT INTO station_line (stationId, lineId, order) VALUES (7, 2, 4);
INSERT INTO station_line (stationId, lineId, order) VALUES (8, 3, 1);
INSERT INTO station_line (stationId, lineId, order) VALUES (9, 3, 2);
INSERT INTO station_line (stationId, lineId, order) VALUES (10, 3, 3);
INSERT INTO station_line (stationId, lineId, order) VALUES (11, 4, 1);
INSERT INTO station_line (stationId, lineId, order) VALUES (4, 4, 2);
INSERT INTO station_line (stationId, lineId, order) VALUES (5, 4, 3);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'DOWN initial seed';
-- +goose StatementEnd
