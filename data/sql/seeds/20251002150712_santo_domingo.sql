-- +goose Up
-- +goose StatementBegin
SELECT 'UP Santo Domingo metro seed';
-- Generated from OpenStreetMap data
-- Reference point: 18.470000, -69.910000
-- Scale: 1 pixel = 100 meters

-- Santo Domingo Metro Stations
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1000, 'Mamá Tingó', 427.90, 262.81, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1001, 'Gregorio Urbano Gilbert', 425.79, 268.37, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1002, 'Gregorio Luperón', 423.69, 273.92, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1003, 'José Francisco Peña Gómez', 421.58, 279.47, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1004, 'Hermanas Mirabal', 419.48, 285.02, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1005, 'Máximo Gómez', 417.37, 290.57, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1006, 'Los Taínos', 415.27, 296.12, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1007, 'Pedro Livio Cedeño', 413.16, 301.67, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1008, 'Manuel Arturo Peña Batlle', 411.05, 306.11, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1009, 'Juan Pablo Duarte', 408.95, 310.55, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1010, 'Profesor Juan Bosch', 406.84, 314.99, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1011, 'Casandra Damirón', 404.74, 319.43, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1012, 'Joaquín Balaguer', 402.63, 323.87, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1013, 'Amín Abel Hasbún', 400.53, 328.31, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1014, 'Francisco Alberto Caamaño Deñó', 398.42, 332.75, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1015, 'Centro de los Héroes', 396.32, 337.19, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1016, 'María Montez', 443.69, 296.12, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1017, 'Pedro Francisco Bonó', 438.43, 297.23, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1018, 'Francisco Gregorio Billini', 433.16, 298.33, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1019, 'Ulises Francisco Espaillat', 427.90, 299.45, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1020, 'Pedro Mir', 422.63, 300.56, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1021, 'Freddy Beras-Goico', 417.37, 301.67, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1022, 'Juan Ulises García Saleta', 412.11, 302.78, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1023, 'Juan Pablo Duarte', 408.95, 310.55, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1024, 'Coronel Rafael Tomás Fernández Domínguez', 403.68, 307.22, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1025, 'Mauricio Báez', 398.42, 308.33, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1026, 'Ramón Cáceres', 393.16, 309.44, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1027, 'Horacio Vásquez', 387.89, 310.55, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1028, 'Manuel de Jesús Galván', 382.63, 311.66, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1029, 'Eduardo Brito', 377.37, 312.77, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1030, 'Ercilia Pepín', 372.10, 313.88, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1031, 'Rosa Duarte', 366.84, 314.99, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1032, 'Trina de Moya de Vásquez', 361.57, 316.10, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1033, 'Concepción Bona', 356.31, 317.21, 0.0, '#FFFFFF');

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
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1009, 1023, 100, 10);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1010, 1010, 100, 11);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1011, 1011, 100, 12);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1012, 1012, 100, 13);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1013, 1013, 100, 14);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1014, 1014, 100, 15);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1015, 1015, 100, 16);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1016, 1016, 101, 1);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1017, 1017, 101, 2);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1018, 1018, 101, 3);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1019, 1019, 101, 4);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1020, 1020, 101, 5);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1021, 1021, 101, 6);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1022, 1022, 101, 7);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1023, 1023, 101, 8);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1024, 1024, 101, 9);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1025, 1025, 101, 10);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1026, 1026, 101, 11);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1027, 1027, 101, 12);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1028, 1028, 101, 13);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1029, 1029, 101, 14);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1030, 1030, 101, 15);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1031, 1031, 101, 16);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1032, 1032, 101, 17);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1033, 1033, 101, 18);

-- Edges (station connections)
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1000, 1000, 1001);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1001, 1001, 1002);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1002, 1002, 1003);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1003, 1003, 1004);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1004, 1004, 1005);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1005, 1005, 1006);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1006, 1006, 1007);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1007, 1007, 1008);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1008, 1008, 1023);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1009, 1023, 1010);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1010, 1010, 1011);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1011, 1011, 1012);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1012, 1012, 1013);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1013, 1013, 1014);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1014, 1014, 1015);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1015, 1016, 1017);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1016, 1017, 1018);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1017, 1018, 1019);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1018, 1019, 1020);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1019, 1020, 1021);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1020, 1021, 1022);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1021, 1022, 1023);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1022, 1023, 1024);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1023, 1024, 1025);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1024, 1025, 1026);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1025, 1026, 1027);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1026, 1027, 1028);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1027, 1028, 1029);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1028, 1029, 1030);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1029, 1030, 1031);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1030, 1031, 1032);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1031, 1032, 1033);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'DOWN Santo Domingo metro seed';
DELETE FROM station WHERE id >= 1000;
DELETE FROM line WHERE id >= 100;
DELETE FROM station_line WHERE id >= 1000;
DELETE FROM edge WHERE id >= 1000;
-- +goose StatementEnd
