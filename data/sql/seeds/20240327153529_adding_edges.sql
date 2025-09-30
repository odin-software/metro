-- +goose Up
-- +goose StatementBegin
SELECT 'UP edges seed';
-- Edges data
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (1, 1, 2);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (2, 2, 3);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (3, 2, 6);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (4, 2, 4);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (5, 4, 5);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (6, 4, 11);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (7, 4, 12);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (8, 6, 7);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (9, 8, 9);
INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (10, 9, 10);
-- Edgespoint data
INSERT OR IGNORE INTO edge_point (id, edgeId, odr, x, y, z) VALUES (1, 1, 1, 50.0, 250.0, 0.0);
INSERT OR IGNORE INTO edge_point (id, edgeId, odr, x, y, z) VALUES (2, 1, 2, 150.0, 200.0, 0.0);
INSERT OR IGNORE INTO edge_point (id, edgeId, odr, x, y, z) VALUES (3, 2, 1, 250.0, 100.0, 0.0);
INSERT OR IGNORE INTO edge_point (id, edgeId, odr, x, y, z) VALUES (4, 3, 1, 300.0, 300.0, 0.0);
INSERT OR IGNORE INTO edge_point (id, edgeId, odr, x, y, z) VALUES (5, 4, 1, 350.0, 200.0, 0.0);
INSERT OR IGNORE INTO edge_point (id, edgeId, odr, x, y, z) VALUES (6, 4, 2, 400.0, 150.0, 0.0);
INSERT OR IGNORE INTO edge_point (id, edgeId, odr, x, y, z) VALUES (7, 4, 3, 400.0, 50.0, 0.0);
INSERT OR IGNORE INTO edge_point (id, edgeId, odr, x, y, z) VALUES (8, 5, 1, 550.0, 100.0, 0.0);
INSERT OR IGNORE INTO edge_point (id, edgeId, odr, x, y, z) VALUES (9, 5, 2, 600.0, 100.0, 0.0);
INSERT OR IGNORE INTO edge_point (id, edgeId, odr, x, y, z) VALUES (10, 7, 1, 600.0, 50.0, 0.0);
INSERT OR IGNORE INTO edge_point (id, edgeId, odr, x, y, z) VALUES (11, 8, 1, 100.0, 500.0, 0.0);
INSERT OR IGNORE INTO edge_point (id, edgeId, odr, x, y, z) VALUES (12, 9, 1, 500.0, 450.0, 0.0);
INSERT OR IGNORE INTO edge_point (id, edgeId, odr, x, y, z) VALUES (13, 10, 1, 500.0, 250.0, 0.0);
INSERT OR IGNORE INTO edge_point (id, edgeId, odr, x, y, z) VALUES (14, 10, 2, 550.0, 200.0, 0.0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'DOWN edges seed';
DELETE FROM edge;
DELETE FROM edge_point;
-- +goose StatementEnd
