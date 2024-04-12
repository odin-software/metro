type RequestStation = {
  id: number;
  position: {
    x: number;
    y: number;
  };
  name: string;
};

type RequestEdge = {
  ID: number;
  Fromid: number;
  Toid: number;
};

type Train = {
  id: number;
};
