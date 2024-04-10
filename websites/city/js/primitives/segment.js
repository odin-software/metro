import { dot, magnitude, normalize } from "../math/utils.js";
import Point from "./point.js";

class Segment {
  constructor(p1, p2, oneWay = false) {
    this.p1 = p1;
    this.p2 = p2;
    this.oneWay = oneWay;
  }

  length() {
    return this.p1.distanceTo(this.p2);
  }

  directionVector() {
    return normalize(Point.sub(this.p2, this.p1));
  }

  equals(other) {
    return this.includes(other.p1) && this.includes(other.p2);
  }

  includes(point) {
    return this.p1.equals(point) || this.p2.equals(point);
  }

  distanceToPoint(point) {
    const proj = this.projectPoint(point);
    if (proj.offset > 0 && proj.offset < 1) {
      return proj.point.distanceTo(point);
    }
    const distToP1 = point.distanceTo(this.p1);
    const distToP2 = point.distanceTo(this.p2);
    return Math.min(distToP1, distToP2);
  }

  projectPoint(point) {
    const a = Point.sub(point, this.p1);
    const b = Point.sub(this.p2, this.p1);
    const normB = normalize(b);
    const scaler = dot(a, normB);
    const proj = {
      point: Point.add(this.p1, Point.scale(normB, scaler)),
      offset: scaler / magnitude(b),
    };
    return proj;
  }

  draw(ctx, { width = 2, color = "white", dash = [] } = {}) {
    ctx.beginPath();
    ctx.lineWidth = width;
    ctx.strokeStyle = color;
    ctx.setLineDash(dash);
    ctx.moveTo(this.p1.x, this.p1.y);
    ctx.lineTo(this.p2.x, this.p2.y);
    ctx.stroke();
  }
}

export default Segment;
