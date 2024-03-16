function getNearestPoint(loc, points, maxDist = Number.MAX_VALUE) {
  let minDist = Infinity;
  let nearest = null;
  for (const point of points) {
    const dist = point.distanceTo(loc);
    if (dist < minDist && dist < maxDist) {
      minDist = dist;
      nearest = point;
    }
  }
  return nearest;
}

function translate(loc, angle, offset) {
  return new Point(
    loc.x + Math.cos(angle) * offset,
    loc.y + Math.sin(angle) * offset 
  )
}

function angle(p) {
  return Math.atan2(p.y, p.x);
}

function normalize(p) {
  return Point.scale(p, 1 / magnitude(p));
}

function magnitude(p) {
  return Math.hypot(p.x, p.y);
}

function lerp(a, b, t) {
  return a + (b - a) * t;
}

function getIntersection(A, B, C, D) {
  const tTop = (D.x - C.x) * (A.y - C.y) - (D.y - C.y) * (A.x - C.x);
  const uTop = (C.y - A.y) * (A.x - B.x) - (C.x - A.x) * (A.y - B.y);
  const bottom = (D.y - C.y) * (B.x - A.x) - (D.x - C.x) * (B.y - A.y);

  const eps = 0.001;
  if (Math.abs(bottom) > eps) {
    const t = tTop / bottom;
    const u = uTop / bottom;
    if (t >= 0 && t <= 1 && u >= 0 && u <= 1) {
      return {
        x: lerp(A.x, B.x, t),
        y: lerp(A.y, B.y, t),
        offset: t,
      }
    }
  }

  return null
}

function getRandomColor() {
  const hue = 290 + Math.random() * 260;
  return `hsl(${hue}, 100%, 60%)`;
}

function average(p1, p2) {
  return new Point((p1.x + p2.x) / 2, (p1.y + p2.y) / 2);
}