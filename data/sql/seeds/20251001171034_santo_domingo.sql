-- +goose Up
-- +goose StatementBegin
SELECT 'UP Santo Domingo metro seed';
-- Generated from OpenStreetMap data
-- Reference point: 18.470000, -69.910000
-- Scale: 1 pixel = 100 meters

-- Santo Domingo Metro Stations
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1000, 'Villa Mella', 386.63, 260.82, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1001, 'Hermanas Mirabal', 395.26, 267.70, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1002, 'Máximo Gómez', 400.32, 279.91, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1003, 'Bartolomé Colón', 405.05, 290.45, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1004, 'Mama Tingó', 407.79, 298.45, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1005, 'Juan Pablo Duarte', 411.37, 303.11, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1006, 'María Montez', 414.21, 309.55, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1007, 'Casandra Damirón', 417.37, 316.76, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1008, 'Francisco del Rosario Sánchez', 420.21, 323.86, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1009, 'Pedro Livio Cedeño', 422.32, 329.64, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1010, 'Los Taínos', 425.06, 336.85, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1011, 'Centro de Los Héroes', 377.89, 295.45, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1012, 'Joaquín Balaguer', 373.36, 286.01, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1013, 'Eduardo Brito', 369.26, 277.02, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1014, 'Francisco Alberto Caamaño Deñó', 377.89, 295.45, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1015, 'La Julia', 362.52, 298.56, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1016, 'Antonio Duvergé', 347.78, 301.89, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1017, 'Hermanas Mirabal 2', 332.62, 305.00, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1018, 'Manuel María Valencia', 317.35, 308.33, 0.0, '#FFFFFF');

-- Santo Domingo Metro Lines
INSERT OR IGNORE INTO line (id, name, color) VALUES (100, 'Línea 1', '#E84B28');
INSERT OR IGNORE INTO line (id, name, color) VALUES (101, 'Línea 2', '#0066B3');

-- Station-Line associations
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1000, 1000, 100, 1);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1001, 1001, 100, 2);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1002, 1002, 100, 3);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1003, 1003, 100, 4);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1004, 1004, 100, 5);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1005, 1005, 100, 6);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1006, 1006, 100, 7);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1007, 1007, 100, 8);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1008, 1008, 100, 9);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1009, 1009, 100, 10);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1010, 1010, 100, 11);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1011, 1011, 100, 12);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1012, 1012, 100, 13);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1013, 1013, 100, 14);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1014, 1014, 101, 1);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1015, 1015, 101, 2);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1016, 1016, 101, 3);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1017, 1017, 101, 4);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1018, 1018, 101, 5);

-- Edges (station connections)
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1000, 1000, 1001);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1001, 1001, 1002);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1002, 1002, 1003);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1003, 1003, 1004);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1004, 1004, 1005);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1005, 1005, 1006);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1006, 1006, 1007);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1007, 1007, 1008);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1008, 1008, 1009);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1009, 1009, 1010);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1010, 1010, 1011);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1011, 1011, 1012);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1012, 1012, 1013);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1013, 1014, 1015);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1014, 1015, 1016);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1015, 1016, 1017);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1016, 1017, 1018);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'DOWN Santo Domingo metro seed';
DELETE FROM station WHERE id >= 1000;
DELETE FROM line WHERE id >= 100;
DELETE FROM station_line WHERE id >= 1000;
DELETE FROM edge WHERE id >= 1000;
-- +goose StatementEnd
