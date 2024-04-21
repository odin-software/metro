type RequestTrain = {
  id: number;
  name: string;
  x: number;
  y: number;
  z: number;
  currentId: number;
  make: string;
  line: string;
};

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

type RequestCreateEdge = {
  fromId: number;
  toId: number;
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
