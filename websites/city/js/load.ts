import { GET_EDGES_URL, GET_STATIONS_URL } from "./utils/consts.js";

export async function getStations() {
  const response = await fetch(GET_STATIONS_URL);
  const data: Station[] = await response.json();

  return data;
}

export async function getEdges() {
  const response = await fetch(GET_EDGES_URL);
  const data: Edge[] = await response.json();

  return data;
}
