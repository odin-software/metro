import Point from "../primitives/point.js";

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
   * Creates a draft station for the city editor.
   * @param x
   * @param y
   * @returns
   */
  static draft(x: number, y: number): Station {
    const name = (Math.random() * 2000).toString();
    const st = new Station(0, name, x, y);
    return st;
  }

  /**
   * Checks if a station is the same as this one.
   * @param station
   */
  equals(station: Station): boolean {
    return station.id == this.id;
  }

  /**
   * DRAFT: Checks if a station is the same as this one.
   * Checks only name since we don't have ID at draft point.
   * @param station
   */
  equalsDraft(station: Station): boolean {
    return station.name == this.name;
  }

  draw(ctx: CanvasRenderingContext2D, { color = "red", size = 19 } = {}) {
    this.position.draw(ctx, { color: color, size: size });
  }
}
