-- +goose Up
-- +goose StatementBegin
SELECT 'UP Santo Domingo metro seed';
-- Generated from OpenStreetMap data
-- Reference point: 18.470000, -69.910000
-- Scale: 1 pixel = 100 meters

-- Santo Domingo Metro Stations
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1000, 'Mamá Tingó', 412.66, 246.84, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1001, 'Rosa Duarte', 445.71, 287.35, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1002, 'Concepción Bona', 458.14, 292.84, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1003, 'Trina de Moya de Vásquez', 452.60, 287.41, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1004, 'José Francisco Peña Gómez', 396.54, 270.46, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1005, 'Gregorio Luperón', 404.75, 265.83, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1006, 'Gregorio Urbano Gilbert', 408.93, 254.23, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1007, 'Juan Pablo Duarte', 398.50, 318.99, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1008, 'Centro de los Héroes', 384.73, 353.16, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1009, 'Francisco Alberto Caamaño Deñó', 388.34, 347.93, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1010, 'Amín Abel Hasbún', 396.45, 343.89, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1011, 'Joaquín Balaguer', 403.18, 337.92, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1012, 'Casandra Damirón', 400.91, 330.54, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1013, 'Profesor Juan Bosch', 399.10, 324.48, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1014, 'Manuel Arturo Peña Batlle', 398.35, 314.37, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1015, 'Pedro Livio Cedeño', 397.78, 305.41, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1016, 'Los Taínos', 397.36, 298.61, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1017, 'Máximo Gómez', 396.96, 290.19, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1018, 'Hermanas Mirabal', 397.83, 278.71, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1019, 'Francisco Gregorio Billini', 356.30, 319.17, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1020, 'María Montez', 341.86, 322.23, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1021, 'Ercilia Pepín', 438.83, 287.63, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1022, 'Manuel de Jesús Galván', 424.27, 298.85, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1023, 'Pedro Francisco Bonó', 348.37, 320.97, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1024, 'Ulises Francisco Espaillat', 364.63, 318.55, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1025, 'Pedro Mir', 370.72, 316.64, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1026, 'Freddy Beras-Goico', 381.19, 317.94, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1027, 'Juan Ulises García Saleta', 392.10, 318.63, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1028, 'Coronel Rafael Tomás Fernández Domínguez', 406.45, 318.78, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1029, 'Mauricio Báez', 408.80, 312.34, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1030, 'Ramón Cáceres', 414.45, 306.45, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1031, 'Horacio Vásquez', 417.67, 303.32, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1032, 'Eduardo Brito', 430.54, 293.74, 0.0, '#FFFFFF');

-- Santo Domingo Metro Lines
INSERT OR IGNORE INTO line (id, name, color) VALUES (100, 'Línea 1', '#E84B28');
INSERT OR IGNORE INTO line (id, name, color) VALUES (101, 'Línea 2', '#0066B3');

-- Station-Line associations
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1000, 1001, 100, 1);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1001, 1002, 100, 2);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1002, 1003, 100, 3);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1003, 1007, 100, 4);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1004, 1008, 100, 5);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1005, 1009, 100, 6);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1006, 1010, 100, 7);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1007, 1011, 100, 8);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1008, 1012, 100, 9);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1009, 1013, 100, 10);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1010, 1014, 100, 11);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1011, 1015, 100, 12);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1012, 1016, 100, 13);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1013, 1017, 100, 14);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1014, 1018, 100, 15);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1015, 1021, 100, 16);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1016, 1022, 100, 17);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1017, 1028, 100, 18);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1018, 1029, 100, 19);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1019, 1030, 100, 20);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1020, 1031, 100, 21);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1021, 1032, 100, 22);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1022, 1009, 101, 1);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1023, 1019, 101, 2);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1024, 1020, 101, 3);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1025, 1023, 101, 4);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1026, 1024, 101, 5);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1027, 1025, 101, 6);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1028, 1026, 101, 7);
INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (1029, 1027, 101, 8);

-- Edges (station connections)
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1000, 1001, 1002);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1001, 1002, 1003);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1002, 1003, 1007);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1003, 1007, 1008);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1004, 1008, 1009);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1005, 1009, 1010);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1006, 1010, 1011);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1007, 1011, 1012);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1008, 1012, 1013);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1009, 1013, 1014);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1010, 1014, 1015);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1011, 1015, 1016);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1012, 1016, 1017);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1013, 1017, 1018);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1014, 1018, 1021);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1015, 1021, 1022);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1016, 1022, 1028);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1017, 1028, 1029);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1018, 1029, 1030);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1019, 1030, 1031);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1020, 1031, 1032);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1021, 1009, 1019);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1022, 1019, 1020);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1023, 1020, 1023);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1024, 1023, 1024);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1025, 1024, 1025);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1026, 1025, 1026);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1027, 1026, 1027);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'DOWN Santo Domingo metro seed';
DELETE FROM station WHERE id >= 1000;
DELETE FROM line WHERE id >= 100;
DELETE FROM station_line WHERE id >= 1000;
DELETE FROM edge WHERE id >= 1000;
-- +goose StatementEnd
