class GraphEditor {
  constructor(viewport, graph) {
    this.viewport = viewport;
    this.canvas = viewport.canvas;
    this.graph = graph;
    this.ctx = viewport.canvas.getContext('2d');
    
    this.selected = null;
    this.hovered = null;
    this.dragging = false;
    this.mouse = null;
    // this.selectedSegment = null;
    // this.dragOffset = null;
  }

  enable() {
    this.#addEventListeners();
  }

  disable() {
    this.#removeEventListeners();
    this.selected = null;
    this.hovered = null;
  }

  #addEventListeners() {
    this.boundMouseDown = e => this.#handleMouseDown(e);
    this.boundMouseMove = e => this.#handleMouseMove(e);
    this.boundMouseUp = _ => this.dragging = false;
    this.boundContextMenu = e => e.preventDefault();
    this.canvas.addEventListener('mousedown', this.boundMouseDown);
    this.canvas.addEventListener('mousemove', this.boundMouseMove);
    this.canvas.addEventListener('mouseup', this.boundMouseUp);
    this.canvas.addEventListener('contextmenu', this.boundContextMenu);
  }

  #removeEventListeners() {
    this.canvas.removeEventListener('mousedown', this.boundMouseDown);
    this.canvas.removeEventListener('mousemove', this.boundMouseMove);
    this.canvas.removeEventListener('mouseup', this.boundMouseUp);
    this.canvas.removeEventListener('contextmenu', this.boundContextMenu);
  }

  #handleMouseDown(e) {
    if (e.button == 2) { // right click
      if (this.selected) {
        this.selected = null;
      } else if (this.hovered) {
        this.#removePoint(this.hovered);
      }
    }
    if (e.button == 0) { // left click
      if (this.hovered) {
        this.#selectPoint(this.hovered);
        this.selected = this.hovered;
        this.dragging = true;
        return;
      }
      this.graph.addPoint(this.mouse);
      this.#selectPoint(this.mouse);
      this.selected = this.mouse;
      this.hovered = this.mouse;
    }
  }

  #handleMouseMove(e) {
    this.mouse = this.viewport.getMouse(e, true); 
    this.hovered = getNearestPoint(this.mouse, this.graph.points, 10 * this.viewport.zoom);
    if (this.dragging) {
      this.selected.x = this.mouse.x;
      this.selected.y = this.mouse.y;
    }
  }

  #removePoint(point) {
    this.graph.removePoint(point);
    this.hovered = null;
    if (this.selected === point) {
      this.selected = null;
    }
  }

  #selectPoint(point) {
    if (this.selected) {
      this.graph.tryAddSegment(new Segment(this.selected, point));
    }
  }

  display() {
    this.graph.draw(this.ctx);

    if (this.hovered) {
      this.hovered.draw(this.ctx, {
        fill: true,
      });
    }

    if (this.selected) {
      const intent = this.hovered ? this.hovered : this.mouse;
      new Segment(this.selected, intent).draw(this.ctx, { dash: [3, 3] });
      this.selected.draw(this.ctx, {
        outline: true,
      });
    }
  }

  dispose() {
    this.graph.dispose();
    this.selected = null;
    this.hovered = null;
    this.dragging = false;
  }
}