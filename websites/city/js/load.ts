import {
  GET_EDGES_URL,
  GET_EDGE_POINTS_URL,
  GET_LINES_URL,
  GET_PAUSE_LOOP,
  GET_PLAY_LOOP,
  GET_STATIONS_URL,
} from "./utils/consts.js";

async function fetchMetro<T>(url: string): Promise<T> {
  const response = await fetch(url);
  if (!response.ok && !(response.status === 404)) {
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

export async function getEdgesPoints(
  id: number
): Promise<RequestEdgePoint[] | null> {
  const response = await fetch(GET_EDGE_POINTS_URL(id));

  if (!response.ok) {
    return null;
  }

  return await (response.json() as Promise<RequestEdgePoint[]>);
}

export async function getLines() {
  const response = await fetchMetro<RequestLine[]>(GET_LINES_URL);

  return response;
}

export function playLoop() {
  fetch(GET_PLAY_LOOP);
}

export function pauseLoop() {
  fetch(GET_PAUSE_LOOP);
}

export async function createStations(sts: RequestCreateStation[]) {
  const response = await fetch(GET_STATIONS_URL, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(sts),
  });

  if (!response.ok) {
    throw new Error("couldn't create stations");
  }

  return (await response.json()) as Promise<RequestStation[]>;
}

export async function createEdges(edges: RequestCreateEdge[]) {
  const response = await fetch(GET_EDGES_URL, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(edges),
  });

  if (!response.ok) {
    throw new Error("couldn't create stations");
  }
}
