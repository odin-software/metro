import Point from "../primitives/point.js";
import Segment from "../primitives/segment.js";
import { Edge } from "./edge.js";
import { Station } from "./station.js";

/**
 * Representation of a network. It has nodes and edges, plus methods
 * to interact with those.
 */
export class Network {
  nodes: Station[];
  edges: Edge[];

  constructor(nodes: Station[] = [], edges: Edge[] = []) {
    this.nodes = nodes;
    this.edges = edges;
  }

  /**
   * String representation of this networkj.
   * @returns {string}
   */
  hash(): string {
    return JSON.stringify(this);
  }

  /**
   * Adds a node to the network.
   * @param {Station} node
   */
  addNode(node: Station) {
    this.nodes.push(node);
  }

  /**
   * Checks if a node is contained within the network.
   * @param {Station} node
   * @returns {boolean}
   */
  containsNode(node: Station): boolean {
    return this.nodes.some((n) => n.equals(node));
  }

  /**
   * Tries to add a node in the network if it doesn't exist.
   * @param {Station} node
   * @returns {boolean}
   */
  tryAddNode(node: Station): boolean {
    if (this.containsNode(node)) {
      return false;
    }
    this.addNode(node);
    return true;
  }

  /**
   * Removes a node and it's related edges from the network.
   * @param {Station} node
   */
  removeNode(node: Station) {
    const segs = this.getEdgesWithNode(node);
    for (const seg of segs) {
      this.removeSegment(seg);
    }
    const index = this.nodes.indexOf(node);
    this.nodes.splice(index, 1);
  }

  /**
   * Adds an edge to the network.
   * @param {Edge} edge
   */
  addEdge(edge: Edge) {
    this.edges.push(edge);
  }

  /**
   *
   * @param edge
   * @returns
   */
  containsEdge(edge: Edge) {
    return this.edges.find((edge) => edge.equals(edge));
  }

  tryAddSegment(segment) {
    if (this.containsSegment(segment)) {
      return false;
    }
    if (segment.p1.equals(segment.p2)) {
      return false;
    }

    this.segments.push(segment);
    return true;
  }

  removeSegment(segment) {
    const index = this.segments.indexOf(segment);
    this.segments.splice(index, 1);
  }

  /**
   * Gets the edges that have the node in either the start or the end.
   * @param node
   */
  getEdgesWithNode(node: Station): Edge[] {
    return this.edges.filter((edge) => edge.includes(node));
  }

  /**
   * Resets the graph by unloading the points and segments array.
   */
  dispose() {
    this.points.length = 0;
    this.segments.length = 0;
  }

  draw(ctx) {
    for (const seg of this.segments) {
      seg.draw(ctx, {});
    }

    for (const point of this.points) {
      point.draw(ctx, { size: 20, color: "white" });
    }
  }
}
