import Point from "../primitives/point.js";

export function getNearestPoint(loc, points, maxDist = Number.MAX_VALUE) {
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

export function getNearestSegment(loc, segments, maxDist = Number.MAX_VALUE) {
  let minDist = Infinity;
  let nearest = null;
  for (const seg of segments) {
    const dist = seg.distanceToPoint(loc);
    if (dist < minDist && dist < maxDist) {
      minDist = dist;
      nearest = seg;
    }
  }
  return nearest;
}

export function perpendicular(p: Point) {
  return new Point(-p.y, p.x);
}

export function translate(loc, angle, offset) {
  return new Point(
    loc.x + Math.cos(angle) * offset,
    loc.y + Math.sin(angle) * offset
  );
}

/**
 * Returns the angle of a Point, based on the X axis.
 * @param {Point} p
 * @returns {number}
 */
export function angle(p: Point): number {
  return Math.atan2(p.y, p.x);
}

/**
 * Returns a normalized vector. A normalized vector is a vector with
 * the same direction but a magnitude of 1.
 * @param {Point} p
 * @returns {Point}
 */
export function normalize(p: Point): Point {
  return Point.scale(p, 1 / magnitude(p));
}

/**
 * Returns the magnitude of a vector, this is done by calculating
 * the squared root of the sum of the squared values of a point.
 * @param {Point} p
 * @returns {number}
 */
export function magnitude(p: Point): number {
  return Math.hypot(p.x, p.y);
}

/**
 * The dot function is used to check if two points are:
 * - answer is more than 0, they are facing in somewhat the same direction
 * - answer is less than 0, they are facing opposite directions
 * - answer is 0, they are perpendicular to one another
 * With unit vectors what happens is:
 * - if answer is -1, they are facing opposite directions
 * - if answer is 1, they are facing the exact direction
 * - if answer is 0, they are perpendicular
 * @param {Point} p1
 * @param {Point} p2
 * @returns {number} answer
 */
export function dot(p1: Point, p2: Point): number {
  return p1.x * p2.x + p1.y * p2.y;
}

export function lerp(a, b, t) {
  return a + (b - a) * t;
}

export function lerp2D(A, B, t) {
  return new Point(lerp(A.x, B.x, t), lerp(A.y, B.y, t));
}

export function invLerp(a, b, x) {
  return (x - a) / (b - a);
}

export function degToRad(deg) {
  return (deg * Math.PI) / 180;
}

export function getIntersection(A, B, C, D) {
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
      };
    }
  }

  return null;
}

export function getRandomColor() {
  const hue = 290 + Math.random() * 260;
  return `hsl(${hue}, 100%, 60%)`;
}

export function average(p1, p2) {
  return new Point((p1.x + p2.x) / 2, (p1.y + p2.y) / 2);
}

export function getFake3DPoint(point, viewPoint, height) {
  const dir = normalize(Point.sub(point, viewPoint));
  const dist = point.distanceTo(viewPoint);
  const scaler = Math.atan(dist / 300) / (Math.PI / 2);
  return Point.add(point, Point.scale(dir, height * scaler));
}

export function polysIntersect(poly1, poly2) {
  for (let i = 0; i < poly1.length; i++) {
    const next = (i + 1) % poly1.length;
    for (let j = 0; j < poly2.length; j++) {
      const next2 = (j + 1) % poly2.length;
      if (getIntersection(poly1[i], poly1[next], poly2[j], poly2[next2])) {
        return true;
      }
    }
  }
  return false;
}

export function getRGBA(value) {
  const alpha = Math.abs(value);
  const R = value < 0 ? 0 : 255;
  const G = R;
  const B = value > 0 ? 0 : 255;
  return `rgba(${R},${G},${B},${alpha})`;
}
