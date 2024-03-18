class Light extends Marking {
  constructor(center, directionVector, width, height) {
    super(center, directionVector, width, 18);

    this.state = "green";
    this.border = this.poly.segments[0];
    this.type = "light";
  }

  draw(ctx) {
    const perp = perpendicular(this.directionVector);
    const line = new Segment(
      Point.add(this.center, Point.scale(perp, this.width / 2)),
      Point.add(this.center, Point.scale(perp, -this.width / 2)),
    )

    const green = lerp2D(line.p1, line.p2, 0.2);
    const yellow = lerp2D(line.p1, line.p2, 0.5);
    const red = lerp2D(line.p1, line.p2, 0.8);

    new Segment(red, green).draw(ctx, {
      width: this.height,
      cap: "round"
    });

    green.draw(ctx, { color: "#060", size: this.height * 0.6 });
    yellow.draw(ctx, { color: "#660", size: this.height * 0.6 });
    red.draw(ctx, { color: "#600", size: this.height * 0.6 });
    
    switch (this.state) {
      case "green":
        green.draw(ctx, { color: "#0f0", size: this.height * 0.6 });
        break;
      case "yellow":
        yellow.draw(ctx, { color: "#ff0", size: this.height * 0.6 });
        break;
      case "red":
        red.draw(ctx, { color: "#f00", size: this.height * 0.6 });
        break;
    }
  }
}