-- +goose Up
-- +goose StatementBegin
SELECT 'UP edges seed';
-- Edges data
INSERT INTO edge (id, fromId, toId) VALUES (1, 3, 0);
-- Edgespoint data
INSERT INTO edge_point (id, edgeId, order, x, y, z) VALUES (1);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down edges seed';
-- +goose StatementEnd