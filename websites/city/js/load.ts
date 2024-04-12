import {
  GET_EDGES_URL,
  GET_EDGE_POINTS_URL,
  GET_STATIONS_URL,
} from "./utils/consts.js";

async function fetchMetro<T>(url: string): Promise<T> {
  const response = await fetch(url);
  if (!response.ok) {
    throw new Error(response.statusText);
  }
  return await (response.json() as Promise<T>);
}

export async function getStations() {
  const response = await fetchMetro<RequestStation[]>(GET_STATIONS_URL);

  return response;
}

export async function getEdges() {
  const response = await fetchMetro<RequestEdge[]>(GET_EDGES_URL);

  return response;
}

export async function getEdgesPoints(id: number) {
  const response = await fetchMetro<RequestEdgePoint[] | null>(
    GET_EDGE_POINTS_URL(id)
  );

  return response;
}
