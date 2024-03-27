-- +goose Up
-- +goose StatementBegin
SELECT 'UP makeid lineid train migration';
ALTER TABLE train ADD COLUMN makeId INTEGER REFERENCES make(id);
ALTER TABLE train ADD COLUMN lineId INTEGER REFERENCES line(id);
ALTER TABLE station_line ADD COLUMN odr INTEGER;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'DOWN makeid lineid train migration';
ALTER TABLE train DROP COLUMN makeId;
ALTER TABLE train DROP COLUMN lineId;
ALTER TABLE station_line DROP COLUMN odr;
-- +goose StatementEnd
