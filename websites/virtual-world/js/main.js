const width = 700;
const height = 700;

const theCanvas = document.getElementById('theCanvas');

theCanvas.width = width;
theCanvas.height = height;

const ctx = theCanvas.getContext('2d');

const graphString = localStorage.getItem('graph');
const graphInfo = graphString ? JSON.parse(graphString) : null;
const graph = graphInfo ? Graph.load(graphInfo) : new Graph();
const viewPort = new Viewport(theCanvas);
const graphEditor = new GraphEditor(viewPort, graph);

animate();

function animate() {
  viewPort.reset();
  graphEditor.display();
  // new Polygon(graph.points).draw(ctx);
  new Envelope(graph.segments[0], 80).draw(ctx);

  requestAnimationFrame(animate);
}

function dispose() {
  graphEditor.dispose();
}

function save() {
  localStorage.setItem('graph', JSON.stringify(graph));
}