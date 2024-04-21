package baso

import (
	"github.com/odin-software/metro/internal/dbstore"
)

type CreateEdge struct {
	Fromid int64 `json:"fromId"`
	Toid   int64 `json:"toId"`
}

func (bs *Baso) ListEdges() ([]dbstore.Edge, error) {
	edges, err := bs.queries.GetEdges(bs.ctx)
	if err != nil {
		return nil, err
	}
	return edges, nil
}

func (bs *Baso) CreateEdge(fromId, toId int64) error {
	_, err := bs.queries.CreateEdge(bs.ctx, dbstore.CreateEdgeParams{
		Fromid: fromId,
		Toid:   toId,
	})

	return err
}

func (bs *Baso) CreateEdges(edgesToCreate []CreateEdge) ([]dbstore.Edge, error) {
	for _, edge := range edgesToCreate {
		_, err := bs.GetStationById(edge.Fromid)
		if err != nil {
			return nil, err
		}
		_, err = bs.GetStationById(edge.Toid)
		if err != nil {
			return nil, err
		}
		err = bs.CreateEdge(edge.Fromid, edge.Toid)
		if err != nil {
			return nil, err
		}
	}

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
