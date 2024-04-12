import Viewport from "./viewport.js";
import { World } from "./models/world.js";
import Point from "./primitives/point.js";
import Graph from "./math/graph.js";

import { initWs } from "./ws/trains.js";
import { getEdges, getEdgesPoints, getStations } from "./load.js";
import Segment from "./primitives/segment.js";
import { Network } from "./models/network.js";
import { Station } from "./models/station.js";
import { Edge } from "./models/edge.js";

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

const mouse = new Point(0, 0);

canvas.addEventListener("mousemove", (event) => {
  mouse.x = event.clientX;
  mouse.y = event.clientY;
});

initWs();
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
    p.draw(ctx, { size: 30, color: "white" });
  });

  requestAnimationFrame(animate);
}

export function setTrains(wsTrains) {
  trains = wsTrains;
}
