import Point from "../primitives/point.js";
import { Network } from "./network.js";

export class World {
  network: Network;
  zoom: number;
  offset: Point;

  constructor(network: Network | null) {
    this.network = network ? network : new Network();
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
        ctx.font = "28px Arial";
        ctx.textAlign = "center";
        ctx.fillText(node.name, node.position.x, node.position.y - 25);
      }
    }
  }

  draw(ctx: CanvasRenderingContext2D, draft = false) {
    for (const edge of this.network.edges) {
      edge.draw(ctx, { color: "white", dash: [], width: 1 });
    }
    for (const node of this.network.nodes) {
      node.draw(ctx);
    }
    if (draft) {
      for (const edge of this.network.draftEdges) {
        edge.draw(ctx, { color: "white", dash: [], width: 1 });
      }
      for (const node of this.network.draftNodes) {
        node.draw(ctx);
      }
    }
  }
}

export default World;
