export const CITY_BASE = "http://localhost:2221";
export const EVENTS_BASE_URL = "ws://localhost:2223";

export const GET_STATIONS_URL = `${CITY_BASE}/stations`;
export const GET_EDGES_URL = `${CITY_BASE}/edges`;
export const GET_EDGE_POINTS_URL = (id: number) => `${CITY_BASE}/edges/${id}`;

export const TRAINS_WS_FEED = `${EVENTS_BASE_URL}/trains`;
