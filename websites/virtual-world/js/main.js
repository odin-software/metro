const width = 700;
const height = 700;

function addRandomPoint() {
  const success = graph.tryAddPoint(new Point(Math.random() * width, Math.random() * height));
  ctx.clearRect(0, 0, width, height);
  graph.draw(ctx);
}

function addRandomSegment() {
  const index1 = Math.floor(Math.random() * graph.points.length); 
  const index2 = Math.floor(Math.random() * graph.points.length); 
  const success = graph.tryAddSegment(new Segment(graph.points[index1], graph.points[index2]));

  ctx.clearRect(0, 0, width, height);
  graph.draw(ctx);
}

function removeRandomSegment() {
  if (graph.segments.length === 0) {
    console.log('No segments to remove');
    return;
  }

  const index = Math.floor(Math.random() * graph.segments.length);
  graph.removeSegment(graph.segments[index]);

  ctx.clearRect(0, 0, width, height);
  graph.draw(ctx);
}

function removeRandomPoint() {
  if (graph.points.length === 0) {
    console.log('No points to remove');
    return;
  }

  const index = Math.floor(Math.random() * graph.points.length);
  graph.removePoint(graph.points[index]);

  ctx.clearRect(0, 0, width, height);
  graph.draw(ctx);
}

function removeAll() {
  graph.dispose();

  ctx.clearRect(0, 0, width, height);
  graph.draw(ctx);
}

const theCanvas = document.getElementById('theCanvas');

theCanvas.width = width;
theCanvas.height = height;

const ctx = theCanvas.getContext('2d');

const p1 = new Point(200, 200);
const p2 = new Point(500, 200);
const p3 = new Point(400, 400);
const p4 = new Point(100, 300);

const s1 = new Segment(p1, p2);
const s2 = new Segment(p1, p3);
const s3 = new Segment(p1, p4);
const s4 = new Segment(p2, p3);

const graph = new Graph([p1, p2, p3, p4], [s1, s2, s3, s4]);

graph.draw(ctx);