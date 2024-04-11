import { setTrains } from "../main.js";

export function initWs() {
  const ws = new WebSocket("ws://localhost:2223/trains");

  ws.onmessage = (ev) => {
    const parsed = JSON.parse(ev.data);
    setTrains(parsed);
  };
}
