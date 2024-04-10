import Graph from "./math/graph.js";
import Point from "./primitives/point.js";
import Segment from "./primitives/segment.js";

class World {
  /**
   * Representation of the world. It has the graph contained, also the initial
   * zoom.
   * @constructor
   * @param {?Graph} [graph]
   */
  constructor(graph) {
    this.graph = graph ? graph : new Graph();
    this.zoom = 1;
    this.offset = new Point(0, 0);
  }

  /**
   * @typedef GraphInfo
   * @property {Point[]} points
   * @property {Segment[]} segments
   */
  /**
   * @typedef WorldInfo
   * @property {GraphInfo} graph
   * @property {number} zoom
   * @property {Point} offset
   */
  /**
   * Loads a world and returns it.
   * @param {WorldInfo} info
   * @returns
   */
  static load(info) {
    const world = new World();

    world.graph = Graph.load(info.graph);

    world.zoom = info.zoom;
    world.offset = info.offset;

    return world;
  }

  /**
   * Set of intructions to run on each tick.
   * @param {CanvasRenderingContext2D} ctx
   * @param {*} vp - mouse point
   */
  update(ctx, vp) {
    for (const point of this.graph.points) {
      if (point.distanceTo(vp) < 60) {
        ctx.fillStyle = "white";
        ctx.font = "48px Arial";
        ctx.fillText(point.name, point.x - 140, point.y - 50);
      }
    }
  }

  draw(ctx) {
    for (const seg of this.graph.segments) {
      seg.draw(ctx, { color: "white", width: 2, dash: [30, 5] });
    }
    for (const point of this.graph.points) {
      point.draw(ctx, { color: "white", size: 30 });
    }
  }
}

export default World;
