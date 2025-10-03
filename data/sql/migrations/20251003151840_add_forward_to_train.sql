-- +goose Up
-- +goose StatementBegin
ALTER TABLE train ADD COLUMN forward INTEGER NOT NULL DEFAULT 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE train DROP COLUMN forward;
-- +goose StatementEnd
