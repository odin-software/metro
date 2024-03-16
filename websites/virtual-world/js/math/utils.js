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