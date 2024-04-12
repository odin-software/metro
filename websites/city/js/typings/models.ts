type Station = {
  id: number;
  position: {
    x: number;
    y: number;
  };
  name: string;
};

type Edge = {
  ID: number;
  Fromid: number;
  Toid: number;
};

type Train = {
  id: number;
};
