package baso

import (
	"context"
	"database/sql"
	"log"

	"github.com/odin-software/metro/internal/dbstore"
)

func Stations() {
	ctx := context.Background()
	d, err := sql.Open("sqlite3", "metro.db")
	if err != nil {
		panic(err)
	}

	queries := dbstore.New(d)
	stations, err := queries.ListStations(ctx)
	if err != nil {
		panic(err)
	}

	log.Println(stations)
}
