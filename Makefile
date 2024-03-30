GOOSE_DRIVER=sqlite3
GOOSE_DBSTRING=data/metro.db

#################
# Data commands #
#################

# Target: run goose migrations files.
run_migrations:
	echo "Running migrations"
	goose -dir "data/sql/migrations" up

# Target: create a new goose migration file with a timestamp.
create_migration:
	echo "Creating migration $(name)"
	goose -dir "data/sql/migrations" create "$(name)" sql 

# Target: create a new seed file with a timestamp.
create_seed:
	echo "Creating seed $(name)"
	goose -dir "data/sql/seeds" create "$(name)" sql 

# Target: run goose seeding files.
run_seeds:
	echo "Running seeds"
	goose -dir "data/sql/seeds" -no-versioning up

run_data: run_migrations run_seeds

# Target: generate sqlc types in go.
generate_sqlc:
	echo "Generating sqlc types"
	sqlc generate