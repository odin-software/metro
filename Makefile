GOOSE_DRIVER=sqlite3
GOOSE_DBSTRING=data/metro.db

.PHONY: help
help:
	@echo "Metro Simulation - Available Commands"
	@echo ""
	@echo "Quick Start:"
	@echo "  make seed_test_city          - Setup test city (12 stations, 4 lines, 5 trains)"
	@echo "  make setup_santo_domingo     - Complete Santo Domingo setup (69 trains)"
	@echo ""
	@echo "City Setup (manual):"
	@echo "  make seed_santo_domingo      - Load Santo Domingo (without regenerating)"
	@echo "  make import_osm              - Import fresh Santo Domingo data from OSM"
	@echo "  make generate_santo_domingo_trains - Generate 69 trains"
	@echo "  make generate_schedules      - Generate schedules for existing trains"
	@echo ""
	@echo "Database Maintenance:"
	@echo "  make run_migrations          - Run database migrations"
	@echo "  make clean_city_data         - Clear all city data (keeps schema)"
	@echo ""
	@echo "Development:"
	@echo "  make generate_sqlc           - Generate Go code from SQL queries"
	@echo "  make create_migration        - Create new migration (name=...)"
	@echo "  make create_seed             - Create new seed file (name=...)"
	@echo ""

#################
# Data commands #
#################

# Target: run goose migrations files.
run_migrations:
	echo "Running migrations"
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir "data/sql/migrations" up

# Target: create a new goose migration file with a timestamp.
create_migration:
	echo "Creating migration $(name)"
	goose -dir "data/sql/migrations" create "$(name)" sql

# Target: create a new seed file with a timestamp.
create_seed:
	echo "Creating seed $(name)"
	goose -dir "data/sql/seeds" create "$(name)" sql

# Target: setup test city (with trains and schedules)
seed_test_city: run_migrations clean_city_data
	@bash data/sql/seeds/test_city.sh $(GOOSE_DBSTRING)

# Target: setup Santo Domingo (with trains and schedules)
seed_santo_domingo: run_migrations clean_city_data
	@bash data/sql/seeds/santo_domingo.sh $(GOOSE_DBSTRING)

# Target: generate Santo Domingo trains (40 for L1, 29 for L2)
generate_santo_domingo_trains:
	@echo "Generating Santo Domingo trains..."
	@cd tools && go run generate_santo_domingo_trains.go
	@echo "✓ Trains generated"

# Target: generate schedules for all trains in the database
generate_schedules:
	@echo "Generating schedules..."
	@cd tools && go run generate_schedules.go > ../data/sql/seeds/schedules_generated.sql
	@sqlite3 $(GOOSE_DBSTRING) < data/sql/seeds/schedules_generated.sql
	@echo "✓ Schedules generated and loaded"

# Target: complete Santo Domingo setup (recommended workflow)
setup_santo_domingo: run_migrations clean_city_data import_osm generate_santo_domingo_trains
	@echo "Loading Santo Domingo data..."
	@bash data/sql/seeds/santo_domingo.sh $(GOOSE_DBSTRING)
	@$(MAKE) generate_schedules
	@echo "✓ Santo Domingo fully configured with 69 trains and schedules"

# Target: clean city-specific data (keeps migrations)
clean_city_data:
	@echo "Cleaning city data..."
	@sqlite3 $(GOOSE_DBSTRING) "DELETE FROM passenger; DELETE FROM train; DELETE FROM edge_point; DELETE FROM edge; DELETE FROM station_line; DELETE FROM line; DELETE FROM station; DELETE FROM schedule;"
	@echo "✓ City data cleaned"

# Target: generate sqlc types in go.
generate_sqlc:
	echo "Generating sqlc types"
	sqlc generate

# Target: import real metro data from OpenStreetMap
import_osm:
	echo "Importing Santo Domingo metro data from OSM"
	cd tools && go run import_osm.go
