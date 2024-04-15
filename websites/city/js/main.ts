import Viewport from "./viewport.js";
import { World } from "./models/world.js";
import Point from "./primitives/point.js";

import { initWs } from "./ws/trains.js";
import { Network } from "./models/network.js";
import { pauseLoop, playLoop } from "./load.js";
import { initEventsWs } from "./ws/events.js";

const canvas = document.getElementById("cityCanvas");
if (!canvas || !(canvas instanceof HTMLCanvasElement)) {
  throw new Error("the element #cityCanvas was not found");
}
canvas.width = window.innerWidth;
canvas.height = window.innerHeight;

const ctx = canvas.getContext("2d");
const world = new World(await Network.load());
const viewport = new Viewport(
  canvas,
  world.zoom,
  world.network.getCenterPoint().invertSign()
);

let trains = [];
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

initWs();
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
  trains.forEach((tr) => {
    const p = new Point(tr.x, tr.y);
    p.draw(ctx, { size: 24, color: "white" });
  });

  requestAnimationFrame(animate);
}

export function setTrains(wsTrains) {
  trains = wsTrains;
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
