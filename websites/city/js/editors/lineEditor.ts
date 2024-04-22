import { getNearestPoint } from "../math/utils.js";
import { Network } from "../models/network.js";
import { Station } from "../models/station.js";
import Point from "../primitives/point.js";
import Viewport from "../viewport.js";
import DialogStore from "../store/dialog.js";
import { createLine } from "../load.js";

export class LineEditor {
  viewport: Viewport;
  canvas: HTMLCanvasElement;
  network: Network;
  ctx: CanvasRenderingContext2D;
  line: Station[];
  lineName: string;
  frames: number;

  selected: Station | null;
  hovered: Station | null;
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
    this.frames = 0;

    this.line = [];
    this.lineName = "";
    this.selected = null;
    this.hovered = null;
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
    this.boundContextMenu = (e) => e.preventDefault();
    this.canvas.addEventListener("mousedown", this.boundMouseDown);
    this.canvas.addEventListener("mousemove", this.boundMouseMove);
    this.canvas.addEventListener("mouseup", this.boundMouseUp);
    this.canvas.addEventListener("contextmenu", this.boundContextMenu);
  }

  #removeEventListeners() {
    this.canvas.removeEventListener("mousedown", this.boundMouseDown);
    this.canvas.removeEventListener("mousemove", this.boundMouseMove);
    this.canvas.removeEventListener("contextmenu", this.boundContextMenu);
  }

  #handleMouseDown(e: MouseEvent) {
    if (e.button == 2) {
      if (this.line.length > 2) {
        DialogStore.dispatch("openDialog", {
          open: true,
          title: "Line",
          body: `Line ${this.lineName}`,
          yesBtn: async () => {
            const createLineReq = this.#createLineRequest();
            await createLine(createLineReq);
            this.line.length = 0;
            DialogStore.dispatch("closeDialog", {});
          },
          noBtn: () => DialogStore.dispatch("closeDialog", {}),
        });
      } else if (this.selected) {
        // right click
        this.selected = null;
      }
    }
    if (e.button == 0) {
      // left click
      if (this.hovered) {
        this.#selectPoint(this.hovered);
        this.selected = this.hovered;
        return;
      }
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
  }

  #selectPoint(st: Station) {
    if (this.selected && this.network.areConnected(this.selected, st)) {
      if (this.line.length == 0) {
        this.lineName = `${Math.random() * 1293}`;
        this.line.push(this.selected);
      }
      this.line.push(st);
    }
  }

  #createLineRequest(): RequestCreateLine {
    const obj: RequestCreateLine = {
      name: this.lineName,
      stations: this.line.map((val, idx) => {
        return {
          stationId: val.id,
          odr: idx + 1,
        };
      }),
    };
    return obj;
  }

  display() {
    this.network.draw(this.ctx, true);

    if (this.hovered) {
      this.hovered.position.draw(this.ctx, {
        fill: true,
      });
    }

    if (this.selected) {
      const nodes = this.network.getConnectedNodes(this.selected);
      const animatedSize = Math.abs(Math.sin(this.frames) * 20);
      nodes.forEach((nd) =>
        nd.draw(this.ctx, { size: animatedSize, color: "green" })
      );
      this.selected.position.draw(this.ctx, {
        outline: true,
      });
      this.frames += 0.04;
    }
  }

  dispose() {
    this.network.dispose();
    this.selected = null;
    this.hovered = null;
  }
}
