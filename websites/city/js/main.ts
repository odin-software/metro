import Viewport from "./viewport.js";
import World from "./world.js";
import Point from "./primitives/point.js";
import Graph from "./math/graph.js";

import { initWs } from "./ws/trains.js";
import { getStations } from "./load.js";

const canvas = document.getElementById("cityCanvas");
if (!canvas || !(canvas instanceof HTMLCanvasElement)) {
  throw new Error("the element #cityCanvas was not found");
}
canvas.width = window.innerWidth;
canvas.height = window.innerHeight;

const ctx = canvas.getContext("2d");
const stations = await getStations();
const world = new World(
  new Graph(
    stations.map((st) => new Point(st.position.x, st.position.y, st.name))
  )
);

const viewport = new Viewport(canvas, world.zoom, world.offset);

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
  stations.forEach((st) => {
    const p = new Point(st.position.x, st.position.y);
    p.draw(ctx, { size: 14, color: "white" });
  });
  trains.forEach((tr) => {
    const p = new Point(tr.x, tr.y);
    p.draw(ctx, { size: 30, color: "white" });
  });

  requestAnimationFrame(animate);
}

export function setTrains(wsTrains) {
  trains = wsTrains;
}
