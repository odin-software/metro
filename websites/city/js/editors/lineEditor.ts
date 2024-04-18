import { getNearestPoint } from "../math/utils.js";
import { Edge } from "../models/edge.js";
import { Network } from "../models/network.js";
import { Station } from "../models/station.js";
import Point from "../primitives/point.js";
import Viewport from "../viewport.js";

export class LineEditor {
  viewport: Viewport;
  canvas: HTMLCanvasElement;
  network: Network;
  ctx: CanvasRenderingContext2D;

  selected: Station | null;
  hovered: Station | null;
  dragging: boolean;
  mouse: Point | null;

  boundMouseDown: (e: MouseEvent) => void;
  boundMouseMove: (e: MouseEvent) => void;
  boundMouseUp: (e: MouseEvent) => void;
  boundContextMenu: (e: MouseEvent) => void;

  constructor(viewport: Viewport, network: Network) {
    this.viewport = viewport;
    this.canvas = viewport.canvas;
    this.network = network;
    this.ctx = viewport.canvas.getContext("2d");

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
    this.boundMouseDown = (e) => this.#handleMouseDown(e);
    this.boundMouseMove = (e) => this.#handleMouseMove(e);
    this.boundMouseUp = (_) => (this.dragging = false);
    this.boundContextMenu = (e) => e.preventDefault();
    this.canvas.addEventListener("mousedown", this.boundMouseDown);
    this.canvas.addEventListener("mousemove", this.boundMouseMove);
    this.canvas.addEventListener("mouseup", this.boundMouseUp);
    this.canvas.addEventListener("contextmenu", this.boundContextMenu);
  }

  #removeEventListeners() {
    this.canvas.removeEventListener("mousedown", this.boundMouseDown);
    this.canvas.removeEventListener("mousemove", this.boundMouseMove);
    this.canvas.removeEventListener("mouseup", this.boundMouseUp);
    this.canvas.removeEventListener("contextmenu", this.boundContextMenu);
  }

  #handleMouseDown(e: MouseEvent) {
    if (e.button == 2) {
      // right click
      if (this.selected) {
        this.selected = null;
      } else if (this.hovered) {
        this.#removePoint(this.hovered);
      }
    }
    if (e.button == 0) {
      // left click
      if (this.hovered) {
        this.#selectPoint(this.hovered);
        this.selected = this.hovered;
        this.dragging = true;
        return;
      }
      const st = Station.draft(this.mouse.x, this.mouse.y);
      this.network.addNode(st);
      this.#selectPoint(st);
      this.selected = st;
      this.hovered = st;
    }
  }

  #handleMouseMove(e: MouseEvent) {
    this.mouse = this.viewport.getMouse(e, true);
    const val = getNearestPoint(
      this.mouse,
      this.network.nodes.map((st) => st.position),
      10 * this.viewport.zoom
    );
    this.hovered = val ? this.network.getNodeFromPosition(val) : null;
    if (this.dragging) {
      this.selected.position.x = this.mouse.x;
      this.selected.position.y = this.mouse.y;
    }
  }

  #removePoint(st: Station) {
    this.network.removeNode(st);
    this.hovered = null;
    if (this.selected === st) {
      this.selected = null;
    }
  }

  #selectPoint(st: Station) {
    if (this.selected) {
      this.network.tryAddEdgeDraft(new Edge(this.selected, st, []));
    }
  }

  display() {
    this.network.draw(this.ctx);

    if (this.hovered) {
      this.hovered.position.draw(this.ctx, {
        fill: true,
      });
    }

    if (this.selected) {
      const mockStation = Station.draft(this.mouse.x, this.mouse.y);
      const intent = this.hovered ? this.hovered : mockStation;
      new Edge(this.selected, intent, []).draw(this.ctx, { dash: [3, 3] });
      this.selected.position.draw(this.ctx, {
        outline: true,
      });
    }
  }

  dispose() {
    this.network.dispose();
    this.selected = null;
    this.hovered = null;
    this.dragging = false;
  }
}
