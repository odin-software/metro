-- name: GetEdges :many
SELECT id, fromId, toId
FROM edge;

-- name: GetEdgePoints :many
SELECT id, edgeId, X, Y, Z, odr
FROM edge_point
WHERE edgeId = ?
ORDER BY odr;