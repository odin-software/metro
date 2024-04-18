package baso

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/odin-software/metro/internal/dbstore"
)

// Baso is a struct that represents the Baso object.
// The Baso object has no fields but does have methods to interact with the database.
type Baso struct {
	ctx     context.Context
	queries *dbstore.Queries
	db      *sql.DB
}

func NewBaso() *Baso {
	ctx := context.Background()
	d, err := sql.Open("sqlite3", "./data/metro.db")
	if err != nil {
		panic(err)
	}

	queries := dbstore.New(d)

	return &Baso{
		ctx:     ctx,
		queries: queries,
		db:      d,
	}
}
