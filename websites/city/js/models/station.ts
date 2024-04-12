import Point from "../primitives/point";

export class Station {
  id: number;
  name: string;
  position: Point;
  trains: number[]; // ids of trains in the station.
  createdAt?: Date;
  updatedAt?: Date;

  constructor(id: number, name: string, x: number, y: number) {
    this.id = id;
    this.name = name;
    this.position = new Point(x, y);
    this.trains = [];
  }

  /**
   * Checks if a station is the same as this one.
   * @param station
   */
  equals(station: Station): boolean {
    return station.id == this.id;
  }

  draw(ctx: CanvasRenderingContext2D) {
    this.position.draw(ctx, {});
  }
}
