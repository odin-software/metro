import { getEdges, getEdgesPoints, getStations } from "../load.js";
import Point from "../primitives/point.js";
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
   * String representation of this network.
   * @returns {string}
   */
  hash(): string {
    return JSON.stringify(this);
  }

  /**
   * Loads the network from the server.
   */
  static async load() {
    const stationsResponse = await getStations();
    const edgesResponse = await getEdges();
    const stations = stationsResponse.map(
      (st) => new Station(st.id, st.name, st.position.x, st.position.y)
    );
    const edges = await Promise.all(
      edgesResponse.map(async (e) => {
        const st1 = stations.find((st) => st.id === e.Fromid);
        const st2 = stations.find((st) => st.id === e.Toid);
        const eps = await getEdgesPoints(e.ID);
        return new Edge(
          st1,
          st2,
          eps ? eps.map((ep) => new Point(ep.X, ep.Y)) : []
        );
      })
    );
    return new Network(stations, edges);
  }

  /**
   * Get a center point of the network to use as offset of the viewport.
   */
  getCenterPoint(): Point {
    const x =
      this.nodes.reduce((prev, cur) => {
        return cur.position.x + prev;
      }, 0) / this.nodes.length;
    const y =
      this.nodes.reduce((prev, cur) => {
        return cur.position.y + prev;
      }, 0) / this.nodes.length;
    return new Point(x, y);
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
   * @returns {boolean} whether it was added or not
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
    const edges = this.getEdgesWithNode(node);
    for (const edge of edges) {
      this.removeEdge(edge);
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
   * Checks if the edge exist in the network.
   * @param edge
   * @returns
   */
  containsEdge(edge: Edge) {
    return this.edges.some((edge) => edge.equals(edge));
  }

  /**
   * Tries to add an edge in the network if it doesn't exist.
   * @param {Edge} edge
   * @returns {boolean} whether it was added or not
   */
  tryAddEdge(edge: Edge): boolean {
    if (this.containsEdge(edge)) {
      return false;
    }
    if (edge.start.equals(edge.end)) {
      return false;
    }
    this.edges.push(edge);
    return true;
  }

  /**
   * Removes an edge from the network.
   * @param {Edge} edge
   */
  removeEdge(edge: Edge) {
    const index = this.edges.indexOf(edge);
    this.edges.splice(index, 1);
  }

  /**
   * Gets the edges that have the node in either the start or the end.
   * @param node
   */
  getEdgesWithNode(node: Station): Edge[] {
    return this.edges.filter((edge) => edge.includes(node));
  }

  /**
   * Resets the graph by unloading the nodes and edges array.
   */
  dispose() {
    this.nodes.length = 0;
    this.edges.length = 0;
  }

  getNodeFromPosition(pos: Point): Station {
    return this.nodes.find(
      (node) => node.position.x === pos.x && node.position.y === pos.y
    );
  }

  draw(ctx: CanvasRenderingContext2D) {
    for (const edge of this.edges) {
      edge.draw(ctx, undefined);
    }

    for (const node of this.nodes) {
      node.draw(ctx);
    }
  }
}
