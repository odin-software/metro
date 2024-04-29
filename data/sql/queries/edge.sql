-- name: GetEdges :many
SELECT id, fromId, toId
FROM edge;

-- name: GetEdgePoints :many
SELECT id, edgeId, X, Y, Z, odr
FROM edge_point
WHERE edgeId = ?
ORDER BY odr;

-- name: CreateEdge :one
INSERT INTO edge (fromId, toId) 
VALUES (?, ?)
RETURNING id;

-- name: DeleteAllEdges :exec
DELETE FROM edge;

-- name: DeleteAllEdgePoints :exec
DELETE FROM edge_point;