import Point from "../primitives/point.js";
import Segment from "../primitives/segment.js";

class Graph {
  /**
   * Representation of a graph. It has points, segments and methods
   * to interact with those.
   * @constructor
   * @param {Point[]} [points=[]]
   * @param {Segment[]} [segments=[]]
   */
  constructor(points = [], segments = []) {
    this.points = points;
    this.segments = segments;
  }

  /**
   * @typedef GraphInfo
   * @property {Point[]} points
   * @property {Segment[]} segments
   */
  /**
   * Creates a graph out of passed points and segments.
   * @param {GraphInfo} info
   * @returns {Graph}
   */
  static load(info) {
    const points = info.points.map((p) => new Point(p.x, p.y, p.name));
    const segments = info.segments.map((s) => {
      return new Segment(
        points.find((p) => p.equals(s.p1)) || new Point(0, 0),
        points.find((p) => p.equals(s.p2)) || new Point(0, 0)
      );
    });
    return new Graph(points, segments);
  }

  /**
   * String representation of this graph.
   * @returns {string}
   */
  hash() {
    return JSON.stringify(this);
  }

  /**
   * Adds a point to the graph.
   * @param {Point} point
   */
  addPoint(point) {
    this.points.push(point);
  }

  /**
   * Checks if a point is contained within the graph.
   * @param {Point} point
   * @returns {boolean}
   */
  containsPoint(point) {
    return this.points.some((p) => p.equals(p));
  }

  /**
   * Tries to add a point in the graph if it doesn't exist.
   * @param {Point} point
   * @returns {boolean}
   */
  tryAddPoint(point) {
    if (this.containsPoint(point)) {
      return false;
    }

    this.addPoint(point);
    return true;
  }

  removePoint(point) {
    const segs = this.getSegmentsWithPoint(point);
    for (const seg of segs) {
      this.removeSegment(seg);
    }
    const index = this.points.indexOf(point);
    this.points.splice(index, 1);
  }

  addSegment(segment) {
    this.segments.push(segment);
  }

  containsSegment(segment) {
    return this.segments.find((s) => s.equals(segment));
  }

  tryAddSegment(segment) {
    if (this.containsSegment(segment)) {
      return false;
    }
    if (segment.p1.equals(segment.p2)) {
      return false;
    }

    this.segments.push(segment);
    return true;
  }

  removeSegment(segment) {
    const index = this.segments.indexOf(segment);
    this.segments.splice(index, 1);
  }

  getSegmentsWithPoint(point) {
    return this.segments.filter((seg) => seg.includes(point));
  }

  /**
   * Resets the graph by unloading the points and segments array.
   */
  dispose() {
    this.points.length = 0;
    this.segments.length = 0;
  }

  draw(ctx) {
    for (const seg of this.segments) {
      seg.draw(ctx);
    }

    for (const point of this.points) {
      point.draw(ctx, { size: 20, color: "white" });
    }
  }
}

export default Graph;
