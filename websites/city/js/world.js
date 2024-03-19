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
      seg.draw(ctx, { color: "white", width: 2, dash: [30, 5]});
    }
    for (const point of this.graph.points) {
      point.draw(ctx, { color: "white", size: 30 });
    }
  }
}
