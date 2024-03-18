const carCanvas = document.getElementById('carCanvas');
carCanvas.width = window.innerWidth - 330;
const networkCanvas = document.getElementById('networkCanvas');
networkCanvas.width = 300;
const miniMapCanvas = document.getElementById('miniMapCanvas');
miniMapCanvas.width = 300;
miniMapCanvas.height = 300;

carCanvas.height = window.innerHeight;
networkCanvas.height = window.innerHeight - 300;

const worldString = localStorage.getItem('world');
const worldInfo = worldString ? JSON.parse(worldString) : null;
let world = worldInfo ? World.load(worldInfo) : new World();
const graph = world.graph;

const viewPort = new Viewport(carCanvas, world.zoom, world.offset);
const miniMap = new MiniMap(miniMapCanvas, world.graph, 300);

const carCtx = carCanvas.getContext('2d');
const networkCtx = networkCanvas.getContext('2d');
const n = 100;
const cars = generateCars(n);
let bestCar = cars[0];

if (localStorage.getItem('bestBrain')) {
  for (let i = 0; i < cars.length; i++) {
    cars[i].brain = JSON.parse(localStorage.getItem('bestBrain'));
    if (i !== 0) {
      NeuralNetwork.mutate(cars[i].brain, 0.1);
    }
  }
}

const traffic = [];
const roadBorders = world.roadBorders.map(s => [s.p1, s.p2]);

animate();

function save() {
  localStorage.setItem('bestBrain', JSON.stringify(bestCar.brain))
}

function discard() {
  localStorage.removeItem('bestBrain');
}

function generateCars(n) {
  const startPoints = world.markings.filter(m => m.type === "start");
  const startPoint = startPoints.length > 0 ? startPoints[0].center : new Point(100, 100);
  const dir = startPoints.length > 0 ? startPoints[0].directionVector : new Point(0, -1);
  const startAngle = -angle(dir) + Math.PI / 2;

  const cars = [];
  for (let i = 1; i <= n; i++) {
    cars.push(new Car(startPoint.x, startPoint.y, 30, 50, "AI", startAngle));
  }
  return cars;
}

function animate(time) {
  // Resizing the canvas
  for (let i = 0; i < traffic.length; i++) {
    traffic[i].update(roadBorders, []);
  }
  for (let i = 0; i < cars.length; i++) {
    cars[i].update(roadBorders, traffic);
  }
  const lessMinFitness = Math.max(...cars.map(c => c.fitness))
  bestCar = cars.find(car => car.fitness == lessMinFitness);

  world.cars = cars;
  world.bestCar = bestCar;

  viewPort.offset.x = -bestCar.x;
  viewPort.offset.y = -bestCar.y;

  viewPort.reset();
  const viewPoint = Point.scale(viewPort.getOffset(), -1);
  world.draw(carCtx, viewPoint, false);

  miniMap.update(viewPoint);

  // Traffic draw
  for (let i = 0; i < traffic.length; i++) {
    traffic[i].draw(carCtx, "red");
  }

  networkCtx.lineDashOffset = - time / 50;
  networkCtx.clearRect(0, 0, networkCanvas.width, networkCanvas.height);
  Visualizer.drawNetwork(networkCtx, bestCar.brain);

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

function save() {
  localStorage.setItem('bestBrain', JSON.stringify(bestCar.brain))
}

function discard() {
  localStorage.removeItem('bestBrain');
}
