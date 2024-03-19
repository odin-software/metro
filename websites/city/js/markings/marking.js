class Marking {
  constructor(center, directionVector, width, height) {
    this.center = center;
    this.directionVector = directionVector;
    this.width = width;
    this.height = height;

    this.support = new Segment(
      translate(center, angle(directionVector), height / 2),
      translate(center, angle(directionVector), -height / 2),
    );
    this.poly = new Envelope(
      this.support,
      width, 
      0
    ).poly;

    this.type = "marking";
  }

  static load(info) {
    const center = new Point(info.center.x, info.center.y);
    const dVec = new Point(info.directionVector.x, info.directionVector.y);
    switch (info.type) {
      case "crossing":
        return new Crossing(center, dVec, info.width, info.height);
      case "stop":
        return new Stop(center, dVec, info.width, info.height);
      case "light":
        return new Light(center, dVec, info.width, info.height);
      case "parking":
        return new Parking(center, dVec, info.width, info.height);
      case "target":
        return new Target(center, dVec, info.width, info.height);
      case "yield":
        return new Yield(center, dVec, info.width, info.height);
      case "start":
        return new Start(center, dVec, info.width, info.height);
      default:
        break;
    }
  }

  draw(ctx) {
    this.poly.draw(ctx);
  }
}