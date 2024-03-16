class World {
  constructor(
    graph, 
    roadWidth = 100, 
    roadRoundness = 8,
    buildingWidth = 150,
    buildingMinLength = 150,
    spacing = 50
  ) {
    this.graph = graph;
    this.roadWidth = roadWidth;
    this.roadRoundness = roadRoundness;
    this.buildingWidth = buildingWidth;
    this.buildingMinLength = buildingMinLength;
    this.spacing = spacing;

    this.envelopes = [];
    this.roadBorders = [];
    this.buildings = [];

    this.generate();
  }

  generate() {
    this.envelopes.length = 0;
    for (const seg of this.graph.segments) {
      this.envelopes.push(
        new Envelope(
          seg, this.roadWidth, this.roadRoundness
        )
      );
    }

    this.roadBorders = Polygon.union(this.envelopes.map(env => env.poly));
    this.buildings = this.#generateBuildings();
  }

  #generateBuildings() {
    const tmpEnvelopes = [];
    for (const seg of this.graph.segments) {
      tmpEnvelopes.push(
        new Envelope(
          seg, 
          this.roadWidth + this.buildingWidth + this.spacing * 2,
          this.roadRoundness
        )
      );
    }

    const guides = Polygon.union(tmpEnvelopes.map(env => env.poly));

    for (let i = 0; i < guides.length; i++) {
      const seg = guides[i];
      if (seg.length() < this.buildingMinLength) {
        guides.splice(i, 1);
        i--;
      }
    }

    const supports = [];
    for (const seg of guides) {
      const len = seg.length() + this.spacing; 
      const buildingCount = Math.floor(len / (this.buildingMinLength + this.spacing));
      const buildingLength = len / buildingCount - this.spacing;

      const dir = seg.directionVector();

      let q1 = seg.p1;
      let q2 = Point.add(q1, Point.scale(dir, buildingLength));
      supports.push(new Segment(q1, q2));

      for (let i = 2; i <= buildingCount; i++) {
        q1 = Point.add(q2, Point.scale(dir, this.spacing));
        q2 = Point.add(q1, Point.scale(dir, buildingLength));
        supports.push(new Segment(q1, q2));
      }
    }

    const bases = [];
    for (const seg of supports) {
      bases.push(
        new Envelope(
          seg, this.buildingWidth
        ).poly
      );
    }

    for (let i = 0; i < bases.length - 1; i++) {
      for (let j = i + 1; j < bases.length; j++) {
        const base1 = bases[i];
        const base2 = bases[j];
        if (base1.intersectsPoly(base2)) {
          bases.splice(j, 1);
          j--;
        }
      }
    }

    return bases;
  }

  draw(ctx) {
    for (const env of this.envelopes) {
      env.draw(ctx, { fill: "#BBB", stroke: "#BBB", lineWidth: 15 });
    }
    for (const seg of this.graph.segments) {
      seg.draw(ctx, { color: "white", width: 4, dash: [10, 10] });
    }
    for (const seg of this.roadBorders) {
      seg.draw(ctx, { color: "white", width: 4 });
    }
    for (const building of this.buildings) {
      building.draw(ctx);
    }
  }
}
