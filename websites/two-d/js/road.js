class Road {
  constructor(x, width, laneCount = 3) {
    this.x = x;
    this.width = width;
    this.laneCount = laneCount; 

    this.left = x - width / 2;
    this.right = x + width / 2;

    const inifity = 10000000;
    this.top = -inifity;
    this.bottom = inifity;

    const topLeft = { x: this.left, y: this.top };
    const topRight = { x: this.right, y: this.top };
    const bottomLeft = { x: this.left, y: this.bottom };
    const bottomRight = { x: this.right, y: this.bottom };

    this.borders = [
      [topLeft, bottomLeft],
      [topRight, bottomRight],
    ]
  }

  getLaneCenter(laneIdx) {
    const laneWidth = this.width / this.laneCount;
    const min = Math.min(laneIdx, this.laneCount - 1);
    return this.left + laneWidth/2 + min * laneWidth;
  }

  draw(ctx) {
    ctx.lineWidth = 5;
    ctx.strokeStyle = 'white';

    // Draw lanes
    for (let i = 1; i <= this.laneCount-1; i++) {
      const x = lerp(this.left, this.right, i / this.laneCount)

      ctx.setLineDash([20, 20]);
      ctx.beginPath();
      ctx.moveTo(x, this.top);
      ctx.lineTo(x, this.bottom);
      ctx.stroke();
    }

    ctx.setLineDash([]);
    this.borders.forEach(([start, end]) => { 
      ctx.beginPath();
      ctx.moveTo(start.x, start.y);
      ctx.lineTo(end.x, end.y);
      ctx.stroke();
    });
  }
}
