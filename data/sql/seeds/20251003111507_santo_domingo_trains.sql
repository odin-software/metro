-- +goose Up
-- +goose StatementBegin
SELECT 'UP Santo Domingo trains seed';
-- Generated: 2025-10-03 11:15:07
-- Real-life fleet: Línea 1 = 40 trains, Línea 2 = 29 trains

-- Línea 1 Trains (40 trains - distributed bidirectionally)
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2000, 'L1-T01', 412.66, 246.84, 0.0, 1000, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2001, 'L1-T02', 384.73, 353.16, 0.0, 1008, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2002, 'L1-T03', 408.93, 254.23, 0.0, 1006, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2003, 'L1-T04', 388.34, 347.93, 0.0, 1009, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2004, 'L1-T05', 404.75, 265.83, 0.0, 1005, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2005, 'L1-T06', 396.45, 343.89, 0.0, 1010, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2006, 'L1-T07', 396.54, 270.46, 0.0, 1004, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2007, 'L1-T08', 403.18, 337.92, 0.0, 1011, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2008, 'L1-T09', 397.83, 278.71, 0.0, 1018, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2009, 'L1-T10', 400.91, 330.54, 0.0, 1012, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2010, 'L1-T11', 396.96, 290.19, 0.0, 1017, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2011, 'L1-T12', 399.10, 324.48, 0.0, 1013, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2012, 'L1-T13', 397.36, 298.61, 0.0, 1016, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2013, 'L1-T14', 398.50, 318.99, 0.0, 1007, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2014, 'L1-T15', 397.78, 305.41, 0.0, 1015, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2015, 'L1-T16', 398.35, 314.37, 0.0, 1014, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2016, 'L1-T17', 398.35, 314.37, 0.0, 1014, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2017, 'L1-T18', 397.78, 305.41, 0.0, 1015, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2018, 'L1-T19', 398.50, 318.99, 0.0, 1007, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2019, 'L1-T20', 397.36, 298.61, 0.0, 1016, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2020, 'L1-T21', 399.10, 324.48, 0.0, 1013, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2021, 'L1-T22', 396.96, 290.19, 0.0, 1017, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2022, 'L1-T23', 400.91, 330.54, 0.0, 1012, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2023, 'L1-T24', 397.83, 278.71, 0.0, 1018, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2024, 'L1-T25', 403.18, 337.92, 0.0, 1011, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2025, 'L1-T26', 396.54, 270.46, 0.0, 1004, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2026, 'L1-T27', 396.45, 343.89, 0.0, 1010, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2027, 'L1-T28', 404.75, 265.83, 0.0, 1005, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2028, 'L1-T29', 388.34, 347.93, 0.0, 1009, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2029, 'L1-T30', 408.93, 254.23, 0.0, 1006, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2030, 'L1-T31', 384.73, 353.16, 0.0, 1008, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2031, 'L1-T32', 412.66, 246.84, 0.0, 1000, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2032, 'L1-T33', 412.66, 246.84, 0.0, 1000, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2033, 'L1-T34', 384.73, 353.16, 0.0, 1008, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2034, 'L1-T35', 408.93, 254.23, 0.0, 1006, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2035, 'L1-T36', 388.34, 347.93, 0.0, 1009, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2036, 'L1-T37', 404.75, 265.83, 0.0, 1005, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2037, 'L1-T38', 396.45, 343.89, 0.0, 1010, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2038, 'L1-T39', 396.54, 270.46, 0.0, 1004, 1, 100);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2039, 'L1-T40', 403.18, 337.92, 0.0, 1011, 1, 100);

-- Línea 2 Trains (29 trains - distributed bidirectionally)
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2040, 'L2-T01', 341.86, 322.23, 0.0, 1020, 1, 101);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2041, 'L2-T02', 458.14, 292.84, 0.0, 1002, 1, 101);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2042, 'L2-T03', 348.37, 320.97, 0.0, 1023, 1, 101);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2043, 'L2-T04', 452.60, 287.41, 0.0, 1003, 1, 101);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2044, 'L2-T05', 356.30, 319.17, 0.0, 1019, 1, 101);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2045, 'L2-T06', 445.71, 287.35, 0.0, 1001, 1, 101);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2046, 'L2-T07', 364.63, 318.55, 0.0, 1024, 1, 101);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2047, 'L2-T08', 438.83, 287.63, 0.0, 1021, 1, 101);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2048, 'L2-T09', 370.72, 316.64, 0.0, 1025, 1, 101);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2049, 'L2-T10', 430.54, 293.74, 0.0, 1032, 1, 101);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2050, 'L2-T11', 381.19, 317.94, 0.0, 1026, 1, 101);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2051, 'L2-T12', 424.27, 298.85, 0.0, 1022, 1, 101);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2052, 'L2-T13', 392.10, 318.63, 0.0, 1027, 1, 101);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2053, 'L2-T14', 417.67, 303.32, 0.0, 1031, 1, 101);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2054, 'L2-T15', 398.50, 318.99, 0.0, 1007, 1, 101);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2055, 'L2-T16', 414.45, 306.45, 0.0, 1030, 1, 101);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2056, 'L2-T17', 406.45, 318.78, 0.0, 1028, 1, 101);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2057, 'L2-T18', 408.80, 312.34, 0.0, 1029, 1, 101);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2058, 'L2-T19', 408.80, 312.34, 0.0, 1029, 1, 101);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2059, 'L2-T20', 406.45, 318.78, 0.0, 1028, 1, 101);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2060, 'L2-T21', 414.45, 306.45, 0.0, 1030, 1, 101);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2061, 'L2-T22', 398.50, 318.99, 0.0, 1007, 1, 101);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2062, 'L2-T23', 417.67, 303.32, 0.0, 1031, 1, 101);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2063, 'L2-T24', 392.10, 318.63, 0.0, 1027, 1, 101);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2064, 'L2-T25', 424.27, 298.85, 0.0, 1022, 1, 101);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2065, 'L2-T26', 381.19, 317.94, 0.0, 1026, 1, 101);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2066, 'L2-T27', 430.54, 293.74, 0.0, 1032, 1, 101);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2067, 'L2-T28', 370.72, 316.64, 0.0, 1025, 1, 101);
INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId) VALUES (2068, 'L2-T29', 438.83, 287.63, 0.0, 1021, 1, 101);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'DOWN Santo Domingo trains seed';
DELETE FROM train WHERE id >= 2000 AND id < 3000;
-- +goose StatementEnd
