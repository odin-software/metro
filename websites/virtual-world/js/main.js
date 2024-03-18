const width = 700;
const height = 700;

const theCanvas = document.getElementById('theCanvas');

theCanvas.width = width;
theCanvas.height = height;

const ctx = theCanvas.getContext('2d');

const worldString = localStorage.getItem('world');
const worldInfo = worldString ? JSON.parse(worldString) : null;
let world = worldInfo ? World.load(worldInfo) : new World();
const graph = world.graph;

const viewPort = new Viewport(theCanvas, world.zoom, world.offset);
const tools = {
  graph: { button: graphBtn, editor: new GraphEditor(viewPort, graph) },
  stop: { button: stopBtn, editor: new StopEditor(viewPort, world) },
  crossing: { button: crossingBtn, editor: new CrossingEditor(viewPort, world) },
  start: { button: startBtn, editor: new StartEditor(viewPort, world) },
  parking: { button: parkingBtn, editor: new ParkingEditor(viewPort, world) },
  light: { button: lightBtn, editor: new LightEditor(viewPort, world) },
  target: { button: targetBtn, editor: new TargetEditor(viewPort, world) },
  yield: { button: yieldBtn, editor: new YieldEditor(viewPort, world) },
}

let oldGraphHash = graph.hash();

setMode('graph');

animate();

function animate() {
  viewPort.reset();
  if (graph.hash() != oldGraphHash) {
    world.generate();
    oldGraphHash = graph.hash();
  }
  const viewPoint = Point.scale(viewPort.getOffset(), -1);
  world.draw(ctx, viewPoint);

  ctx.globalAlpha = 0.5;
  for (const tool of Object.values(tools)) {
    tool.editor.display();
  }
  ctx.globalAlpha = 1;

  requestAnimationFrame(animate);
}

function dispose() {
  tools["graph"].editor.dispose();
  world.markings.length = 0;
}

function save() {
  world.zoom = viewPort.zoom;
  world.offset = viewPort.getOffset();

  const element = document.createElement('a');
  element.setAttribute("href", "data:application/json;charset=utf-8," + encodeURIComponent(JSON.stringify(world)));

  const fileName = "name.world";
  element.setAttribute("download", fileName);

  element.click();

  localStorage.setItem('world', JSON.stringify(world));
}

function load(event) {
  const file = event.target.files[0];
  if (!file) return;

  const reader = new FileReader();
  reader.readAsText(file);
  reader.onload = function (event) {
    const fileContent = event.target.result;
    const worldInfo = JSON.parse(fileContent);
    world = World.load(worldInfo);

    localStorage.setItem('world', JSON.stringify(world));
    location.reload();
  };
}

function setMode(mode) {
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