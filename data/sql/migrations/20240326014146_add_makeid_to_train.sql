-- +goose Up
-- +goose StatementBegin
SELECT 'UP makeid lineid train migration';
ALTER TABLE train ADD COLUMN makeId INTEGER;
ALTER TABLE train ADD COLUMN lineId INTEGER;
ALTER TABLE train ADD FOREIGN KEY(makeId) REFERENCES make(id);
ALTER TABLE train ADD FOREIGN KEY(lineId) REFERENCES line(id);
ALTER TABLE station_line ADD COLUMN order INTEGER;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'DOWN makeid lineid train migration';
ALTER TABLE train DROP FOREIGN KEY makeId;
ALTER TABLE train DROP FOREIGN KEY lineId;
ALTER TABLE train DROP COLUMN makeId;
ALTER TABLE train DROP COLUMN lineId;
ALTER TABLE station_line DROP COLUMN order;
-- +goose StatementEnd
