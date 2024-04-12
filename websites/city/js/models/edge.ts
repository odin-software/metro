import Point from "../primitives/point.js";
import { Station } from "./station";

export class Edge {
  start: Station;
  end: Station;
  edgePoints: Point[];

  constructor(start: Station, end: Station, edgePoints: Point[]) {
    this.start = start;
    this.end = end;
    this.edgePoints = edgePoints;
  }

  /**
   * Checks if a edge has the same points as this one.
   * @param {Edge} edge
   * @returns {boolean}
   */
  equals(edge: Edge): boolean {
    return this.includes(edge.start) && this.includes(edge.end);
  }

  /**
   * Returns true if a node is included in one of this edge's point.
   * @param {Station} node
   * @returns {boolean}
   */
  includes(node: Station): boolean {
    return this.start.equals(node) || this.end.equals(node);
  }

  draw(
    ctx: CanvasRenderingContext2D,
    { color = "white", dash = [22, 11], width = 2 }: EdgeStyling | undefined
  ) {
    ctx.beginPath();
    ctx.lineWidth = width;
    ctx.strokeStyle = color;
    ctx.setLineDash(dash);
    ctx.moveTo(this.start.position.x, this.start.position.y);
    for (let i = 0; i < this.edgePoints.length; i++) {
      ctx.lineTo(this.edgePoints[i].x, this.edgePoints[i].y);
    }
    ctx.lineTo(this.end.position.x, this.end.position.y);
    ctx.stroke();
  }
}

type EdgeStyling = {
  color: string;
  dash: number[];
  width: number;
};
