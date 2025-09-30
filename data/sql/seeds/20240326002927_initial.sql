-- +goose Up
-- +goose StatementBegin
SELECT 'UP initial seed';
-- Stations data
INSERT INTO station (id, name, x, y, z, color) VALUES (1, 'Station 1', 50.0, 350.0, 0.0, "#FFFFFF");
INSERT INTO station (id, name, x, y, z, color) VALUES (2, 'Station 2', 250.0, 200.0, 0.0, "#F4F4F4");
INSERT INTO station (id, name, x, y, z, color) VALUES (3, 'Station 3', 150.0, 100.0, 0.0, "#DFDFDF");
INSERT INTO station (id, name, x, y, z, color) VALUES (4, 'Station 4', 500.0, 50.0, 0.0, "#D1D1D1");
INSERT INTO station (id, name, x, y, z, color) VALUES (5, 'Station 5', 650.0, 150.0, 0.0, "#BFBFBF");
INSERT INTO station (id, name, x, y, z, color) VALUES (6, 'Station 6', 200.0, 400.0, 0.0, "#AFAFAF");
INSERT INTO station (id, name, x, y, z, color) VALUES (7, 'Station 7', 200.0, 500.0, 0.0, "#FFFFFF");
INSERT INTO station (id, name, x, y, z, color) VALUES (8, 'Station 8', 400.0, 450.0, 0.0, "#F4F4F4");
INSERT INTO station (id, name, x, y, z, color) VALUES (9, 'Station 9', 500.0, 350.0, 0.0, "#DFDFDF");
INSERT INTO station (id, name, x, y, z, color) VALUES (10, 'Station 10', 650.0, 300.0, 0.0, "#D1D1D1");
INSERT INTO station (id, name, x, y, z, color) VALUES (11, 'Station 11', 450.0, 150.0, 0.0, "#BFBFBF");
INSERT INTO station (id, name, x, y, z, color) VALUES (12, 'Station 12', 700.0, 50.0, 0.0, "#AFAFAF");
-- Lines data
INSERT INTO line (id, name, color) VALUES (1, 'Line 1', "#FFFF00");
INSERT INTO line (id, name, color) VALUES (2, 'Line 2', "#23FF00");
INSERT INTO line (id, name, color) VALUES (3, 'Line 3', "#D29100");
INSERT INTO line (id, name, color) VALUES (4, 'Line 4', "#F14900");
-- Make data
INSERT INTO make (id, name, description, acceleration, top_speed, color) VALUES (1, '4-Legged-chu', 'A type of fast train.', 0.003, 0.01, "#0000DD");
INSERT INTO make (id, name, description, acceleration, top_speed, color) VALUES (2, '1-Legged-chu', 'Another type of fast train.', 0.004, 0.06, "#113298");
-- Trains data
INSERT INTO train (name, x, y, z, makeId, currentId, lineId) VALUES ('Cha', 50.0, 350.0, 0.0, 1, 1, 1);
INSERT INTO train (name, x, y, z, makeId, currentId, lineId) VALUES ('Che', 500.0, 50.0, 0.0, 2, 4, 4);
INSERT INTO train (name, x, y, z, makeId, currentId, lineId) VALUES ('Chi', 250.0, 200.0, 0.0, 2, 2, 1);
INSERT INTO train (name, x, y, z, makeId, currentId, lineId) VALUES ('Cho', 400.0, 450.0, 0.0, 1, 8, 3);
INSERT INTO train (name, x, y, z, makeId, currentId, lineId) VALUES ('Chu', 250.0, 200.0, 0.0, 2, 2, 2);
-- Station lines data
INSERT INTO station_line (stationId, lineId, odr) VALUES (1, 1, 1);
INSERT INTO station_line (stationId, lineId, odr) VALUES (2, 1, 2);
INSERT INTO station_line (stationId, lineId, odr) VALUES (4, 1, 3);
INSERT INTO station_line (stationId, lineId, odr) VALUES (12, 1, 4);
INSERT INTO station_line (stationId, lineId, odr) VALUES (3, 2, 1);
INSERT INTO station_line (stationId, lineId, odr) VALUES (2, 2, 2);
INSERT INTO station_line (stationId, lineId, odr) VALUES (6, 2, 3);
INSERT INTO station_line (stationId, lineId, odr) VALUES (7, 2, 4);
INSERT INTO station_line (stationId, lineId, odr) VALUES (8, 3, 1);
INSERT INTO station_line (stationId, lineId, odr) VALUES (9, 3, 2);
INSERT INTO station_line (stationId, lineId, odr) VALUES (10, 3, 3);
INSERT INTO station_line (stationId, lineId, odr) VALUES (11, 4, 1);
INSERT INTO station_line (stationId, lineId, odr) VALUES (4, 4, 2);
INSERT INTO station_line (stationId, lineId, odr) VALUES (5, 4, 3);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'DOWN initial seed';
DELETE FROM station;
DELETE FROM line;
DELETE FROM make;
DELETE FROM trains;
DELETE FROM station_line;
-- +goose StatementEnd
