-- +goose Up
-- +goose StatementBegin
SELECT 'UP Santo Domingo metro seed';
-- Generated from OpenStreetMap data
-- Reference point: 18.470000, -69.910000
-- Scale: 1 pixel = 100 meters

-- Santo Domingo Metro Stations
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1000, 'Rosa Duarte', 450.97, 271.42, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1001, 'Concepción Bona', 463.39, 276.91, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1002, 'Trina de Moya de Vásquez', 457.86, 271.48, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1003, 'Los Beisbolistas', 336.61, 298.39, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1004, 'Los Beisbolistas', 336.67, 298.34, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1005, 'Juan Pablo Duarte', 403.75, 303.06, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1006, 'Centro de los Héroes', 389.98, 337.22, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1007, 'Francisco Alberto Caamaño Deñó', 393.59, 332.00, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1008, 'Amín Abel Hasbún', 401.70, 327.96, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1009, 'Joaquín Balaguer', 408.43, 321.98, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1010, 'Casandra Damirón', 406.16, 314.61, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1011, 'Profesor Juan Bosch', 404.36, 308.55, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1012, 'Manuel Arturo Peña Batlle', 403.61, 298.43, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1013, 'Pedro Livio Cedeño', 403.04, 289.48, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1014, 'Los Taínos', 402.61, 282.68, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1015, 'Máximo Gómez', 402.21, 274.26, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1016, 'Hermanas Mirabal', 403.08, 262.78, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1017, 'Francisco Gregorio Billini', 361.55, 303.23, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1018, 'María Montez', 347.11, 306.30, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1019, 'Ercilia Pepín', 444.09, 271.69, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1020, 'Manuel de Jesús Galván', 429.53, 282.92, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1021, 'Pedro Francisco Bonó', 353.62, 305.04, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1022, 'Ulises Francisco Espaillat', 369.88, 302.62, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1023, 'Pedro Mir', 375.97, 300.71, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1024, 'Freddy Beras-Goico', 386.44, 302.01, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1025, 'Juan Ulises García Saleta', 397.35, 302.70, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1026, 'Coronel Rafael Tomás Fernández Domínguez', 411.70, 302.85, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1027, 'Mauricio Báez', 414.05, 296.41, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1028, 'Ramón Cáceres', 419.70, 290.52, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1029, 'Horacio Vásquez', 422.92, 287.39, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1030, 'Eduardo Brito', 435.79, 277.81, 0.0, '#FFFFFF');

-- Santo Domingo Metro Lines
INSERT OR IGNORE INTO line (id, name, color) VALUES (100, 'Línea 1', '#E84B28');
INSERT OR IGNORE INTO line (id, name, color) VALUES (101, 'Línea 2', '#0066B3');

-- Station-Line associations
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1000, 1000, 100, 1);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1001, 1001, 100, 2);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1002, 1002, 100, 3);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1003, 1005, 100, 4);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1004, 1006, 100, 5);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1005, 1007, 100, 6);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1006, 1008, 100, 7);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1007, 1009, 100, 8);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1008, 1010, 100, 9);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1009, 1011, 100, 10);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1010, 1012, 100, 11);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1011, 1013, 100, 12);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1012, 1014, 100, 13);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1013, 1015, 100, 14);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1014, 1016, 100, 15);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1015, 1019, 100, 16);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1016, 1020, 100, 17);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1017, 1026, 100, 18);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1018, 1027, 100, 19);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1019, 1028, 100, 20);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1020, 1029, 100, 21);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1021, 1030, 100, 22);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1022, 1003, 101, 1);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1023, 1004, 101, 2);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1024, 1007, 101, 3);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1025, 1017, 101, 4);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1026, 1018, 101, 5);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1027, 1021, 101, 6);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1028, 1022, 101, 7);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1029, 1023, 101, 8);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1030, 1024, 101, 9);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1031, 1025, 101, 10);

-- Edges (station connections)
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1000, 1000, 1001);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1001, 1001, 1002);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1002, 1002, 1005);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1003, 1005, 1006);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1004, 1006, 1007);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1005, 1007, 1008);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1006, 1008, 1009);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1007, 1009, 1010);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1008, 1010, 1011);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1009, 1011, 1012);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1010, 1012, 1013);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1011, 1013, 1014);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1012, 1014, 1015);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1013, 1015, 1016);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1014, 1016, 1019);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1015, 1019, 1020);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1016, 1020, 1026);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1017, 1026, 1027);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1018, 1027, 1028);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1019, 1028, 1029);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1020, 1029, 1030);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1021, 1003, 1004);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1022, 1004, 1007);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1023, 1007, 1017);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1024, 1017, 1018);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1025, 1018, 1021);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1026, 1021, 1022);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1027, 1022, 1023);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1028, 1023, 1024);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1029, 1024, 1025);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'DOWN Santo Domingo metro seed';
DELETE FROM station WHERE id >= 1000;
DELETE FROM line WHERE id >= 100;
DELETE FROM station_line WHERE id >= 1000;
DELETE FROM edge WHERE id >= 1000;
-- +goose StatementEnd
