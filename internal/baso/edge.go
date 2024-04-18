package baso

import (
	"github.com/odin-software/metro/internal/dbstore"
)

func (bs *Baso) ListEdges() ([]dbstore.Edge, error) {
	edges, err := bs.queries.GetEdges(bs.ctx)
	if err != nil {
		return nil, err
	}
	return edges, nil
}

func (bs *Baso) ListEdgePoints(id int64) ([]dbstore.GetEdgePointsRow, error) {
	eps, err := bs.queries.GetEdgePoints(bs.ctx, id)
	if err != nil {
		return nil, err
	}
	return eps, nil
}
