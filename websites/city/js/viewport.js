class Viewport {
  constructor(canvas, zoom = 1, offset = null) {
    this.canvas = canvas;
    this.ctx = canvas.getContext('2d');

    this.zoom = zoom;
    this.center = new Point(canvas.width / 2, canvas.height / 2);
    this.offset = offset ? offset : Point.scale(this.center, -1);
    this.drag = {
      start: new Point(0, 0),
      end: new Point(0, 0),
      offset: new Point(0, 0),
      active: false
    }

    this.#addEventListeners();
  }

  getMouse(e, substractDragOffset = false) {
    const p = new Point(
      (e.offsetX - this.center.x) * this.zoom - this.offset.x, 
      (e.offsetY - this.center.y) * this.zoom - this.offset.y
    );

    return substractDragOffset ? Point.sub(p, this.drag.offset) : p;
  }

  getOffset() {
    return Point.add(this.offset, this.drag.offset);
  }

  reset() {
    this.ctx.restore();

    this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);

    this.ctx.save();
    this.ctx.translate(this.center.x, this.center.y);

    this.ctx.scale(1 / this.zoom, 1 / this.zoom);
    const offset = this.getOffset();
    this.ctx.translate(offset.x, offset.y);
  }

  #addEventListeners() {
    this.canvas.addEventListener('mousewheel', e =>  this.#onMouseWheel(e));
    this.canvas.addEventListener('mousedown', e =>  this.#onMouseDown(e));
    this.canvas.addEventListener('mousemove', e =>  this.#onMouseMove(e));
    this.canvas.addEventListener('mouseup', e =>  this.#onMouseUp(e));
  }

  #onMouseWheel(e) {
    const dir = Math.sign(e.deltaY);
    const step = 0.1;
    this.zoom += dir * step;
    this.zoom = Math.max(1, Math.min(5, this.zoom));
  }

  #onMouseDown(e) {
    if (e.button == 1) {
      this.drag.start = this.getMouse(e);
      this.drag.active = true;
    }
  }

  #onMouseMove(e) {
    if (this.drag.active) {
      this.drag.end = this.getMouse(e);
      this.drag.offset = Point.sub(this.drag.end, this.drag.start);
    }
  }

  #onMouseUp(e) {
    if (this.drag.active) {
      this.offset = Point.add(this.offset, this.drag.offset);
      this.drag = {
        start: new Point(0, 0),
        end: new Point(0, 0),
        offset: new Point(0, 0),
        active: false
      };
    }
  }
}