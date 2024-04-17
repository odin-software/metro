import Graph from "./math/graph.js";
import Point from "./primitives/point.js";

// Others
export type WorldInfo = {
  graph: Graph;
  zoom: number;
  offset: Point;
};
