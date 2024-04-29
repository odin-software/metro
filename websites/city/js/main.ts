import Viewport from "./viewport.js";
import { World } from "./models/world.js";
import Point from "./primitives/point.js";

import Store from "./store/store.js";
import { initWs } from "./ws/trains.js";
import { Network } from "./models/network.js";
import { getLines, pauseLoop, playLoop } from "./load.js";
import { initEventsWs } from "./ws/events.js";
import { Line } from "./models/line.js";
import { TrainStore } from "./typings/store.js";
import { trainStoreParams } from "./store/trains.js";

const canvas = document.getElementById("cityCanvas");
if (!canvas || !(canvas instanceof HTMLCanvasElement)) {
  throw new Error("the element #cityCanvas was not found");
}
canvas.width = window.innerWidth;
canvas.height = window.innerHeight;

const trainStore = new Store<TrainStore>(trainStoreParams);

const ctx = canvas.getContext("2d");
const world = new World(await Network.load());
const viewport = new Viewport(
  canvas,
  world.zoom,
  world.network.getCenterPoint().invertSign()
);
const lines = await getLines();
console.log(lines);

let logs = [];
const mouse = new Point(0, 0);
const list = document.querySelector("#logsList") as HTMLUListElement;

canvas.addEventListener("mousemove", (event) => {
  mouse.x = event.clientX;
  mouse.y = event.clientY;
});
document.querySelector("#pauseBtn").addEventListener("click", async () => {
  pauseLoop();
});
document.querySelector("#playBtn").addEventListener("click", async () => {
  playLoop();
});

initWs(trainStore);
initEventsWs();
animate();

function animate() {
  if (!ctx) {
    return;
  }
  viewport.reset();
  const gm = viewport.getMouseFromPoint(mouse);
  world.update(ctx, gm);
  world.draw(ctx);
  if (trainStore.state.trains.length > 0) {
    trainStore.state.trains.forEach((tr) => {
      tr.draw(ctx);
    });
  }
  lines.forEach((ln) => {
    const line = new Line(ln.points.map((l) => new Point(l.x, l.y)));
    line.draw(ctx, { color: "yellow" });
  });

  requestAnimationFrame(animate);
}

export function setLogs(wsLogs) {
  logs = wsLogs;
  list.innerHTML = "";

  for (let i = 0; i < logs.length; i++) {
    const li = document.createElement("li");
    const p = document.createTextNode(logs[i]);
    li.appendChild(p);
    list.appendChild(li);
  }
}
