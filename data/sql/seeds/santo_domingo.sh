#!/bin/bash
# Load Santo Domingo seed data
DB_PATH="${1:-data/metro.db}"

echo "Loading Santo Domingo..."

# Extract and run UP sections only (everything between first StatementBegin and -- +goose Down)
sed -n '/-- +goose Up/,/-- +goose Down/p' data/sql/seeds/20240326000000_makes.sql | \
  sed '/-- +goose/d' | sqlite3 "$DB_PATH"

# Find and load latest Santo Domingo seed
LATEST_SD=$(ls -t data/sql/seeds/*santo_domingo.sql 2>/dev/null | head -1)
if [ -n "$LATEST_SD" ]; then
  sed -n '/-- +goose Up/,/-- +goose Down/p' "$LATEST_SD" | \
    sed '/-- +goose/d' | sqlite3 "$DB_PATH"
  echo "✓ Santo Domingo loaded (19 stations, 2 lines)"
  echo "Note: No trains or schedules yet - run 'cd tools && go run generate_schedules.go' to create them"
else
  echo "Error: No Santo Domingo seed file found. Run 'make import_osm' first."
  exit 1
fi
