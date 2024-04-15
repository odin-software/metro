import { setLogs } from "../main.js";
import { LOGS_WS_FEED } from "../utils/consts.js";

export function initEventsWs() {
  const ws = new WebSocket(LOGS_WS_FEED);

  ws.onmessage = (ev) => {
    const parsed = JSON.parse(ev.data);
    setLogs(parsed);
  };

  ws.onclose = function (e) {
    console.warn(
      "Socket connection is closed. Reconnection attempt in 2 seconds",
      e.reason
    );
    setTimeout(function () {
      initEventsWs();
    }, 2000);
  };

  ws.onerror = function (err) {
    console.error("Socket encountered error, closing.");
    ws.close();
  };
}
