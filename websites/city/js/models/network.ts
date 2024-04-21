import {
  createEdges,
  createStations,
  getEdges,
  getEdgesPoints,
  getStations,
} from "../load.js";
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
  lines: Record<string, Station[]>;
  draftNodes: Station[];
  draftEdges: Edge[];

  constructor(
    nodes: Station[] = [],
    edges: Edge[] = [],
    lines: Record<string, Station[]> = {}
  ) {
    this.nodes = nodes;
    this.edges = edges;
    this.lines = lines;
    this.draftNodes = [];
    this.draftEdges = [];
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

  async saveDrafts() {
    try {
      const stationsToCreate = this.draftNodes.map((val) => {
        return {
          name: val.name,
          x: val.position.x,
          y: val.position.y,
          z: 0,
        };
      });
      const newStations = await createStations(stationsToCreate);
      const edgesToCreate = this.draftEdges.map((val) => {
        return {
          fromId: newStations.find(
            (newSt) =>
              val.start.position.x === newSt.position.x &&
              val.start.position.y === newSt.position.y
          ).id,
          toId: newStations.find(
            (newSt) =>
              val.end.position.x === newSt.position.x &&
              val.end.position.y === newSt.position.y
          ).id,
        };
      });
      await createEdges(edgesToCreate);
      this.nodes.push(...this.draftNodes);
      this.draftNodes.length = 0;
      this.edges.push(...this.draftEdges);
      this.draftEdges.length = 0;
    } catch (err) {
      console.error(err);
    }
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
   * Adds a node to the draft node list on the network.
   * @param {Station} node
   */
  addDraftNode(node: Station) {
    this.draftNodes.push(node);
  }

  /**
   * Checks if a node is a draft in the network.
   * @param {Station} node
   * @returns {boolean}
   */
  checkNodeIsDraft(node: Station): boolean {
    return this.draftNodes.some((n) => n.equalsDraft(node));
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
    const draftEdges = this.getDraftEdgesWithNode(node);
    for (const edge of edges) {
      this.removeEdge(edge);
    }
    for (const edge of draftEdges) {
      this.removeDraftEdge(edge);
    }
    const index = this.nodes.indexOf(node);
    this.nodes.splice(index, 1);
  }

  /**
   * Removes a draft node and it's related edges from the network.
   * @param {Station} node
   */
  removeDraftNode(node: Station) {
    const edges = this.getEdgesWithNode(node);
    const draftEdges = this.getDraftEdgesWithNode(node);
    for (const edge of edges) {
      this.removeEdge(edge);
    }
    for (const edge of draftEdges) {
      this.removeDraftEdge(edge);
    }
    const index = this.draftNodes.indexOf(node);
    this.draftNodes.splice(index, 1);
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
    return this.edges.some((ed) => ed.equals(edge));
  }

  /**
   * DRAFT:
   * Checks if the edge exist checking the draft way in the network.
   * @param edge
   * @returns
   */
  containsDraftEdge(edge: Edge) {
    return this.edges.some((ed) => ed.equalsDraft(edge));
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
   * Tries to add an edge in the network if it doesn't exist.
   * @param {Edge} edge
   * @returns {boolean} whether it was added or not
   */
  tryAddEdgeDraft(edge: Edge): boolean {
    if (this.containsDraftEdge(edge)) {
      return false;
    }
    if (edge.start.equalsDraft(edge.end)) {
      return false;
    }
    this.draftEdges.push(edge);
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
   * Removes a draft edge from the network.
   * @param {Edge} edge
   */
  removeDraftEdge(edge: Edge) {
    const index = this.draftEdges.indexOf(edge);
    this.draftEdges.splice(index, 1);
  }

  /**
   * Gets the edges that have the node in either the start or the end.
   * @param node
   */
  getEdgesWithNode(node: Station): Edge[] {
    return this.edges.filter((edge) => edge.includes(node));
  }

  /**
   * Gets the edges that have the node in either the start or the end.
   * @param node
   */
  getDraftEdgesWithNode(node: Station): Edge[] {
    return this.draftEdges.filter((edge) => edge.includesDraft(node));
  }

  /**
   * Resets the graph by unloading the nodes and edges array.
   */
  dispose() {
    this.nodes.length = 0;
    this.edges.length = 0;
  }

  /**
   * Finds a node from its position vector.
   * @param pos
   * @returns
   */
  getNodeFromPosition(pos: Point): Station {
    return this.nodes.find(
      (node) => node.position.x === pos.x && node.position.y === pos.y
    );
  }

  /**
   * Finds a draft node from its position vector.
   * @param pos
   * @returns
   */
  getDraftNodeFromPosition(pos: Point): Station {
    return this.draftNodes.find(
      (node) => node.position.x === pos.x && node.position.y === pos.y
    );
  }

  draw(ctx: CanvasRenderingContext2D, draft = false) {
    for (const edge of this.edges) {
      edge.draw(ctx, { color: "white", dash: [], width: 1 });
    }
    for (const node of this.nodes) {
      node.draw(ctx);
    }
    if (draft) {
      for (const edge of this.draftEdges) {
        edge.draw(ctx, { color: "white", dash: [], width: 1 });
      }
      for (const node of this.draftNodes) {
        node.draw(ctx);
      }
    }
  }
}
