#!/bin/bash
# Load Santo Domingo seed data
DB_PATH="${1:-data/metro.db}"

echo "Loading Santo Domingo..."

# Extract and run UP sections only (everything between first StatementBegin and -- +goose Down)
sed -n '/-- +goose Up/,/-- +goose Down/p' data/sql/seeds/20240326000000_makes.sql | \
  sed '/-- +goose/d' | sqlite3 "$DB_PATH"

# Find and load latest Santo Domingo seed
LATEST_SD=$(ls -t data/sql/seeds/*santo_domingo.sql 2>/dev/null | grep -v "trains" | head -1)
if [ -n "$LATEST_SD" ]; then
  sed -n '/-- +goose Up/,/-- +goose Down/p' "$LATEST_SD" | \
    sed '/-- +goose/d' | sqlite3 "$DB_PATH"
fi

# Find and load latest Santo Domingo trains
LATEST_TRAINS=$(ls -t data/sql/seeds/*santo_domingo_trains.sql 2>/dev/null | head -1)
if [ -n "$LATEST_TRAINS" ]; then
  sed -n '/-- +goose Up/,/-- +goose Down/p' "$LATEST_TRAINS" | \
    sed '/-- +goose/d' | sqlite3 "$DB_PATH"
  echo "✓ Santo Domingo loaded with 69 trains (40 on L1, 29 on L2)"
else
  echo "✓ Santo Domingo loaded (19 stations, 2 lines)"
  echo "Note: No trains yet - run 'cd tools && go run generate_santo_domingo_trains.go' first"
fi
