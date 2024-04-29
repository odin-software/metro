import {
  GET_EDGES_URL,
  GET_EDGE_POINTS_URL,
  GET_LINES_URL,
  GET_PAUSE_LOOP,
  GET_PLAY_LOOP,
  GET_STATIONS_URL,
  GET_TRAINS_URL,
  UPDATE_TRAIN_LINE_URL,
} from "./utils/consts.js";

async function fetchMetro<T>(url: string): Promise<Response> {
  const response = await fetch(url);
  if (!response.ok && !(response.status === 404)) {
    throw new Error(response.statusText);
  }

  return response;
}

export async function getStations() {
  const response = await fetchMetro<RequestStation[]>(GET_STATIONS_URL);
  if (!response.ok) {
    return [];
  }

  return await response.json();
}

export async function getEdges() {
  const response = await fetchMetro<RequestEdge[]>(GET_EDGES_URL);
  if (!response.ok) {
    return [];
  }

  return await response.json();
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
  if (!response.ok) {
    return [];
  }

  return await response.json();
}

export async function getTrains() {
  const response = await fetchMetro<RequestTrain[]>(GET_TRAINS_URL);
  if (!response.ok) {
    return [];
  }

  return await response.json();
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

export async function createLine(sts: RequestCreateLine) {
  const response = await fetch(GET_LINES_URL, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(sts),
  });

  if (!response.ok) {
    throw new Error("couldn't create line");
  }
}

export async function updateTrainLine(trainId: number, lineId: number) {
  const response = await fetch(UPDATE_TRAIN_LINE_URL, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      trainId,
      lineId,
    }),
  });

  if (!response.ok) {
    throw new Error("couldn't move train");
  }
}
