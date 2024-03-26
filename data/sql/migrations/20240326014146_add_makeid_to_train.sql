-- +goose Up
-- +goose StatementBegin
SELECT 'UP creating make edges table migration';
ALTER TABLE train ADD COLUMN makeId INTEGER;
ALTER TABLE train ADD FOREIGN KEY(makeId) REFERENCES make(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'DOWN creating make edges table migration';
-- +goose StatementEnd
