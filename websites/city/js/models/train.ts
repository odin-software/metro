import Point from "../primitives/point.js";

export class Train {
  id: number;
  name: string;
  position: Point;
  currentStation: number;
  line: string;
  make: string;

  constructor(
    id: number,
    name: string,
    x: number,
    y: number,
    curr: number,
    line: string,
    make: string
  ) {
    this.id = id;
    this.name = name;
    this.position = new Point(x, y);
    this.currentStation = curr;
    this.line = line;
    this.make = make;
  }

  draw(ctx: CanvasRenderingContext2D) {
    this.position.draw(ctx, { size: 24, color: "white" });
  }
}
