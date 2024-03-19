class Point {
  constructor(x, y) {
    this.x = x;
    this.y = y;
  }

  static add(p1, p2) {
    return new Point(p1.x + p2.x, p1.y + p2.y);
  }

  static sub(p1, p2) {
    return new Point(p1.x - p2.x, p1.y - p2.y);
  }
  
  static scale(p, s) {
    return new Point(p.x * s, p.y * s);
  }

  equals(other) {
    return this.x === other.x && this.y === other.y;
  }

  distanceTo(other) {
    const dx = this.x - other.x;
    const dy = this.y - other.y;
    return Math.sqrt(dx * dx + dy * dy);
  }

  draw(ctx, { size = 18, color = "black", outline = false, fill = false } = {}) {
    const radius = size / 2;
    ctx.beginPath();
    ctx.fillStyle = color;
    ctx.arc(this.x, this.y, radius, 0, Math.PI * 2);
    ctx.fill();
    if (outline) {
      ctx.beginPath();
      ctx.lineWidth = 2;
      ctx.strokeStyle = "yellow";
      ctx.arc(this.x, this.y, radius * 0.6, 0, Math.PI * 2);
      ctx.stroke();
    }
    if (fill) {
      ctx.beginPath();
      ctx.arc(this.x, this.y, radius * 0.4, 0, Math.PI * 2);
      ctx.fillStyle = "yellow";
      ctx.fill();
    }
  }
}