class World {
  constructor(graph) {
    this.graph = graph ? graph : new Graph();
  }

  static load(info) {
    const world = new World(new Graph());

    world.graph = Graph.load(info.graph);

    world.zoom = info.zoom;
    world.offset = info.offset;
        
    return world;
  }

  draw(ctx, viewPoint) {
    for (const seg of this.graph.segments) {
      seg.draw(ctx, { color: "white", width: 3 });
    }
    for (const p of this.graph.points) {
      p.draw(ctx, { color: "white", radius: 12 });
    }
  }
}
