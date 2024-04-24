import { Network } from "../models/network.js";
import { Train } from "../models/train.js";

export type TrainStore = {
  trains: Train[];
};

export type DialogStore = {
  open: boolean;
  title: string;
  body: string;
  input: string;
  yesBtn: () => void;
  noBtn: () => void;
};

export type NetworkStore = {
  network: Network;
};
