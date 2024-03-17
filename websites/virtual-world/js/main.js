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
const graphEditor = new GraphEditor(viewPort, graph);

let oldGraphHash = graph.hash();
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
  graphEditor.display();
  ctx.globalAlpha = 1;

  requestAnimationFrame(animate);
}

function dispose() {
  graphEditor.dispose();
}

function save() {
  localStorage.setItem('graph', JSON.stringify(graph));
}