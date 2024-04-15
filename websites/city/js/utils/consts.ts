export const CITY_BASE = "http://localhost:2221";
export const EVENTS_BASE_URL = "ws://localhost:2223";

export const GET_STATIONS_URL = `${CITY_BASE}/stations`;
export const GET_EDGES_URL = `${CITY_BASE}/edges`;
export const GET_EDGE_POINTS_URL = (id: number) => `${CITY_BASE}/edges/${id}`;

export const GET_PLAY_LOOP = `${CITY_BASE}/resume`;
export const GET_PAUSE_LOOP = `${CITY_BASE}/pause`;

export const TRAINS_WS_FEED = `${EVENTS_BASE_URL}/trains`;
export const LOGS_WS_FEED = `${EVENTS_BASE_URL}/logs`;
