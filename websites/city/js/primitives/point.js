class Point {
  /**
   * The point class represents a vector, it contains static methods
   * for basic vector calculations.
   * @param {number} x
   * @param {number} y
   * @param {string} name - the name representing this point.
   */
  constructor(x, y, name = "") {
    this.x = x;
    this.y = y;
    this.name = name;
  }

  /**
   * Adds two points and returns a new Point with the sum.
   * @param {Point} p1
   * @param {Point} p2
   * @returns {Point}
   */
  static add(p1, p2) {
    return new Point(p1.x + p2.x, p1.y + p2.y);
  }

  /**
   * Substracts two points and returns a new Point with the difference.
   * @param {Point} p1
   * @param {Point} p2
   * @returns {Point}
   */
  static sub(p1, p2) {
    return new Point(p1.x - p2.x, p1.y - p2.y);
  }

  /**
   * Scalar multiplication of a vector, returns a new one.
   * @param {Point} p
   * @param {number} s - number to multiply both components of the point.
   * @returns {Point}
   */
  static scale(p, s) {
    return new Point(p.x * s, p.y * s);
  }

  /**
   * Asserts whether the specified point is equal to this one.
   * @param {Point} p
   * @returns {boolean}
   */
  equals(p) {
    return this.x === p.x && this.y === p.y;
  }

  /**
   * Returns the distance between the specified point and this one.
   * @param {Point} p
   * @returns {number}
   */
  distanceTo(p) {
    const dx = this.x - p.x;
    const dy = this.y - p.y;
    return Math.sqrt(dx * dx + dy * dy);
  }

  /**
   * Styling options for drawing a point.
   * @typedef {Object} Styles
   * @property {number} [size=18] - diameter of the point
   * @property {string} [color=black] - color of the point
   * @property {boolean} [outline=false] - whether the point has an outline
   * @property {boolean} [fill=false] - whether the point has an fill
   */
  /**
   * Function to draw a Point with options on styling.
   * @param {CanvasRenderingContext2D} ctx
   * @param {Styles} style
   */
  draw(
    ctx,
    { size = 18, color = "black", outline = false, fill = false } = {}
  ) {
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

export default Point;
