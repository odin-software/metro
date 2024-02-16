# Metro

Metro is a simulation project of a transport system and the complex interactions between the different components of the system. The project is developed in Go, and has a web interface to interact with the system.

## Installation

To install the project, you need to have Go installed in your system. You can download it from the [official website](https://golang.org/).

Once you have Go installed, you can clone the repository and install the dependencies with the following commands:

```bash
git clone https://github.com/odin-software/metro
cd metro
go mod download
```

## Usage

To run the project, you can use the following command:

```bash
go build && ./metro
```

This will start the report interface at `http://localhost:4440` and the simulation will
start running in the background.

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

To run the migrations, you can use the following command:

```bash
goose -dir migrations up
```

## Tools and libraries

For the most part we want to use a minimal quantity of dependencies, but we do use some libraries and tools that we need to give credit to:

- [Go](https://golang.org/): The programming language used to develop the project.
- [echo](https://echo.labstack.com/): The web framework used to develop the web interface.
- [goose](https://github.com/pressly/goose): The library used to manage the database migrations.
