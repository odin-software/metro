GOOSE_DRIVER=sqlite3
GOOSE_DBSTRING=data/metro.db

# Target: run goose migrations files.
run_migrations:
	echo "Running migrations"
	goose -dir "data/sql/migrations" up

# Target: create a new goose migration file with a timestamp.
create_migration:
	echo "Creating migration $(name)"
	goose -dir "data/sql/migrations" create "$(name)" sql 

# Target: generate sqlc types in go.
generate_sqlc:
	echo "Generating sqlc types"
	sqlc generate