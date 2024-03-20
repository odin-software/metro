const canvas = document.getElementById('theCanvas');
canvas.width = window.innerWidth;
canvas.height = window.innerHeight;

const ctx = canvas.getContext('2d');

const worldString = localStorage.getItem('world');
const worldInfo = worldString ? JSON.parse(worldString) : null;
let world = worldInfo ? World.load(worldInfo) : new World(new Graph());
const graph = world.graph;

const viewPort = new Viewport(canvas, world.zoom, world.offset);
const mouse = new Point(0, 0);

canvas.addEventListener('mousemove', (event) => {
  mouse.x = event.clientX;
  mouse.y = event.clientY;
});

animate();

function animate() {
  viewPort.reset();
  const gm = viewPort.getMouseFromPoint(mouse);
  world.update(ctx, gm);
  world.draw(ctx);

  requestAnimationFrame(animate);
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
