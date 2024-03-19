class Start extends Marking {
  constructor(center, directionVector, width, height) {
    super(center, directionVector, width, height);

    this.type = "start";
  }

  draw(ctx) {
    ctx.save();
    ctx.translate(this.center.x, this.center.y);
    ctx.rotate(angle(this.directionVector) - Math.PI / 2);

    ctx.fill = "white";
    ctx.rect(0, 0, 30, 30);
    ctx.restore();
  }
}