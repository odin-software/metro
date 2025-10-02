-- +goose Up
-- +goose StatementBegin
SELECT 'UP Santo Domingo metro seed';
-- Generated from OpenStreetMap data
-- Reference point: 18.470000, -69.910000
-- Scale: 1 pixel = 100 meters

-- Santo Domingo Metro Stations
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1000, 'Rosa Duarte', 442.60, 255.45, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1001, 'Concepción Bona', 455.03, 260.94, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1002, 'Trina de Moya de Vásquez', 449.49, 255.52, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1003, 'Los Beisbolistas', 328.23, 282.43, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1004, 'Los Beisbolistas', 328.30, 282.38, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1005, 'Juan Pablo Duarte', 395.38, 287.10, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1006, 'Centro de los Héroes', 381.61, 321.26, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1007, 'Francisco Alberto Caamaño Deñó', 385.22, 316.04, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1008, 'Amín Abel Hasbún', 393.33, 311.99, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1009, 'Joaquín Balaguer', 400.06, 306.02, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1010, 'Casandra Damirón', 397.79, 298.65, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1011, 'Profesor Juan Bosch', 395.99, 292.58, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1012, 'Manuel Arturo Peña Batlle', 395.24, 282.47, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1013, 'Pedro Livio Cedeño', 394.66, 273.52, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1014, 'Los Taínos', 394.24, 266.72, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1015, 'Máximo Gómez', 393.84, 258.29, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1016, 'Hermanas Mirabal', 394.71, 246.81, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1017, 'Francisco Gregorio Billini', 353.17, 287.27, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1018, 'María Montez', 338.74, 290.33, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1019, 'Ercilia Pepín', 435.72, 255.73, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1020, 'Manuel de Jesús Galván', 421.16, 266.96, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1021, 'Pedro Francisco Bonó', 345.24, 289.08, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1022, 'Ulises Francisco Espaillat', 361.50, 286.66, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1023, 'Pedro Mir', 367.60, 284.74, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1024, 'Freddy Beras-Goico', 378.06, 286.05, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1025, 'Juan Ulises García Saleta', 388.98, 286.74, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1026, 'Coronel Rafael Tomás Fernández Domínguez', 403.33, 286.88, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1027, 'Mauricio Báez', 405.68, 280.45, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1028, 'Ramón Cáceres', 411.33, 274.55, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1029, 'Horacio Vásquez', 414.55, 271.43, 0.0, '#FFFFFF');
INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (1030, 'Eduardo Brito', 427.42, 261.84, 0.0, '#FFFFFF');

-- Santo Domingo Metro Lines

-- Station-Line associations

-- Edges (station connections)

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'DOWN Santo Domingo metro seed';
DELETE FROM station WHERE id >= 1000;
DELETE FROM line WHERE id >= 100;
DELETE FROM station_line WHERE id >= 1000;
DELETE FROM edge WHERE id >= 1000;
-- +goose StatementEnd
