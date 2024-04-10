import Graph from "./math/graph.js";
import Point from "./primitives/point.js";

class World {
  constructor(graph) {
    this.graph = graph ? graph : new Graph();
    this.zoom = 1;
    this.offset = new Point(0, 0);
  }

  static load(info) {
    const world = new World(new Graph());

    world.graph = Graph.load(info.graph);

    world.zoom = info.zoom;
    world.offset = info.offset;

    return world;
  }

  update(ctx, viewPoint) {
    for (const point of this.graph.points) {
      if (point.distanceTo(viewPoint) < 60) {
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
