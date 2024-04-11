import Point from "./point.js";
import { dot, magnitude, normalize } from "../math/utils.js";
import { SegmentStyle } from "../typings.js";

/**
 * The Segment class represents a connection between two points inside a graph.
 */
class Segment {
  p1: Point;
  p2: Point;
  oneWay: boolean;

  constructor(p1: Point, p2: Point, oneWay: boolean = false) {
    this.p1 = p1;
    this.p2 = p2;
    this.oneWay = oneWay;
  }

  /**
   * Length of the segment, which is the distance between the two points connected.
   * @returns {number}
   */
  length(): number {
    return this.p1.distanceTo(this.p2);
  }

  /**
   * The direction vector of the two points of this segment.
   * @returns {Point}
   */
  directionVector(): Point {
    return normalize(Point.sub(this.p2, this.p1));
  }

  /**
   * Checks if a segment has the same points as this one.
   * @param {Segment} seg
   * @returns {boolean}
   */
  equals(seg: Segment): boolean {
    return this.includes(seg.p1) && this.includes(seg.p2);
  }

  /**
   * Returns if a point is included in one of this segment's point.
   * @param {Point} point
   * @returns {boolean}
   */
  includes(point: Point): boolean {
    return this.p1.equals(point) || this.p2.equals(point);
  }

  /**
   * Distance to a point.
   * @param {Point} point
   * @returns {number}
   */
  distanceToPoint(point: Point): number {
    const proj = this.projectPoint(point);
    if (proj.offset > 0 && proj.offset < 1) {
      return proj.point.distanceTo(point);
    }
    const distToP1 = point.distanceTo(this.p1);
    const distToP2 = point.distanceTo(this.p2);
    return Math.min(distToP1, distToP2);
  }

  /**
   * Project a point in this segment.
   * @param {Point} point
   * @returns
   */
  projectPoint(point: Point) {
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

  /**
   * Styling options for drawing a segment.
   * @typedef {Object} Styles
   * @property {number} [width=2] - with of the segment
   * @property {string} [color=white] - color of the segment
   * @property {number[]} [dash=[]] - defines the dash style of the segment
   */
  /**
   * Function to draw a Segment with options on styling.
   * @param {CanvasRenderingContext2D} ctx
   * @param {SegmentStyle} style
   */
  draw(
    ctx: CanvasRenderingContext2D,
    { width = 2, color = "white", dash = [] }: SegmentStyle
  ) {
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
