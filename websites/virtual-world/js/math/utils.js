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