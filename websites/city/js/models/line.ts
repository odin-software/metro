import Point from "../primitives/point.js";

export class Line {
  edgePoints: Point[];

  constructor(edgePoints: Point[]) {
    this.edgePoints = edgePoints;
  }

  draw(
    ctx: CanvasRenderingContext2D,
    { color = "yellow", dash = [22, 11], width = 2 } = {}
  ) {
    ctx.beginPath();
    ctx.lineWidth = width;
    ctx.strokeStyle = color;
    ctx.setLineDash(dash);
    ctx.moveTo(this.edgePoints[0].x, this.edgePoints[0].y);
    for (let i = 1; i < this.edgePoints.length; i++) {
      ctx.lineTo(this.edgePoints[i].x, this.edgePoints[i].y);
    }
    ctx.stroke();
  }
}
