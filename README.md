# Metro v0.0.5

Metro is a simulation project of a transport system and the complex interactions between the different components of the system. The project is developed in Go, and has a web interface to interact with the system.

## Installation

To install the project, you need to have Go installed in your system. You can download it from the [official website](https://golang.org/).

Once you have Go installed, you can clone the repository and install the dependencies with the following commands:

```bash
git clone https://github.com/odin-software/metro
cd metro
go mod download
```

## Quick Start

### 1. Setup Database

Choose which city to simulate:

**Test City** (12 stations, 4 lines, 5 trains):

```bash
make seed_test_city
```

**Santo Domingo** (19 stations, 2 lines, real coordinates):

```bash
make seed_santo_domingo
```

### 2. Run Simulation

```bash
go build && ./metro
```

This will start the visual simulation with Ebiten rendering.

## City Selection

The simulation supports multiple cities:

- **Test City**: Synthetic network for development (70x50 km, unrealistic scale)
- **Santo Domingo**: Real metro data from OpenStreetMap (~12x7 km, realistic scale)

To switch cities, clean the database and reload:

```bash
make seed_santo_domingo  # or seed_test_city
```

To import fresh Santo Domingo data from OpenStreetMap:

```bash
make import_osm
make seed_santo_domingo
```

## Development

We use [goose](https://github.com/pressly/goose) to manage the database migrations. To install it, you can use the following command:

```bash
go get -u github.com/pressly/goose/cmd/goose
```

In order to take advantage of the pre-configured goose options you can set up the following environment variables:

```bash
export GOOSE_DRIVER=sqlite3
export GOOSE_DBSTRING=metro.db
export GOOSE_MIGRATIONS_DIR=migrations
```

To create a migration, you can use the following command:

```bash
goose -dir migrations create <migration_name> sql
```

To run the migrations, you can use the following command:

```bash
goose -dir migrations up
```

## Tools and libraries

For the most part we want to use a minimal quantity of dependencies, but we do use some libraries and tools that we need to give credit to:

- [Go](https://golang.org/): The programming language used to develop the project.
- [echo](https://echo.labstack.com/): The web framework used to develop the web interface.
- [goose](https://github.com/pressly/goose): The library used to manage the database migrations.
