const width = 700;
const height = 700;

const theCanvas = document.getElementById('theCanvas');

theCanvas.width = width;
theCanvas.height = height;

const ctx = theCanvas.getContext('2d');

const graphString = localStorage.getItem('graph');
const graphInfo = graphString ? JSON.parse(graphString) : null;
const graph = graphInfo ? Graph.load(graphInfo) : new Graph();
const world = new World(graph);

const viewPort = new Viewport(theCanvas);
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
  localStorage.setItem('graph', JSON.stringify(graph));
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