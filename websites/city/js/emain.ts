import { Dialog } from "./components/dialog.js";
import { NetworkEditor } from "./editors/networkEditor.js";
import { Network } from "./models/network.js";
import World from "./models/world.js";
import Point from "./primitives/point.js";
import Viewport from "./viewport.js";
import DialogStore from "./store/dialog.js";
import { saveDraftTemplate } from "./utils/template.js";
import { Line } from "./models/line.js";
import { getLines } from "./load.js";
import { LineEditor } from "./editors/lineEditor.js";

const canvas = document.getElementById("editorCanvas");
if (!canvas || !(canvas instanceof HTMLCanvasElement)) {
  throw new Error("the element #editorCanvas was not found");
}
canvas.width = window.innerWidth;
canvas.height = window.innerHeight - 150;

const dialog = new Dialog();

const ctx = canvas.getContext("2d");
const world = new World(await Network.load());
const viewport = new Viewport(
  canvas,
  world.zoom,
  world.network.getCenterPoint().invertSign()
);
const lines = await getLines();

const graphBtn = document.getElementById("graphBtn");
graphBtn.addEventListener("click", async () => {
  setMode("graph");
});
const lineBtn = document.getElementById("lineBtn");
lineBtn.addEventListener("click", async () => {
  setMode("line");
});

const saveBtn = document.getElementById("saveBtn");
saveBtn.addEventListener("click", async () => {
  DialogStore.commit("openDialog", {
    open: true,
    title: "Saving Drafts",
    body: saveDraftTemplate(
      world.network.draftNodes.length,
      world.network.draftEdges.length
    ),
    yesBtn: async () => {
      await world.network.saveDrafts();
      DialogStore.dispatch("closeDialog", {});
    },
    noBtn: () => DialogStore.dispatch("closeDialog", {}),
  });
});

const tools = {
  graph: {
    button: graphBtn,
    editor: new NetworkEditor(viewport, world.network),
  },
  line: {
    button: lineBtn,
    editor: new LineEditor(viewport, world.network),
  },
};

setMode("line");

animate();

function animate() {
  if (!ctx) {
    return;
  }
  viewport.reset();
  const viewPoint = Point.scale(viewport.getOffset(), -1);
  // world.draw(ctx, true);
  lines.forEach((ln) => {
    const line = new Line(ln.points.map((l) => new Point(l.x, l.y)));
    line.draw(ctx, { color: "yellow" });
  });

  ctx.globalAlpha = 0.5;
  for (const tool of Object.values(tools)) {
    tool.editor.display();
  }
  ctx.globalAlpha = 1;

  requestAnimationFrame(animate);
}

function dispose() {
  tools["graph"].editor.dispose();
}

function setMode(mode: string) {
  disableEditors();
  tools[mode].button.style.backgroundColor = "white";
  tools[mode].button.style.filter = "";
  tools[mode].editor.enable();
}

function disableEditors() {
  for (const tool of Object.values(tools)) {
    tool.button.style.backgroundColor = "gray";
    tool.button.style.filter = "grayscale(100%)";
    tool.editor.disable();
  }
}

// function openOsmPanel() {
//   osmPanel.style.display = "block";
// }

// function closeOsmPanel() {
//   osmPanel.style.display = "none";
// }

// function loadOsmData() {
//   if (osmDataContainer.value == "") {
//     alert("Please enter valid OSM data");
//     return;
//   }

//   const res = Osm.parseRoads(JSON.parse(osmDataContainer.value));
//   graph.points = res.points;
//   graph.segments = res.segments;
//   closeOsmPanel();
// }
