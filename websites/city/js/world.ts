import Graph from "./math/graph.js";
import Point from "./primitives/point.js";
import Segment from "./primitives/segment.js";
import { WorldInfo } from "./typings.js";

class World {
  graph: Graph;
  zoom: number;
  offset: Point;

  constructor(graph: Graph | null) {
    this.graph = graph ? graph : new Graph();
    this.zoom = 1;
    this.offset = new Point(0, 0);
  }

  /**
   * Loads a world and returns it.
   * @param {WorldInfo} info
   * @returns
   */
  static load(info: WorldInfo) {
    const world = new World(null);

    world.graph = Graph.load(info.graph);

    world.zoom = info.zoom;
    world.offset = info.offset;

    return world;
  }

  /**
   * Set of intructions to run on each tick.
   * @param {CanvasRenderingContext2D} ctx
   * @param {Point} vp - mouse point
   */
  update(ctx: CanvasRenderingContext2D, vp: Point) {
    for (const point of this.graph.points) {
      if (point.distanceTo(vp) < 60) {
        ctx.fillStyle = "white";
        ctx.font = "48px Arial";
        ctx.textAlign = "center";
        ctx.fillText(point.name, point.x, point.y - 50);
      }
    }
  }

  draw(ctx: CanvasRenderingContext2D) {
    for (const seg of this.graph.segments) {
      seg.draw(ctx, { color: "white", width: 2, dash: [30, 5] });
    }
    for (const point of this.graph.points) {
      point.draw(ctx, { color: "white", size: 30 });
    }
  }
}

export default World;
