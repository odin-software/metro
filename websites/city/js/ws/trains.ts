import { setTrains } from "../main.js";
import { TRAINS_WS_FEED } from "../utils/consts.js";

export function initWs() {
  const ws = new WebSocket(TRAINS_WS_FEED);

  ws.onmessage = (ev) => {
    const parsed = JSON.parse(ev.data);
    setTrains(parsed);
  };

  ws.onclose = function (e) {
    console.warn(
      "Socket connection is closed. Reconnection attempt in 2 seconds",
      e.reason
    );
    setTimeout(function () {
      initWs();
    }, 2000);
  };

  ws.onerror = function (err) {
    console.error("Socket encountered error, closing.");
    ws.close();
  };
}
