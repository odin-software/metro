const canvas = document.getElementById('myCanvas');
canvas.width = 200;

const ctx = canvas.getContext('2d');
const road = new Road(canvas.width / 2, canvas.width * 0.9);
const car = new Car(road.getLaneCenter(1), 100, 30, 50);

animate();

function animate() {
  // Resizing the canvas
  ctx.clearRect(0, 0, canvas.width, canvas.height);
  canvas.height = window.innerHeight;
  car.update(road.borders);

  ctx.save(); 
  ctx.translate(0, -car.y + canvas.height *0.7);

  // Road draw
  road.draw(ctx);

  // Car movement and draw
  car.draw(ctx);

  ctx.restore();
  requestAnimationFrame(animate);
}