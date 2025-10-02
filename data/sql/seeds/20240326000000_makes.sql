-- +goose Up
-- +goose StatementBegin
SELECT 'UP train makes seed';
-- Make data (Realistic metro speeds: ~60-70 km/h)
-- Scale: 1 pixel = 100 meters, 60 ticks/second
-- Make 1: 70 km/h = 0.0032 pixels/tick, accel = 0.0002 pixels/tick² (~1.2 m/s²)
-- Make 2: 60 km/h = 0.0028 pixels/tick, accel = 0.0002 pixels/tick² (~1.2 m/s²)
INSERT OR IGNORE INTO make (id, name, description, acceleration, top_speed, color)
VALUES (1, '4-Legged-chu', 'Fast metro train (70 km/h)', 0.0002, 0.0032, "#0000DD");
INSERT OR IGNORE INTO make (id, name, description, acceleration, top_speed, color)
VALUES (2, '1-Legged-chu', 'Standard metro train (60 km/h)', 0.0002, 0.0028, "#113298");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'DOWN train makes seed';
DELETE FROM make;
-- +goose StatementEnd
