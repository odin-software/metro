#!/bin/bash
# Load test city seed data
DB_PATH="${1:-data/metro.db}"

echo "Loading test city..."

# Extract and run UP sections only (everything between first StatementBegin and -- +goose Down)
sed -n '/-- +goose Up/,/-- +goose Down/p' data/sql/seeds/20240326000000_makes.sql | \
  sed '/-- +goose/d' | sqlite3 "$DB_PATH"

sed -n '/-- +goose Up/,/-- +goose Down/p' data/sql/seeds/20240326002927_initial.sql | \
  sed '/-- +goose/d' | sqlite3 "$DB_PATH"

sed -n '/-- +goose Up/,/-- +goose Down/p' data/sql/seeds/20240327153529_adding_edges.sql | \
  sed '/-- +goose/d' | sqlite3 "$DB_PATH"

sed -n '/-- +goose Up/,/-- +goose Down/p' data/sql/seeds/20251001200000_schedules.sql | \
  sed '/-- +goose/d' | sqlite3 "$DB_PATH"

TRAIN_COUNT=$(sqlite3 "$DB_PATH" "SELECT COUNT(*) FROM train;")
echo "✓ Test city loaded (12 stations, 4 lines, $TRAIN_COUNT trains with schedules)"
