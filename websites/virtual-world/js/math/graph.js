class Graph {
  constructor(points = [], segments = []) {
    this.points = points;
    this.segments = segments;
  }

  static load(info) {
    const points = info.points.map(p => new Point(p.x, p.y));
    const segments = info.segments.map(s => {
      return new Segment(
        points.find(p => p.equals(s.p1)),
        points.find(p => p.equals(s.p2))
      );
    })
    return new Graph(points, segments);
  }

  addPoint(point) {
    this.points.push(point);
  }

  containsPoint(point) {
    return this.points.find(p => p.equals(point));
  }

  tryAddPoint(point) {
    if (this.containsPoint(point)) {
      return false;
    }

    this.addPoint(point);
    return true;
  }

  removePoint(point) {
    const segs = this.getSegmentsWithPoint(point);
    for (const seg of segs) {
      this.removeSegment(seg);
    }
    const index = this.points.indexOf(point);
    this.points.splice(index, 1);
  }

  addSegment(segment) {
    this.segments.push(segment);
  }

  containsSegment(segment) {
    return this.segments.find(s => s.equals(segment));
  }

  tryAddSegment(segment) {
    if (this.containsSegment(segment)) {
      return false;
    }
    if (segment.p1.equals(segment.p2)) {
      return false;
    }

    this.segments.push(segment);
    return true;
  }

  removeSegment(segment) {
    const index = this.segments.indexOf(segment);
    this.segments.splice(index, 1);
  }

  getSegmentsWithPoint(point) {
    return this.segments.filter(seg => seg.includes(point));
  }

  dispose() {
    this.points.length = 0;
    this.segments.length = 0;
  }

  draw(ctx) {
    for (const seg of this.segments) {
      seg.draw(ctx);
    }

    for (const point of this.points) {
      point.draw(ctx);
    }
  }
}