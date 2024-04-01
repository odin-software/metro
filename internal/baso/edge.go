package baso

import (
	"log"

	"github.com/odin-software/metro/internal/dbstore"
)

func (bs *Baso) ListEdges() []dbstore.Edge {
	edges, err := bs.queries.GetEdges(bs.ctx)
	if err != nil {
		log.Fatal(err)
	}
	return edges
}

func (bs *Baso) ListEdgePoints(id int64) []dbstore.GetEdgePointsRow {
	eps, err := bs.queries.GetEdgePoints(bs.ctx, id)
	if err != nil {
		log.Fatal(err)
	}
	return eps
}
