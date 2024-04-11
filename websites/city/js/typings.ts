import Graph from "./math/graph";
import Point from "./primitives/point";

// Primitives
export type PointStyle = {
  size?: number;
  color?: string;
  outline?: boolean;
  fill?: boolean;
};

export type SegmentStyle = {
  width?: number;
  color?: string;
  dash?: number[];
};

// Others
export type WorldInfo = {
  graph: Graph;
  zoom: number;
  offset: Point;
};
