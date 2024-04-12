import Point from "../primitives/point.js";
import { Network } from "./network.js";

export class World {
  network: Network;
  zoom: number;
  offset: Point;

  constructor(graph: Network | null) {
    this.network = graph ? graph : new Network();
    this.zoom = 1;
    this.offset = new Point(0, 0);
  }

  /**
   * Set of intructions to run on each tick.
   * @param {CanvasRenderingContext2D} ctx
   * @param {Point} vp - mouse point
   */
  update(ctx: CanvasRenderingContext2D, vp: Point) {
    for (const node of this.network.nodes) {
      if (node.position.distanceTo(vp) < 60) {
        ctx.fillStyle = "white";
        ctx.font = "48px Arial";
        ctx.textAlign = "center";
        ctx.fillText(node.name, node.position.x, node.position.y - 50);
      }
    }
  }

  draw(ctx: CanvasRenderingContext2D) {
    for (const edge of this.network.edges) {
      edge.draw(ctx, { color: "white", dash: [], width: 1 });
    }
    for (const node of this.network.nodes) {
      node.draw(ctx);
    }
  }
}

export default World;
