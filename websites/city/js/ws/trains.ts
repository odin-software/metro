import { Train } from "../models/train.js";
import Store from "../store/store.js";
import { TrainStore } from "../typings/store.js";
import { TRAINS_WS_FEED } from "../utils/consts.js";

export function initWs(trainStore: Store<TrainStore>) {
  const ws = new WebSocket(TRAINS_WS_FEED);

  ws.onmessage = (ev) => {
    const parsed = JSON.parse(ev.data) as RequestTrain[];
    trainStore.dispatch("updateAllTrains", {
      trains: Train.FromRequestToModel(parsed),
    });
  };

  ws.onclose = function (e) {
    console.warn(
      "Socket connection is closed. Reconnection attempt in 2 seconds",
      e.reason
    );
    setTimeout(function () {
      initWs(trainStore);
    }, 2000);
  };

  ws.onerror = function (err) {
    console.error("Socket encountered error, closing.");
    ws.close();
  };
}
