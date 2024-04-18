type RequestStation = {
  id: number;
  position: {
    x: number;
    y: number;
  };
  name: string;
};

type RequestCreateStation = {
  name: string;
  x: number;
  y: number;
  z: number;
};

type RequestEdge = {
  ID: number;
  Fromid: number;
  Toid: number;
};

type RequestEdgePoint = {
  ID: number;
  Edgeid: number;
  X: number;
  Y: number;
  Z: number;
  Odr: number;
};

type RequestLine = {
  name: string;
  points: {
    x: number;
    y: number;
  }[];
};

type Train = {
  id: number;
};
